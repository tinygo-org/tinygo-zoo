// blink program for the BBC micro:bit that uses the LED matrix
package main

import (
	"machine"
	"time"
)

func main() {
	machine.InitLEDMatrix()

	for {
		machine.ClearLEDMatrix()
		time.Sleep(time.Millisecond * 500)

		machine.SetLEDMatrix(2, 2)
		time.Sleep(time.Millisecond * 500)
	}
}
