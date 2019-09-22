package main

import (
	dev "device/rpi3"
	"unsafe"
)

func main() {
	print("starting at\n")
	dev.UART0TimeDateString(dev.Now())
	print("\n")

	if !dev.InitFramebuffer() {
		print("unable to init framebuffer!")
	}
	font := (*psf)(unsafe.Pointer(&_binary_font_psf_start))
	print("magic=", font.magic, " headersize=", font.headersize, " height=", font.height, " width=", font.width, "bytesperglyph=", font.bytesperglyph, " \n")
	print("done\n")
}

//go:extern _binary_font_psf_start
var _binary_font_psf_start *uint8

type psf struct {
	magic         uint32
	version       uint32
	headersize    uint32
	flags         uint32
	numglyph      uint32
	bytesperglyph uint32
	height        uint32
	width         uint32
	glyphs        *uint8
}
