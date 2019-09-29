all: kernel8.img

# note that we are building kernel8.img as well as kernel8.elf.
# on raspberry pi (which uses bcm2837) kernel8.img is one that can be copied
# onto an SD card and booted.  The name of the kernel to load is kernel8.img on a RPI3 but
# if you have changed options in config.txt on the SD card, it might different
# for you.
#
# kernel8.elf is used if you want to run and use QEMU as a simulator (-M raspi3)
#

## NOTE: this uses the bare metal version of RPI (--target rpi3) but all the things
## NOTE: loaded WITH this bootloader should supply the --target rpi3_bl
kernel8.img: clean
	tinygo build --target=rpi3 -o kernel8.elf .
	llvm-objcopy -O binary kernel8.elf kernel8.img

clean:
	rm kernel8.elf kernel8.img *.o >/dev/null 2>/dev/null || true

## for running the simulator, use this target BEFORE starting the host-side
## in ../anticipation
runqemu:
	qemu-system-aarch64 --chardev 'pty,id=char0' -M raspi3 -kernel kernel8.img -serial chardev:char0
