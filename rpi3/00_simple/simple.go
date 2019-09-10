package main

import (
	dev "device/rpi3"
)

func main() { 
	print("echoing what you type\n")
	for {
		c := dev.UART0Getc()
		dev.UART0Send(c)
	}
}
