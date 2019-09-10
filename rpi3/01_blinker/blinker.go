package main

import (
	dev "device/rpi3"
)

var on = false
var interval uint32

// it's more than a little unclear that the values passed to this
// interrupt handler are ever going to be useful
func blinker(foo uint32, bar uint64) uint32 {
	dev.LEDSet(on)
	on = !on
	return interval
}

//
// MAIN_QEMU hits the QEMU version of the counters directly.  It won't work
// on real hardware.
//
func main() {

	dev.GetRPIID()
	blinkerSetup()

	freq := dev.QEMUCounterFreq() //getting my freq on
	print("frequency per second is:")
	dev.UART0Hex(freq) //this number is number of ticks/sec

	//sets the timer for  N secs, where N is multiplier
	interval = 1 * freq
	dev.QEMUSetCounterTargetInterval(interval, blinker)
	dev.QEMUCore0CounterToCore0Irq()
	dev.QEMUEnableCounter()
	dev.EnableTimerIRQ()
	for {
		dev.WaitForInterrupt()
	}
}
