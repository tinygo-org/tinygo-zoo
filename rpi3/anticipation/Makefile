all:
	go build -o anticipation ./cmd/anticipation

## QEMU info
## 1) you have to use the right TTY that is created by anticipationbl when
## it starts (via QEMU).  You should run anticipationbl first so you can see the output
## of the device (tty) to use to connect to it.
## 2) You need to choose what you want the bootloader to load.  In the
## example below, we are loading the kernel (must be an ELF file) created
## by building in the 00_simple directory.
runqemu:
	./anticipation -d /dev/ttys007  ../00_simple/kernel8.elf

## RPI3 (hardware) info
## 1) you need to know where your Serial-To-USB device is.  I have
## listed the device below as it exists on my Mac.
## 2) You need to choose what you want the bootloader to load.  In the
## example below, we are loading the kernel (must be an ELF file) created
## by building in the 00_simple directory.
runrpi3:
	./anticipation -d /dev/tty.SLAB_USBtoUART ../00_simple/kernel8.elf
