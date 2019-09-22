package anticipation

import (
	"bytes"
	"debug/elf"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/pkg/term"
)

const normalSize = 32
const maxDataLineFails = 2 //this is multiplied by maxFailsOneLine for # fails on a single data line
const maxFailsOneLine = 5
const maxCommandFails = 5 // usually the other side is in a bad state when this happens
const maxPings = 25
const sanityBlockSize = 0x20
const notificationInterval = 0x800
const pageSize = 0x10000

var debugRcvd = false
var debugRcvdChars = false
var debugSent = false

var cleanupNeeded = []*term.Term{}

var signalChan = make(chan os.Signal, 1)

type antcState int

const (
	waiting antcState = iota
	sectionSend
	done
)

const ready = "?"
const sectionTransmit = "section"
const fetchTransmit = "fetch"
const inflateTransmit = "inflate"
const launchTransmit = "launch"
const pingTransmit = "ping"

var state = waiting

type sectionInfo struct {
	offset   uint64
	physAddr uint64
	size     uint64
}

// Main is where the action is.
func Main(device string, kpath string, byteLimit int) {

	log.SetOutput(os.Stdout)
	var terminal io.ReadWriter

	if strings.HasPrefix(device, "/dev/") {
		terminal = setupTTY(device, false)
	} else {
		file, err := os.Open(device)
		if err != nil {
			log.Fatalf("unable to open %s: %v", device, err)
		}
		terminal = file
	}

	pingLoop(terminal)

	//
	// Open the kernel file for our transmissions
	//
	in, err := os.Open(kpath)
	if err != nil {
		fatalf(terminal, "unable to open %s: %v", kpath, err)
	}

	// end of the program is the end of what we sent to the bootloader (highest address used)
	endOfProgram := uint64(0)
	//
	// these are the sections that require us to copy bytes to target in hex format
	// executable, rwdata, rodata
	//
	var execInfo, rwInfo, roInfo *sectionInfo
	names := []string{"executable", "rw data", "readonly data"}
	execInfo, err = getExecutable(kpath)
	if err != nil {
		fatalf(terminal, "loading executable: %v", err)
	}
	if execInfo.physAddr+execInfo.size > endOfProgram {
		endOfProgram = execInfo.physAddr + execInfo.size
	}
	entryPoint := execInfo.physAddr
	rwInfo, err = getRWData(kpath)
	if err != nil {
		log.Printf("WARNING: no RW section found in the elf file %v", kpath)
		//fatalf(terminal, "loading rw: %v", err)
	} else {
		if rwInfo.physAddr+rwInfo.size > endOfProgram {
			endOfProgram = rwInfo.physAddr + rwInfo.size
		}
	}
	roInfo, err = getSectionByName(kpath, ".rodata")
	if err != nil {
		fatalf(terminal, "loading ro: %v", err)
	}
	if roInfo.physAddr+roInfo.size > endOfProgram {
		endOfProgram = roInfo.physAddr + roInfo.size
	}
	for i, info := range []*sectionInfo{execInfo, rwInfo, roInfo} {
		if info == nil {
			continue
		}
		log.Printf("sending %s: 0x%04x bytes total, starts at offset 0x%04x in elf file\n", names[i], info.size, info.offset)
		if !sendSection(terminal, info, in) {
			fatalf(terminal, "giving up trying to send %s", names[i])
		}
		sanityCheck(terminal, info, in, names[i])
	}

	//
	// Inflate the BSS data and zero it
	//
	info, err := getSectionByName(kpath, ".bss")
	if err != nil {
		fatalf(terminal, "unable to load the .bss section: %v", err)
	}
	if !sendInflate(terminal, info.physAddr, info.size) {
		fatalf(terminal, "unable to inflate bss section, giving up")
	}
	if info.physAddr+info.size > endOfProgram {
		endOfProgram = info.physAddr + info.size
	}

	mask := ^int64(pageSize - 1)
	heap := (endOfProgram + pageSize) & uint64(mask) //upward
	stack := heap + pageSize
	log.Printf("highest address used by loaded program: 0x%04x... launching at %04x with SP %04x and heap at %04x", endOfProgram, entryPoint, stack, heap)
	if !sendLaunch(terminal, entryPoint, stack, heap) {
		fatalf(terminal, "failed to launch, giving up")
	}
	simpleTerminal(terminal, device, byteLimit)
}

const sanityCheckSize = 0x100

func sanityCheck(terminal io.ReadWriter, prog *sectionInfo, in *os.File, name string) {
	// do a sanity check on the first few bytes.... seek back to origin
	_, err := in.Seek(int64(prog.offset), 0)
	if err != nil {
		fatalf(terminal, "failed to seek to do sanity check: %v", err)
	}
	minSize := uint64(sanityCheckSize)
	if prog.size < minSize {
		minSize = prog.size
	}

	// run some records through fetch and compare to disk version
	for p := prog.physAddr; p < prog.physAddr+minSize; p += sanityBlockSize {
		ok, line := sendFetch(terminal, p)
		if !ok {
			fatalf(terminal, "fetch failed, giving up")
		}
		rcvd, err := hex.DecodeString(line)
		if err != nil {
			fatalf(terminal, "unable to decode hex from bl: %v", err)
		}
		onDisk := make([]byte, len(rcvd))
		n, err := in.Read(onDisk)
		if err != nil || n != len(rcvd) {
			fatalf(terminal, "unable to read the disk to do comparison: %v", err)
		}
		for i := 0; i < len(rcvd) && uint64(i) < minSize && p+uint64(i) < prog.physAddr+minSize; i++ {
			if rcvd[i] != onDisk[i] {
				fatalf(terminal, "byte mismatch found at %04x: %x expected but got %x (iteration %d, minsize %d, %d vs %d)", int(p)+i, onDisk[i], rcvd[i], i, minSize, len(onDisk), len(rcvd))
			}
		}
	}
	log.Printf("sanity check completed on '%s' (checked first 0x%04x bytes)", name, minSize)
}

func sendSection(terminal io.ReadWriter, info *sectionInfo, fp *os.File) bool {
	switch state {
	case waiting:
		line := readLine(terminal)
		if line == ready {
			if sendSectionCommand(terminal) {
				state = sectionSend
			} else {
				//failed to send the command, user will see us trying a few times
				return false
			}
		}
		fallthrough
	case sectionSend:
		if err := transmitSection(terminal, info, fp); err == false {
			return false
		}
		state = waiting
		return true
	default:
		fatalf(terminal, "unknown anticipation state: %d", int(state))
	}
	return true //should never happen
}

func cleanupAndExit() {
	for _, terminal := range cleanupNeeded {
		terminal.Write([]byte{0x10, 0x21, 0x10, 0x04})
		terminal.Restore()
	}
	os.Exit(0)
}

func readOneByte(t io.ReadWriter) byte {
	c := make([]byte, 1)
	n, err := t.Read(c)
	if err != nil {
		fatalf(t, "error reading character from terminal: %v", err)
	}
	if n == 0 && err == io.EOF {
		cleanupAndExit()
	}
	if debugRcvdChars {
		if c[0] == 10 {
			log.Printf("<<<<< debugChar: LF")
		} else {
			log.Printf("<<<< debugChar: %c", c[0])
		}
	}
	return c[0]
}

func readLine(t io.ReadWriter) string {
	var buffer bytes.Buffer
	for {
		c := readOneByte(t)
		if c == 10 {
			break
		}
		if err := buffer.WriteByte(c); err != nil {
			fatalf(t, "error writing to buffer: %v", err)
		}
	}
	l := buffer.String()
	if len(l) == 0 {
		if debugRcvd {
			log.Printf("<----------- EMPTY LINE received!")
		}
		return l
	}
	if debugRcvd {
		log.Printf("<----------- %s", l)
	}
	return l
}

func fatalf(t io.ReadWriter, s string, args ...interface{}) {
	log.Printf(s, args...)
	cleanupAndExit()
}

func transmitSection(t io.ReadWriter, info *sectionInfo, fp *os.File) bool {
	//setup file pointer
	ret, err := fp.Seek(int64(info.offset), 0)
	if err != nil || ret != int64(info.offset) {
		fatalf(t, "bad seek: %v (ret was %d, size is %d)", err, ret, info.offset)
	}

	return transmitFile(t, info.physAddr, info.size, fp)
}

func transmitFile(t io.ReadWriter, paddr uint64, filesz uint64, fp *os.File) bool {
	if !sendHexESA(t, paddr) {
		return false
	}

	current := uint64(0)
	buffer := make([]byte, normalSize)

	//we will try data lines up to 5 times
	fails := 0

	for current < filesz && fails < maxDataLineFails {
		size := 32 //normal case
		if int(filesz-current) < 32 {
			size = int(filesz - current)
		}
		n, err := fp.Read(buffer[:size])
		if n != int(size) || err != nil {
			fatalf(t, "bad read from data file: %v (with size=%d and bytes read=%d)", err, size, n)
		}
		if !sendHexData(t, size, current, buffer) {
			//move fp back to previous position
			_, err := fp.Seek(int64(-size), 1)
			if err != nil {
				fatalf(t, "bad seek for rewind: %v (n was %d, size is %d)", err, n, size)
			}
			fails++
		} else {
			fails = 0
			if current%notificationInterval == 0 { //we iter in 0x20 increments
				log.Printf("sent block @ 0x%08x!\n", current+paddr)
			}
			current += uint64(size)
		}
	}
	if fails == maxDataLineFails {
		return false
	}
	if !sendHexEOF(t) {
		return false
	}
	log.Printf("completed sending file... %04x bytes", filesz)
	return true
}

// returns the line sent in either case
func confirm(t io.ReadWriter) (bool, string) {
	l := readLine(t)
	if strings.HasPrefix(l, ".") {
		return true, l
	}
	return false, l
}

func writeWithChecksum(t io.ReadWriter, payload string) {
	decoded, err := hex.DecodeString(payload)
	if err != nil {
		fatalf(t, "bad hex string: %v", err)
	}
	//figure out checksum
	sum := uint64(0)
	for i := 0; i < len(decoded); i++ {
		sum += uint64(decoded[i])
	}
	complement := ((^sum) + 1)
	checksum := uint8(complement) & 0xff

	line := fmt.Sprintf(":%s%02x\n", payload, checksum)
	b := []byte(line)

	current := 0
	if debugSent {
		log.Printf("------------> %+v\n", b)
	}
	for current < len(b) {
		n, err := t.Write(b[current:])
		if err != nil {
			fatalf(t, "failed writing line to bl: %v", err)
		}
		current += n
	}

}

func sendHexEOF(t io.ReadWriter) bool {
	payload := fmt.Sprintf("00000001")
	ok, _ := sendSingleCommand(t, payload, "EOF", maxFailsOneLine, false)
	return ok
}

func sendHexESA(t io.ReadWriter, paddr uint64) bool {
	executableLocation := paddr >> 4
	if executableLocation > 0xffff {
		fatalf(t, "unable to use hex record type 02 (ESA) because executable physical address too large: %x", paddr)
	}
	loc16 := uint16(executableLocation) //checked above
	payload := fmt.Sprintf("02000002%04x", loc16)
	ok, _ := sendSingleCommand(t, payload, "ESA", maxDataLineFails, false)
	return ok
}

func sendHexData(t io.ReadWriter, size int, current uint64, buffer []byte) bool {
	payload := fmt.Sprintf("%02x%04x00", size, current)
	for i := 0; i < size; i++ {
		payload += fmt.Sprintf("%02x", buffer[i])
	}
	prev := debugSent
	if prev {
		log.Printf("------------> DATA @ 0x%04x00", current)
		debugSent = false
	}
	ok, _ := sendSingleCommand(t, payload, "DATA", maxFailsOneLine, false)
	if prev {
		debugSent = true
	}
	return ok
}

// bool is success/failure and the string is the response in the failure case
func sendSingleCommand(t io.ReadWriter, payload string, name string, maxFails int, quiet bool) (bool, string) {
	tries := 0
	confirmLine := ""
	var ok bool
	for tries < maxFails {
		if isCommand(payload) {
			p := fmt.Sprintf(":%s\n", payload)
			b := []byte(p)
			if debugSent {
				log.Printf("------------> %+v\n", b)
			}
			n, err := t.Write(b)
			if err != nil || n != len(p) {
				fatalf(t, "unable to send command %s (%d bytes sent): %v", payload, n, err)
			}
		} else {
			writeWithChecksum(t, payload)
		}
		ok, confirmLine = confirm(t)
		if ok {
			break
		}
		if !quiet {
			log.Printf("attempt %d of %s command failed, response: %s", tries, name, confirmLine)
		}
		tries++
	}
	if tries == maxFails {
		return false, confirmLine
	}
	return true, confirmLine
}

func setupTTY(device string, cbreak bool) io.ReadWriter {
	// tty shenanigans
	tty, err := term.Open(device)
	if err != nil {
		log.Fatalf("unable to open %s: %v", device, err)
	}

	if err := tty.SetFlowControl(term.NONE); err != nil {
		log.Fatalf("unable to set flow control none:%v", err)
	}
	if err := tty.SetSpeed(115200); err != nil {
		log.Fatalf("unable to set speed:%v", err)
	}

	if cbreak {
		if err := tty.SetCbreak(); err != nil {
			log.Fatalf("unable to set cbreak on %s: %v", device, err)
		}
	} else {
		if err := tty.SetRaw(); err != nil {
			log.Fatalf("unable to set raw on %s: %v", device, err)
		}
		a, err := tty.Available()
		if err != nil {
			log.Fatalf("unable to check Available on %s: %v", device, err)
		}
		log.Printf("available? %d\n", a)
	}
	cleanupNeeded = append(cleanupNeeded, tty)

	signal.Notify(signalChan, os.Interrupt)
	go func(t *term.Term) {
		<-signalChan
		cleanupAndExit()
		os.Exit(0)
	}(tty)
	return tty
}

func sendSectionCommand(t io.ReadWriter) bool {
	ok, _ := sendSingleCommand(t, sectionTransmit, "SECTION", maxCommandFails, false)
	return ok
}

func isCommand(payload string) bool {
	switch {
	case sectionTransmit == payload:
		return true
	case pingTransmit == payload:
		return true
	case strings.HasPrefix(payload, fetchTransmit):
		return true
	case strings.HasPrefix(payload, inflateTransmit):
		return true
	case strings.HasPrefix(payload, launchTransmit):
		return true
	default:
		return false
	}
}

func sendFetch(t io.ReadWriter, addr uint64) (bool, string) {
	cmd := fmt.Sprintf("%s %08x", fetchTransmit, addr)
	ok, response := sendSingleCommand(t, cmd, "FETCH ", maxCommandFails, false)
	if !ok {
		return false, response
	}
	return true, response[len(".ok "):]
}

func sendLaunch(t io.ReadWriter, addr uint64, stack uint64, heap uint64) bool {
	cmd := fmt.Sprintf("%s %08x %08x %08x %08x", launchTransmit, addr, stack, heap, time.Now().Unix())
	ok, _ := sendSingleCommand(t, cmd, "LAUNCH ", maxCommandFails, false)
	if !ok {
		return false
	}
	return true
}

func sendInflate(t io.ReadWriter, addr uint64, size uint64) bool {
	cmd := fmt.Sprintf("%s %08x %08x", inflateTransmit, addr, size)
	ok, _ := sendSingleCommand(t, cmd, "INFLATE ", maxCommandFails, false)
	if !ok {
		return false
	}
	return true
}

func getRWData(kpath string) (*sectionInfo, error) {
	return getProgramSectionByHeader(kpath, false, true, true)
}

func getExecutable(kpath string) (*sectionInfo, error) {
	//return getSectionByName(kpath, ".text")
	return getProgramSectionByHeader(kpath, true, false, true)
}

func getSectionByName(kpath string, sectionName string) (*sectionInfo, error) {
	elfFile, err := elf.Open(kpath)
	if err != nil {
		log.Fatalf("whoa!?!? can't read elf file but checked it before: %v", err)
	}
	defer elfFile.Close()
	s := elfFile.Section(sectionName)
	if s == nil {
		return nil, fmt.Errorf("unable to find section %s", sectionName)
	}
	return &sectionInfo{
		offset:   s.Offset,
		physAddr: s.Addr,
		size:     s.Size,
	}, nil
}

func getProgramSectionByHeader(kpath string, targExec bool, targWrite bool, targRead bool) (*sectionInfo, error) {
	elfFile, err := elf.Open(kpath)
	if err != nil {
		log.Fatalf("whoa!?!? can't read elf file but checked it before: %v", err)
	}
	defer elfFile.Close()
	for _, prog := range elfFile.Progs {
		if prog.Type == elf.PT_LOAD {
			isExecutable := false
			isReadable := false
			isWritable := false
			if prog.Flags&elf.PF_X == elf.PF_X {
				isExecutable = true
			}
			if prog.Flags&elf.PF_W == elf.PF_W {
				isWritable = true
			}
			if prog.Flags&elf.PF_R == elf.PF_R {
				isReadable = true
			}
			if targExec == isExecutable && targRead == isReadable && targWrite == isWritable {
				info := &sectionInfo{
					physAddr: prog.Paddr,
					offset:   prog.Off,
					size:     prog.Filesz,
				}
				return info, nil
			}
		}
	}
	return nil, fmt.Errorf("no executable program found in %s with attributes x=%v,w=%v,r==%v",
		kpath, targExec, targWrite, targRead)
}

func simpleTerminal(terminal io.ReadWriter, device string, byteLimit int) {
	t, ok := terminal.(*term.Term)
	if !ok {
		fatalf(terminal, "unable to run simple terminal when terminal is a file!")
	}
	byteCount := 0
	log.Printf("starting terminal loop....\n")
	//	log.Printf("hackery: %s",readLine(terminal))

	if device != "/dev/tty" {
		k := setupTTY("/dev/tty", true)
		kbd := k.(*term.Term)

		var wg sync.WaitGroup
		wg.Add(2)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				one := make([]byte, 1)
				_, err := kbd.Read(one)
				if err != nil {
					fatalf(terminal, "unable to read from /dev/tty: %v", err)
				}
				_, err = terminal.Write(one)
				if err != nil {
					fatalf(terminal, "unable to write to device: %v", err)
				}
				if one[0] == 0x04 {
					cleanupAndExit()
				}
			}
		}(&wg)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			one := make([]byte, 1)
			for {
				n, err := t.Read(one)
				if err != nil {
					fatalf(t, "unable to read from device: %v", err)
					return
				}
				if n == 0 {
					fmt.Printf("read failed (no error, but no data read)\n")
					continue
				}
				if one[0] == 0 {
					fmt.Printf("nul ")
					continue
				}
				if one[0] < 32 && one[0] != '\n' {
					fmt.Printf("[%02x]", one[0])
					continue
				}
				_, err = kbd.Write(one)
				if err != nil {
					fatalf(terminal, "unable to write to kbd terminal: %v", err)
				}
				byteCount++
				if byteLimit > 0 && byteCount >= byteLimit {
					log.Printf("byte limit of %d characters reached", byteLimit)
					cleanupAndExit()
				}
			}
		}(&wg)
		wg.Wait() //currently this will wait forever because there is no exit protocol
		cleanupAndExit()
	}
}

func pingLoop(t io.ReadWriter) {
	// ping is useful for being sure our connection is ok and giving us a chance
	// to futz with the set hardware flow control
	attempts := 0
	for attempts < maxPings {
		log.Printf("sending ping %d\n", attempts)
		ok, _ := sendSingleCommand(t, pingTransmit, "PING", 3, false)
		if ok {
			break
		}
		attempts++
		terminal, ok := t.(*term.Term) // just in case it ends up being a file or something
		if ok {
			if err := terminal.SetFlowControl(term.NONE); err != nil {
				log.Fatalf("unable to set flow control none:%v", err)
			}
		}
	}
	if attempts == maxPings {
		fatalf(t, "giving, cannot reach the device with ping...(%d attempts)", attempts)
	}
	log.Printf("connection established")
}
