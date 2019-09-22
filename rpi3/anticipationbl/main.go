package main

import (
	dev "device/rpi3"
	"unsafe"
)

//first character of line sent:
//? in command mode
//# in hex mode
//. ok
//! error

const ready = "?\n"
const readyHex = "#\n"
const sectionCommand = ":section"
const fetchCommand = ":fetch"
const inflateCommand = ":inflate"
const launchCommand = ":launch"
const pingCommand = ":ping"
const confirm = ".ok\n"
const fetchSize = 0x20

var buffer [1025]byte // : is first so to get 512 bytes, we add 1
var converted [512]uint8
var hexArgs [4]uint64
var fetchBuffer [0x40]byte

type blstate int

const (
	waiting blstate = iota
	section
)

type hexLineType int

const (
	dataLine hexLineType = iota
	endOfFile
	extendedSegmentAddress
	badBufferType
)

var tests = 0

func main() {
	state := waiting
	for {
		switch state {
		case waiting:
			length, ok := readLine()
			if !ok {
				print(ready)
				continue
			}
			if checkLine(sectionCommand, length) {
				print(confirm)
				state = section
				continue
			}
			if checkLine(pingCommand, length) {
				print(confirm)
				continue
			}

			if checkLine(fetchCommand, 0) {
				performFetch(length)
				continue
			}
			if checkLine(inflateCommand, 0) {
				performInflate(length)
				continue
			}
			if checkLine(launchCommand, 0) {
				performLaunch(length)
				continue
			}
			//s := string(buffer[:length])
			//print("! bad line:", s, "\n")
		case section:
			readHex() //succeed or fail, we go back to ready for more sections
			state = waiting
		}
	}
	//dev.QEMUTryExit()
}

func checkLine(s string, l int) bool {
	if l != 0 && len(s) != l {
		return false
	}
	for i, b := range []byte(s) {
		if buffer[i] != b {
			return false
		}
	}
	return true
}

// returns second value true if this timed out
// otherwise, first return is length of line in buffer, with L/F removed
func readLine() (int, bool) {
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 0
	}
	current := 0
	miss := 0
	for miss < 5000000 {
		if dev.UART0DataAvailable() {
			miss = 0
			b := dev.UART0Getc()
			if b == 10 {
				if current == 0 {
					print("!empty line! ")
					dev.Abort()
				}
				return current, true
			}
			buffer[current] = b
			current++
			continue
		}
		miss++
	}
	return -181711, false
}

const failLimitOnHexMode = 10

// read a file using Intel Hex format https://en.wikipedia.org/wiki/Intel_HEX
// only uses record types 00, 01, and 02
//go:export readHex
func readHex() bool {
	fails := 0
	baseAddr := uint64(0)

	//this loop reads a purported line of hex protocol and either complains or executes
	//the appropriate action
	for fails < failLimitOnHexMode {
		l, ok := readLine()
		if !ok {
			print(readyHex)
			fails++
			continue
		}
		if buffer[0] != ':' {
			print("!fail line not started with colon but was ", buffer[0], ",", buffer[1], ",", buffer[2], " and ", l, "\n")
			fails++
			continue
		}
		if !convertBuffer(l) {
			fails++
			continue
		}
		if !checkBufferLength(l) {
			fails++
			continue
		}
		if !checkChecksum(l) {
			fails++
			continue
		}
		t := lineType()
		if t == badBufferType {
			fails++
			continue
		}
		err, done := processLine(t, &baseAddr)
		if err {
			fails++
			continue
		}
		//line was ok if we get here
		fails = 0
		print(confirm)
		//bail out because we got EOF?
		if done {
			return true
		}
	}
	//too many fails
	return false

}

// deal with a received hex line and return (error?,done?)
func processLine(t hexLineType, baseAddr *uint64) (bool, bool) {
	switch t {
	case dataLine:
		l := converted[0]
		offset := (uint64(converted[1]) * 256) + (uint64(converted[2]))
		offset += *baseAddr
		var addr *uint8
		var val uint8
		for i := uint8(0); i < l; i++ {
			addr = (*uint8)(unsafe.Pointer(uintptr(offset) + uintptr(i)))
			val = converted[4+i]
			*addr = val
		}
		return false, false
	case endOfFile:
		return false, true
	case extendedSegmentAddress:
		len := converted[0]
		if len != 2 {
			print("!ESA value has too many bytes:", len, "\n")
			return true, false
		}
		esaAddr := uint64(converted[4])*256 + uint64(converted[5])
		esaAddr = esaAddr << 4 //it's assumed to be a multiple of 16
		if esaAddr == 0x80000 {
			print("!ESA value ,", esaAddr, "would load code over the bootloader\n")
			return true, false
		}
		*baseAddr = esaAddr
		return false, false
	}
	print("!internal error, unexpected line type", t, "\n")
	return true, false
}

// received a line, check that it has a hope of being syntactically correct
func checkBufferLength(l int) bool {
	total := uint8(11) //size of just framing in characters (colon, 2 len chars, 4 addr chars, 2 type chars, 2 checksum chars)
	if uint8(l) < total {
		print("!bad buffer length, can't be smaller than", total, ":", l, "\n")
		return false
	}
	total += converted[0] * 2
	if uint8(l) != total {
		print("!bad buffer length, expected", total, "but got", l, " based on ", converted[0], "\n")
		return false
	}
	return true
}

// verify line's checksum
func checkChecksum(l int) bool {
	sum := uint64(0)
	limit := (l - 1) / 2
	for i := 0; i < limit; i++ {
		sum += uint64(converted[i])
	}
	complement := ^sum
	complement++
	checksum := uint8(complement & 0xff)
	if checksum != 0 {
		print("!bad checksum, expected ", checksum, "but got ", converted[limit-1], "\n")
		return false
	}
	return true
}

// extract the line type, 00 (data), 01 (eof), or 02 (esa)
func lineType() hexLineType {
	switch converted[3] {
	case 0:
		return dataLine
	case 1:
		return endOfFile
	case 2:
		return extendedSegmentAddress
	default:
		print("!bad buffer type:", converted[3], "\n")
		return badBufferType
	}
}

// change buffer of ascii->converted bytes by taking the ascii values (2 per byte) and making them proper bytes
func convertBuffer(l int) bool {
	//l-1 because the : is skipped so the remaining number of characters must be even
	if (l-1)%2 == 1 {
		print("!bad payload, expected even number of hex bytes (length read minus LF is=", l, ")")
		for i := 0; i < l; i++ {
			print(i, "=", buffer[i], " ")
		}
		print("\n")
		return false
	}
	//skip first colon
	for i := 1; i < l; i += 2 {
		v, ok := bufferValue(i)
		if !ok {
			return false // they already sent the error to the other side
		}
		converted[(i-1)/2] = v
	}
	return true
}

func performInflate(length int) {
	start := len(inflateCommand + " ")
	if !splitHexArguments(2, start, length, "inflate") {
		return
	}
	addr := hexArgs[0]
	size := hexArgs[1]
	for i := int(addr); i < int(addr)+int(size); i++ {
		ptr := (*uint8)((unsafe.Pointer)((uintptr)(i)))
		*ptr = 0
	}
	print(confirm)
}

func performLaunch(length int) {
	start := len(launchCommand + " ")
	if !splitHexArguments(4, start, length, "launch") {
		return
	}
	print(confirm)
	addr := hexArgs[0]
	stackPtr := hexArgs[1] - 0x10 //writes are from lower to higher, but we have be 16 byte aligned
	heapPtr := hexArgs[2]
	now := hexArgs[3]
	//print("xxx about to jump:", addr, " ", stackPtr, " ", heapPtr, " ", now, "\n")
	dev.AsmFull(`mov x1,{stackPtr}
		mov x2,{heapPtr}
		mov x3,{now}
		mov x0,{addr}
		br x0`, map[string]interface{}{"addr": addr, "stackPtr": stackPtr, "heapPtr": heapPtr, "now": now})

	print("!bad launch addr=", addr, " sp=", stackPtr, " heap=", heapPtr, " time=", now, "\n")
}

func performFetch(length int) {
	start := len(fetchCommand + " ")
	if !splitHexArguments(1, start, length, "fetch") {
		return
	}
	for i := int(hexArgs[0]); i < int(hexArgs[0])+fetchSize; i++ {
		index := i - int(hexArgs[0])
		ptr := (*uint8)((unsafe.Pointer)((uintptr)(i)))
		thisByte := *ptr
		hi := thisByte >> 4
		lo := thisByte & 0xf
		for j, v := range []uint8{hi, lo} {
			switch v {
			case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
				fetchBuffer[(index*2)+j] = byte(48 + v)
			case 10, 11, 12, 13, 14, 15:
				fetchBuffer[(index*2)+j] = byte(65 + (v - 10))
			default:
				print("!bad value in hex?!? ", v, "\n")
			}
		}
	}
	print(".ok ", string(fetchBuffer[:]), "\n")
}

// this hits buffer[i] and buffer[i+1] to convert an ascii byte
// returns false to mean you had a bad character in the input
func bufferValue(i int) (uint8, bool) {
	total := uint8(0)
	switch buffer[i] {
	case '0':
	case '1':
		total += 16 * 1
	case '2':
		total += 16 * 2
	case '3':
		total += 16 * 3
	case '4':
		total += 16 * 4
	case '5':
		total += 16 * 5
	case '6':
		total += 16 * 6
	case '7':
		total += 16 * 7
	case '8':
		total += 16 * 8
	case '9':
		total += 16 * 9
	case 'a', 'A':
		total += 16 * 10
	case 'b', 'B':
		total += 16 * 11
	case 'c', 'C':
		total += 16 * 12
	case 'd', 'D':
		total += 16 * 13
	case 'e', 'E':
		total += 16 * 14
	case 'f', 'F':
		total += 16 * 15
	default:
		print("!bad character in payload hi byte(number #", i, "):", buffer[i], "\n")
		return 0xff, false
	}
	switch buffer[i+1] {
	case '0':
	case '1':
		total++
	case '2':
		total += 2
	case '3':
		total += 3
	case '4':
		total += 4
	case '5':
		total += 5
	case '6':
		total += 6
	case '7':
		total += 7
	case '8':
		total += 8
	case '9':
		total += 9
	case 'a', 'A':
		total += 10
	case 'b', 'B':
		total += 11
	case 'c', 'C':
		total += 12
	case 'd', 'D':
		total += 13
	case 'e', 'E':
		total += 14
	case 'f', 'F':
		total += 15
	default:
		print("!bad character in payload low byte (number #", i+1, "):", buffer[i+1], "\n")
		return 0xff, false
	}
	return total, true
}

func bufferASCIIToUint(start int, end int) (uint64, bool) {
	var ok bool
	total := uint64(0)
	var thisByte uint8
	placeValue := uint64(24)
	for i := start; i < end; i += 2 {
		thisByte, ok = bufferValue(i)
		if !ok {
			return 0, false
		}
		total += uint64(thisByte) * (1 << placeValue)
		placeValue -= 8
	}
	return total, true
}

func splitHexArguments(expected int, start int, length int, name string) bool {
	for argNum := 0; argNum < expected; argNum++ {
		if argNum == expected-1 { //last arg
			end := length
			if (end-start)%2 != 0 {
				print("!", name, ":argument ", argNum, " [last] has odd number of hex digits, distance=", (end - start), "(length of total command was ", length, " and start was ", start, ") \n")
				return false
			}
			total, ok := bufferASCIIToUint(start, end)
			if !ok {
				return false
			}
			hexArgs[argNum] = total
			return true
		}
		//we are not at last argument, look for next space
		found := false
		end := -1
		for i := start; i < length; i++ {
			if buffer[i] == 0x20 {
				end = i
				found = true
				break
			}
		}
		if !found {
			print("!unable to find end of argument number ", argNum, "\n")
			return false
		}
		if (end-start)%2 != 0 {
			print("!argument ", argNum, "has odd number of hex digits:", (end - start), "(length of total command was ", length, ", start was ", start, " and end was ", end, ") \n")
			return false
		}

		total, ok := bufferASCIIToUint(start, end)
		if !ok {
			return false
		}
		hexArgs[argNum] = total
		start = end + 1
	}
	return true
}
