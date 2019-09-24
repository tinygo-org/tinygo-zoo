#
# Use the bootloader to test our delays.
#

all: clean kernel8.elf

# note the target is rpi3_bl in this directory
kernel8.elf: rpi3_bl.o
	tinygo build -target 'rpi3_bl.json' -ldflags 'rpi3_bl.o' -no-debug  -o kernel8.elf .

rpi3_bl.o: rpi3_bl.S
	clang --target=aarch64-elf -c rpi3_bl.S

clean:
	rm kernel8.elf kernel8.img *.o >/dev/null 2>/dev/null || true
