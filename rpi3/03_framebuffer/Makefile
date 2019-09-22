#
# Simplest bootloader example.
#

all: clean kernel8.elf

# note the target is rpi3_bl
kernel8.elf: clean font.o
	tinygo build -ldflags "-m aarch64elf font.o" --target=rpi3_bl -o kernel8.elf .

font.o: font.psf
	ld.lld -m aarch64elf -r -b binary -o font.o font.psf

#ld.lld -m aarch64elf -nostdlib start.o font.o $(OBJS) -T link.ld -o kernel8.elf

clean:
	rm kernel8.elf kernel8.img *.o >/dev/null 2>/dev/null || true
