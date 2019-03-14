# TinyGo Zoo

Various sample programs for microcontrollers using TinyGo (http://tinygo.org)

## Installation requirements

### Install Docker image

    docker pull tinygo/tinygo

### Install tinygo-drivers into your LOCAL golang installation

    go get -d github.com/tinygo-org/drivers

### Install flashing tools

- BBC micro:bit

        sudo apt-get install openocd

- Arduino

        sudo apt-get install avrdude

- PCA10040
Install nrfjprog as follows.
        
First install the J-Link Software and Documentation Pack from Segger:
https://www.segger.com/downloads/jlink/#J-LinkSoftwareAndDocumentationPack

Then install the nRF5x Command-Line Tools:
https://docs.zephyrproject.org/latest/tools/nordic_segger.html#nordic-segger


## Blink for BBC micro:bit

    make microbit-blink

## Pixel buttons for BBC micro:bit

    make microbit-pixelbuttons

## Color lamp for Arduino Uno

    make arduino-colorlamp

## Color lamp for PCA10040 (NRF52-DK)

    make nrf-colorlamp

## Built-in accelerometer on BBC micro:bit

    make microbit-accel

## Built-in accelerometer on reel board

    make reelboard-accel
