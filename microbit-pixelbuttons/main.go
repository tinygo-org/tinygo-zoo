// Example using the buttons and LED matrix on a BBC:Microbit
package main

import (
	"machine"
	"time"
)

func main() {
	machine.InitLEDMatrix()

	left := machine.GPIO{machine.BUTTONA}
	left.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT})

	right := machine.GPIO{machine.BUTTONB}
	right.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT})

	var (
		x uint8 = 2
		y uint8 = 2
	)

	machine.ClearLEDMatrix()

	for {
		if !left.Get() {
			switch {
			case x > 0:
				x--
			case x == 0:
				if y > 0 {
					x = 4
					y--
				}
			}
		}

		if !right.Get() {
			switch {
			case x < 4:
				x++
			case x == 4:
				if y < 4 {
					x = 0
					y++
				}
			}
		}

		machine.SetLEDMatrix(x, y)
		time.Sleep(time.Millisecond * 100)
	}
}
