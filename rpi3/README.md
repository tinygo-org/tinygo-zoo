# Raspberry PI Samples

## Requirements
You have to have three things in your PATH.  A copy of tinygo (that includes
the rpi3 devices, version 0.9+), a "normal" copy of go of at least version 1.11,
and `llvm-copy`.

The "normal" go will be used (only) to compile/run a program that runs on the host
computer and talks to  the RPI3 over serial.

`llvm-objcopy` is used to extract a bootable image from an elf file created by
tinygo, which is using llvm under the covers.  If you built llvm as part
of installing tinygo, then you can probably find the binary you need in
`llvm-build/bin/llvm-objcopy` within your tinygo source tree.

## Hardware vs. Emulation

It may be advantageous to have a recent copy of QEMU (4.10+) so you can run
in emulation mode without needing an actual piece of hardware.

If you are using a real hardware version of the RPI3, it must be connected to the
host computer over a serial port. See this tutorial for how to install a
serial cable connection to the host and RPI3:
https://learn.adafruit.com/adafruits-raspberry-pi-lesson-5-using-a-console-cable/overview

## How To Use

Run the "device" part of the bootloader `anticipationbl` either on your RPI3 (connected via
serial to the host) or on QEMU (see `make runqemu`).  Then you run the host side,
passing the appropriate device (see `make runqemu` in `anticipation` for an example
with QEMU).  You need to make sure you use the right device to talk to the device
side, either serial port like `/dev/tty.SLAB_USBtoUART` or the tty device created
by QEMU when you started the device side.  

Once you have started the host side, it will transfer the elf file you supplied
on the command line.  Once that is completed, the device-side bootloader jumps
to the downloaded file's entry point and the downloaded file starts executing.

## What's Here

* `anticipation` and `anticipationbl` which together allow bootloading over
the serial port on either hardware or on qemu.

* `00_simple` Simplest possible bootloaded program.  It prints out one line
to the terminal, then echos back whatever you type at it.  Works on QEMU.

* `01_simple` The inevitable blinking light example.  This one uses the timers
in QEMU or on the hardware (they aren't the same) and handles the timer
interrupt to blink the light.

* `02_delay` Demonstrates how to use the delays on the system and how to get
the system time.  Note that RPI3 does not have a battery-powered clock, so
the time is copied from the host to the running program via the bootloader.

* `03_framebuffer` Initializes the framebuffer, and then uses a font to write
to the "console" on the screen. It does this in a loop so it scrolls a lot.

## To bootload or not to bootload

The bootloader provided, anticipation, has two primary advantages over running
your own code on bare metal by creating a `kernel8.img` and then booting from that.

1) Less hassle when doing development.  Constantly changing SD cards from the
RPI3 to the host and back to update code is irritating. Further, this increases
the mechanical wear and tear on the RPI3's card slot--which is less than robust.

2) The bootloader can get your system into a known, useful state without you
having to worry about initialization.  Primarily this means that devices are
initialized, your code is running with a sensible stack and heap pointers, the
time is set, and so on.

### If you don't want to bootload

Copy the makefile from `anticipationbl` and go for it! You can see there that
you use tinygo to create a self contained kernel image in the same way that
the bootloader builds.
