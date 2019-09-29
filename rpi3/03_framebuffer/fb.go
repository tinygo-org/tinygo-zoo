package main

import (
	dev "device/rpi3"
	"unsafe"
)

const width = 1024
const height = 768

var twoDigits = []byte{99, 99}

func main() {
	dev.UART0TimeDateString(dev.Now())
	print("\n")

	if !dev.InitFramebuffer() {
		print("unable to init framebuffer!")
		dev.Abort()
	}

	font := dev.NewPSFFontViaLinker(unsafe.Pointer(&_binary_font_psf_start), &dev.FrameBufferInfo)
	for {
		display48(font)
	}

	print("finished ok")
}

func display48(font *dev.PSFFont) {
	font.ConsolePrint("0 Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	font.ConsolePrint("1 Sed vehicula lacinia malesuada.")
	font.ConsolePrint("2 Phasellus sagittis nisl nisl, nec placerat lectus rutrum nec.")
	font.ConsolePrint("3")
	font.ConsolePrint("4")
	font.ConsolePrint("5")
	font.ConsolePrint("6 Donec sed nibh ut tortor finibus ultricies et non tellus.")
	font.ConsolePrint("7 Vivamus eget suscipit nibh.")
	font.ConsolePrint("8 Duis egestas, velit non hendrerit eleifend, libero neque bibendum nibh, bibendum cursus odio metus ac sapien.")
	font.ConsolePrint("9")
	font.ConsolePrint("0")
	font.ConsolePrint("1 Nullam enim turpis, egestas vitae mi vel, scelerisque interdum dolor.")
	font.ConsolePrint("2 Aenean vestibulum tortor vel congue pulvinar.")
	font.ConsolePrint("3 Suspendisse lobortis varius convallis.")
	font.ConsolePrint("4")
	font.ConsolePrint("5")
	font.ConsolePrint("6 Mauris quis consequat dui.")
	font.ConsolePrint("7 In sagittis elit at felis cursus, eget aliquam dui aliquam.")
	font.ConsolePrint("8 Curabitur augue ante, ullamcorper hendrerit nulla sit amet, sollicitudin euismod ante.")
	font.ConsolePrint("9")
	font.ConsolePrint("0")
	font.ConsolePrint("1")
	font.ConsolePrint("2 Nam varius ultricies condimentum.")
	font.ConsolePrint("3 Interdum et malesuada fames ac ante ipsum primis in faucibus.")
	font.ConsolePrint("4 Vivamus rhoncus laoreet molestie.")
	font.ConsolePrint("5")
	font.ConsolePrint("6")
	font.ConsolePrint("7 Nunc mattis nec elit at varius.")
	font.ConsolePrint("8 Aenean faucibus aliquam augue ac gravida.")
	font.ConsolePrint("9 Ut non tellus luctus, laoreet nulla eget, congue dolor.")
	font.ConsolePrint("0")
	font.ConsolePrint("1")
	font.ConsolePrint("2 Pellentesque fringilla tincidunt rutrum.")
	font.ConsolePrint("3 Duis ultricies auctor fringilla. Nunc lacus arcu, scelerisque ac arcu a, finibus viverra turpis.")
	font.ConsolePrint("4 Fusce auctor eleifend erat, sit amet iaculis augue.")
	font.ConsolePrint("5")
	font.ConsolePrint("6 Mauris orci quam, ornare eget orci ut, ullamcorper sollicitudin nunc.")
	font.ConsolePrint("7 Morbi eu nibh urna.")
	font.ConsolePrint("8 Sed vel aliquam nisl.")
	font.ConsolePrint("9 ")
	font.ConsolePrint("0 Mauris ornare tellus eu metus blandit congue.")
	font.ConsolePrint("1 Nunc metus lacus, laoreet finibus tempus quis, gravida non lectus.")
	font.ConsolePrint("2 Sed sed rutrum neque, sed venenatis libero.")
	font.ConsolePrint("3")
	font.ConsolePrint("4 Interdum et malesuada fames ac ante ipsum primis in faucibus.")
	font.ConsolePrint("5 Nam pretium tristique lectus. Vestibulum venenatis tellus ut euismod hendrerit.")
	font.ConsolePrint("6 Suspendisse quis nibh blandit, tincidunt mauris eget, eleifend lorem.")
	font.ConsolePrint("7 <------") //48th line
}

//go:extern _binary_font_psf_start
var _binary_font_psf_start *uint8

//go:export sync_el1h_handler
func interruptHandler(n int, esr uint64, address uint64) {
	//if you see this, probably you are hitting the GPU's mailbox with bad params
	//see README.md
	print("unexpected interrupt: synchronous el1h: esr:", esr, " address 0x")
	dev.UART0Hex64(address)
	sp := dev.ReadRegister("sp")
	print("sp is ")
	dev.UART0Hex64(uint64(sp))
	dev.Abort()
}
