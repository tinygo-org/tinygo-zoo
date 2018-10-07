
# aliases
.PHONY: clean

TARGET ?= unix

ifeq ($(TARGET),unix)
# Regular *nix system.

else ifeq ($(TARGET),pca10040)
# PCA10040: nRF52832 development board
OBJCOPY = arm-none-eabi-objcopy
TGOFLAGS += -target $(TARGET)

else ifeq ($(TARGET),microbit)
# BBC micro:bit
OBJCOPY = arm-none-eabi-objcopy
TGOFLAGS += -target $(TARGET)

else ifeq ($(TARGET),bluepill)
# "blue pill" development board
# See: https://wiki.stm32duino.com/index.php?title=Blue_Pill
OBJCOPY = arm-none-eabi-objcopy
TGOFLAGS += -target $(TARGET)

else ifeq ($(TARGET),arduino)
OBJCOPY = avr-objcopy
TGOFLAGS += -target $(TARGET)

else
$(error Unknown target)

endif

ifeq ($(TARGET),pca10040)
flash-%: build/%.hex
	nrfjprog -f nrf52 --sectorerase --program $< --reset
else ifeq ($(TARGET),microbit)
flash-%: build/%.hex
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program $< reset exit'
else ifeq ($(TARGET),arduino)
flash-%: build/%.hex
	avrdude -c arduino -p atmega328p -P /dev/ttyACM0 -U flash:w:$<
else ifeq ($(TARGET),bluepill)
flash-%: build/%.hex
	openocd -f interface/stlink-v2.cfg -f target/stm32f1x.cfg -c 'program $< reset exit'
endif

clean:
	@rm -rf build

# ELF file that can run on a microcontroller.
build/%.elf: % %/*.go
	tinygo build $(TGOFLAGS) -size=short -o $@ ./$<

# Convert executable to Intel hex file (for flashing).
build/%.hex: build/%.elf
	$(OBJCOPY) -O ihex $^ $@
