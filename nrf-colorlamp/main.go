// This program runs on an PCA10040 that has the following four devices connected:
// - Button connected to P0.11
// - Rotary analog dial connected to P0.03 (A0)
// - RGB LED connected to P0.23, P0.24, and P0.25 used as PWM pins
// - BlinkM I2C RGB LED
//
// Pushing the button switches which color is selected.
// Rotating the dial changes the value for the currently selected color.
// Changing the color value updates the color displayed on both the
// PWM-controlled RGB LED and the I2C-controlled BlinkM.
package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/blinkm"
)

const (
	buttonPin = 11
	redPin    = 24
	greenPin  = 25
	bluePin   = 23

	red   = 0
	green = 1
	blue  = 2
)

func main() {
	machine.InitADC()
	machine.InitPWM()
	machine.I2C0.Configure(machine.I2CConfig{})

	// Init BlinkM
	blm := blinkm.New(machine.I2C0)
	blm.StopScript()

	button := machine.Pin(buttonPin)
	button.Configure(machine.PinConfig{Mode: machine.PinInput})

	dial := machine.ADC{machine.ADC0}
	dial.Configure()

	redLED := machine.PWM{redPin}
	redLED.Configure()

	greenLED := machine.PWM{greenPin}
	greenLED.Configure()

	blueLED := machine.PWM{bluePin}
	blueLED.Configure()

	selectedColor := red
	colors := []uint16{0, 0, 0}

	for {
		// If we pushed the button, switch active color.
		if !button.Get() {
			if selectedColor == blue {
				selectedColor = red
			} else {
				selectedColor++
			}
		}

		// Change the intensity for the currently selected color based on the dial setting.
		colors[selectedColor] = (dial.Get())

		// Update the RGB LED.
		redLED.Set(colors[red])
		greenLED.Set(colors[green])
		blueLED.Set(colors[blue])

		// Update the BlinkM.
		blm.SetRGB(byte(colors[red]>>8), byte(colors[green]>>8), byte(colors[blue]>>8))

		time.Sleep(time.Millisecond * 100)
	}
}
