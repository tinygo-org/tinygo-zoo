
# aliases
.PHONY: clean arduino-colorlamp microbit-blink microbit-pixelbuttons nrf-colorlamp

clean:
	mkdir -p build
	rm -rf build/**

arduino-colorlamp:
	make clean
	docker run --rm -v $(PWD):/src hybridgroup/tinygo-all build -o /src/build/arduino-colorlamp.hex -target arduino /src/arduino-colorlamp/main.go
	avrdude -c arduino -p atmega328p -P /dev/ttyACM0 -U flash:w:build/arduino-colorlamp.hex

microbit-blink:
	make clean
	docker run --rm -v $(PWD):/src hybridgroup/tinygo-all build -o /src/build/microbit-blink.hex -target microbit /src/microbit-blink/main.go
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-blink.hex reset exit'

microbit-pixelbuttons:
	make clean
	docker run --rm -v $(PWD):/src hybridgroup/tinygo-all build -o /src/build/microbit-pixelbuttons.hex -target microbit /src/microbit-pixelbuttons/main.go
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-pixelbuttons.hex reset exit'

nrf-colorlamp:
	make clean
	docker run --rm -v $(PWD):/src hybridgroup/tinygo-all build -o /src/build/nrf-colorlamp.hex -target pca10040 /src/nrf-colorlamp/main.go
	nrfjprog -f nrf52 --sectorerase --program build/nrf-colorlamp.hex --reset
