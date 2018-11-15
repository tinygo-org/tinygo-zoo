# TinyGo Zoo

Various sample programs for microcontrollers using TinyGo (http://tinygo.org)

## Installation requirements

### Install Docker image

    docker pull tinygo/tinygo

### Install flashing tools for BBC micro:bit

    sudo apt-get install openocd

### Install flashing tools for Arduino

    sudo apt-get install avrdude

### Install flashing tools for PCA10040

Install nrfjprog.

## Blink for BBC micro:bit

    make microbit-blink

## Pixel buttons for BBC micro:bit

    make microbit-pixelbuttons

## Color lamp for Arduino Uno

    make arduino-colorlamp

## Color lamp for PCA10040 (NRF52-DK)

    make nrf-colorlamp
