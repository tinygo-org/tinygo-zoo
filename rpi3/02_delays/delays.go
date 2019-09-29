package main

import (
	dev "device/rpi3"
)

var interval uint32

// not sure how useful the current value is, but the virtual timer should
// have increased by exactly interval (or very close to it)
// return the amount of time to wait for the next call to this callback
// return 0 if you don't want it anymore
func myHandler(current uint32, virtualTimer uint64) uint32 {
	print("my handler, current: ")
	dev.UART0Hex(current)
	print("my handler, virtual timer: ")
	dev.UART0Hex64(virtualTimer)
	return interval
}

func main() {
	dev.UART0TimeDateString(dev.Now())
	print("\n")
	testDelays()
	// dev.BComTimerInit(interval)
	// dev.EnableInterruptController()
	// dev.EnableTimerIRQ()
	//
	// for {
	// 	dev.WaitForInterrupt()
	// }
	dev.Abort()
}

func testDelays() {
	print("Waiting 3000000 CPU cycles (ARM CPU): ")
	dev.WaitCycles(3000000)
	print("OK\n")

	print("Waiting 3000000 microsec (ARM CPU): ")
	dev.WaitMuSec(3000000)
	print("OK\n")

	print("Waiting 3000000 microsec (BCM System Timer): ")
	if dev.SysTimer() == 0 {
		print("Not available\n")
	} else {
		dev.WaitMuSecST(3000000)
		print("OK\n")
	}
}
