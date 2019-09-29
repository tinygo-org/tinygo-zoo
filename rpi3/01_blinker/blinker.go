package main

import (
	dev "device/rpi3"
)

var on = false
var interval uint32

// it's more than a little unclear that the values passed to this
// interrupt handler are ever going to be useful
func blinker(_ uint32, _ uint64) uint32 {
	dev.LEDSet(on)
	on = !on
	print("set LED to ", on, "\n")
	return interval
}

func main() {
	hi, lo := dev.GetRPIID()
	if hi == 0 && lo == 0 {
		print("you are running on QEMU so you aren't going to see any lights blinking...\n")
	}

	freq := dev.CounterFreq() //getting my freq on
	print("frequency per second is:")
	dev.UART0Hex(freq) //this number is number of ticks/sec

	//sets the timer for  N secs, where N is multiplier
	interval = 1 * freq
	dev.SetCounterTargetInterval(interval, blinker)
	dev.Core0CounterToCore0Irq()
	dev.EnableCounter()
	dev.EnableTimerIRQ()
	for {
		dev.WaitForInterrupt()
	}
}
