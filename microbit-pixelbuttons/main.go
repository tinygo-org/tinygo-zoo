// Example using the buttons and LED matrix on a BBC:Microbit
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/microbitmatrix"
)

func main() {
	display := microbitmatrix.New()
	display.Configure(microbitmatrix.Config{})
	display.ClearDisplay()

	left := machine.BUTTONA
	left.Configure(machine.PinConfig{Mode: machine.PinInput})

	right := machine.BUTTONB
	right.Configure(machine.PinConfig{Mode: machine.PinInput})

	var (
		x int16 = 2
		y int16 = 2
		c       = color.RGBA{255, 255, 255, 255}
	)

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

		display.SetPixel(x, y, c)
		time.Sleep(time.Millisecond * 100)
	}
}
