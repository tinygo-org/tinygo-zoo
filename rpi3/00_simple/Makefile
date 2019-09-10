#
# Simplest bootloader example.
#

all: clean kernel8.img

# note the target is rpi3_bl
kernel8.img: simple.go
	tinygo build --target=rpi3_bl -o kernel8.elf .

clean:
	rm kernel8.elf kernel8.img *.o >/dev/null 2>/dev/null || true
