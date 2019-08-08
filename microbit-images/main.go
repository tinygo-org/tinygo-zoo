package main

import (
	"image/color"
	"time"

	"tinygo.org/x/drivers/microbitmatrix"
)

var display microbitmatrix.Device

// Reuse from https://github.com/bbcmicrobit/micropython/blob/master/source/microbit/microbitconstimage.cpp
var ImageHeart = [25]int16{
	0, 1, 0, 1, 0,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageHeartSmall = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
}

// smilies

var ImageHappy = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	1, 0, 0, 0, 1,
	0, 1, 1, 1, 0,
}

var ImageSmile = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	1, 0, 0, 0, 1,
	0, 1, 1, 1, 0,
}

var ImageSad = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	0, 1, 1, 1, 0,
	1, 0, 0, 0, 1,
}

var ImageConfused = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	1, 0, 1, 0, 1,
}

var ImageAngry = [25]int16{
	1, 0, 0, 0, 1,
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	1, 1, 1, 1, 1,
	1, 0, 1, 0, 1,
}

var ImageAsleep = [25]int16{
	0, 0, 0, 0, 0,
	1, 1, 0, 1, 1,
	0, 0, 0, 0, 0,
	0, 1, 1, 1, 0,
	0, 0, 0, 0, 0,
}

var ImageSurprised = [25]int16{
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageSilly = [25]int16{
	1, 0, 0, 0, 1,
	0, 0, 0, 0, 0,
	1, 1, 1, 1, 1,
	0, 0, 1, 0, 1,
	0, 0, 1, 1, 1,
}

var ImageFabulous = [25]int16{
	1, 1, 1, 1, 1,
	1, 1, 0, 1, 1,
	0, 0, 0, 0, 0,
	0, 1, 0, 1, 0,
	0, 1, 1, 1, 0,
}

var ImageMeh = [25]int16{
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 1, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 0, 0,
}

// yes/no

var ImageYes = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 1,
	0, 0, 0, 1, 0,
	1, 0, 1, 0, 0,
	0, 1, 0, 0, 0,
}

var ImageNo = [25]int16{
	1, 0, 0, 0, 1,
	0, 1, 0, 1, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	1, 0, 0, 0, 1,
}

// clock hands

var ImageClock12 = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock1 = [25]int16{
	0, 0, 0, 1, 0,
	0, 0, 0, 1, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock2 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 1, 1,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock3 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 1, 1,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock4 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 1, 1,
	0, 0, 0, 0, 0,
}

var ImageClock5 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 1, 0,
	0, 0, 0, 1, 0,
}

var ImageClock6 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
}

var ImageClock7 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 0, 0,
	0, 1, 0, 0, 0,
}

var ImageClock8 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	1, 1, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock9 = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	1, 1, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock10 = [25]int16{
	0, 0, 0, 0, 0,
	1, 1, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageClock11 = [25]int16{
	0, 1, 0, 0, 0,
	0, 1, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

// arrows

var ImageArrowN = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	1, 0, 1, 0, 1,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
}

var ImageArrowNE = [25]int16{
	0, 0, 1, 1, 1,
	0, 0, 0, 1, 1,
	0, 0, 1, 0, 1,
	0, 1, 0, 0, 0,
	1, 0, 0, 0, 0,
}

var ImageArrowE = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 0, 1, 0,
	1, 1, 1, 1, 1,
	0, 0, 0, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageArrowSE = [25]int16{
	1, 0, 0, 0, 0,
	0, 1, 0, 0, 0,
	0, 0, 1, 0, 1,
	0, 0, 0, 1, 1,
	0, 0, 1, 1, 1,
}

var ImageArrowS = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	1, 0, 1, 0, 1,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageArrowSW = [25]int16{
	0, 0, 0, 0, 1,
	0, 0, 0, 1, 0,
	1, 0, 1, 0, 0,
	1, 1, 0, 0, 0,
	1, 1, 1, 0, 0,
}

var ImageArrowW = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 0, 0, 0,
	1, 1, 1, 1, 1,
	0, 1, 0, 0, 0,
	0, 0, 1, 0, 0,
}

var ImageArrowNW = [25]int16{
	1, 1, 1, 0, 0,
	1, 1, 0, 0, 0,
	1, 0, 1, 0, 0,
	0, 0, 0, 1, 0,
	0, 0, 0, 0, 1,
}

// geometry

var ImageTriangle = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	1, 1, 1, 1, 1,
	0, 0, 0, 0, 0,
}

var ImageTriangleLeft = [25]int16{
	1, 0, 0, 0, 0,
	1, 1, 0, 0, 0,
	1, 0, 1, 0, 0,
	1, 0, 0, 1, 0,
	1, 1, 1, 1, 1,
}

var ImageChessboard = [25]int16{
	0, 1, 0, 1, 0,
	1, 0, 1, 0, 1,
	0, 1, 0, 1, 0,
	1, 0, 1, 0, 1,
	0, 1, 0, 1, 0,
}

var ImageDiamond = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	1, 0, 0, 0, 1,
	0, 1, 0, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageDiamondSmall = [25]int16{
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 0, 0,
}

var ImageSquare = [25]int16{
	1, 1, 1, 1, 1,
	1, 0, 0, 0, 1,
	1, 0, 0, 0, 1,
	1, 0, 0, 0, 1,
	1, 1, 1, 1, 1,
}

var ImageSquareSmall = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 1, 1, 0,
	0, 1, 0, 1, 0,
	0, 1, 1, 1, 0,
	0, 0, 0, 0, 0,
}

// animals

var ImageRabbit = [25]int16{
	1, 0, 1, 0, 0,
	1, 0, 1, 0, 0,
	1, 1, 1, 1, 0,
	1, 1, 0, 1, 0,
	1, 1, 1, 1, 0,
}

var ImageCow = [25]int16{
	1, 0, 0, 0, 1,
	1, 0, 0, 0, 1,
	1, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
}

// musical notes

var ImageMusicCrotchet = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	1, 1, 1, 0, 0,
	1, 1, 1, 0, 0,
}

var ImageMusicQuaver = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 1, 1, 0,
	0, 0, 1, 0, 1,
	1, 1, 1, 0, 0,
	1, 1, 1, 0, 0,
}

var ImageMusicQuavers = [25]int16{
	0, 1, 1, 1, 1,
	0, 1, 0, 0, 1,
	0, 1, 0, 0, 1,
	1, 1, 0, 1, 1,
	1, 1, 0, 1, 1,
}

// other icons

var ImagePitchfork = [25]int16{
	1, 0, 1, 0, 1,
	1, 0, 1, 0, 1,
	1, 1, 1, 1, 1,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
}

var ImageXmas = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	1, 1, 1, 1, 1,
}

var ImagePacman = [25]int16{
	0, 1, 1, 1, 1,
	1, 1, 0, 1, 0,
	1, 1, 1, 0, 0,
	1, 1, 1, 1, 0,
	0, 1, 1, 1, 1,
}

var ImageTarget = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	1, 1, 0, 1, 1,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
}

/*
The following images were designed by Abbie Brooks.
*/

var ImageTshirt = [25]int16{
	1, 1, 0, 1, 1,
	1, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 1, 1, 1, 0,
	0, 1, 1, 1, 0,
}

var ImageRollerskate = [25]int16{
	0, 0, 0, 1, 1,
	0, 0, 0, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	0, 1, 0, 1, 0,
}

var ImageDuck = [25]int16{
	0, 1, 1, 0, 0,
	1, 1, 1, 0, 0,
	0, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 0, 0, 0, 0,
}

var ImageHouse = [25]int16{
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	1, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 1, 0, 1, 0,
}

var ImageTortoise = [25]int16{
	0, 0, 0, 0, 0,
	0, 1, 1, 1, 0,
	1, 1, 1, 1, 1,
	0, 1, 0, 1, 0,
	0, 0, 0, 0, 0,
}

var ImageButterfly = [25]int16{
	1, 1, 0, 1, 1,
	1, 1, 1, 1, 1,
	0, 0, 1, 0, 0,
	1, 1, 1, 1, 1,
	1, 1, 0, 1, 1,
}

var ImageStickfigure = [25]int16{
	0, 0, 1, 0, 0,
	1, 1, 1, 1, 1,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	1, 0, 0, 0, 1,
}

var ImageGhost = [25]int16{
	1, 1, 1, 1, 1,
	1, 0, 1, 0, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 0, 1, 0, 1,
}

var ImageSword = [25]int16{
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 1, 0, 0,
	0, 1, 1, 1, 0,
	0, 0, 1, 0, 0,
}

var ImageGiraffe = [25]int16{
	1, 1, 0, 0, 0,
	0, 1, 0, 0, 0,
	0, 1, 0, 0, 0,
	0, 1, 1, 1, 0,
	0, 1, 0, 1, 0,
}

var ImageSkull = [25]int16{
	0, 1, 1, 1, 0,
	1, 0, 1, 0, 1,
	1, 1, 1, 1, 1,
	0, 1, 1, 1, 0,
	0, 1, 1, 1, 0,
}

var ImageUmbrella = [25]int16{
	0, 1, 1, 1, 0,
	1, 1, 1, 1, 1,
	0, 0, 1, 0, 0,
	1, 0, 1, 0, 0,
	0, 1, 1, 0, 0,
}

var ImageSnake = [25]int16{
	1, 1, 0, 0, 0,
	1, 1, 0, 1, 1,
	0, 1, 0, 1, 0,
	0, 1, 1, 1, 0,
	0, 0, 0, 0, 0,
}

// CreateImage function to get display device
func CreateImage(display *microbitmatrix.Device, image [25]int16) {
	c := color.RGBA{255, 255, 255, 255}
	w, _ := display.Size()
	// Pins Down
	display.SetRotation(1)
	for index, value := range image {
		x := (index % int(w))
		y := (index / int(w))
		if value == 1 {
			display.SetPixel(int16(x), int16(y), c)
		}
	}
}

func main() {
	display = microbitmatrix.New()
	display.Configure(microbitmatrix.Config{})
	then := time.Now()

	faces := [2][25]int16{ImageSad, ImageHappy}
	f := 0
	for {
		if time.Since(then).Nanoseconds() > 800000000 {
			display.ClearDisplay()
			CreateImage(&display, faces[f])
			f = (f + 1) % 2
			then = time.Now()
		}
		display.Display()
	}
}
