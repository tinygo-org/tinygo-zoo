#
# Simplest bootloader example.
#

all: clean kernel8.elf

# note the target is rpi3_bl
kernel8.elf: clean font.o rpi3_bl.o
	tinygo build -target 'rpi3_bl.json' -ldflags 'font.o rpi3_bl.o' -no-debug  -o kernel8.elf .

font.o: font.psf
	ld.lld -m aarch64elf -r -b binary -o font.o font.psf

rpi3_bl.o: rpi3_bl.S
	clang --target=aarch64-elf -c rpi3_bl.S

clean:
	rm kernel8.elf kernel8.img *.o >/dev/null 2>/dev/null || true
