package main

import (
	dev "device/rpi3"
)

func main() {
	print("starting at\n")
	dev.UART0TimeDateString(dev.Now())
	print("\n")

	if !dev.InitFramebuffer() {
		print("unable to init framebuffer!")
	}
	print("done\n")
}
