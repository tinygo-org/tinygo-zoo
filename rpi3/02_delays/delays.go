package main

import (
	dev "device/rpi3"
	"runtime/volatile"
)


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

func mainForRPI3OnHardware() {
	//dev.IRQVectorInitEL1(HandleIRQ)
	dev.TimerInit(interval)
	dev.EnableInterruptController()
	dev.EnableTimerIRQ()
}

func testDelays() {
	print("hello from main\n")
	print("Waiting 1000000 CPU cycles (ARM CPU): ")
	dev.WaitCycles(1000000)
	print("OK\n")

	print("Waiting 1000000 microsec (ARM CPU): ")
	dev.WaitMuSec(1000000)
	print("OK\n")

	print("Waiting 1000000 microsec (BCM System Timer): ")
	if dev.SysTimer() == 0 {
		print("Not available\n")
	} else {
		dev.WaitMuSecST(1000000)
		print("OK\n")
	}
}
