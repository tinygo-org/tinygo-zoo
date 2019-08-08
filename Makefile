
# aliases
.PHONY: clean arduino-colorlamp build-arduino-colorlamp flash-arduino-colorlamp microbit-blink build-microbit-blink flash-microbit-blink microbit-pixelbuttons build-microbit-pixelbuttons flash-microbit-pixelbuttons nrf-colorlamp build-nrf-colorlamp flash-nrf-colorlamp microbit-accel build-microbit-accel flash-microbit-accel reelboard-accel build-reelboard-accel flash-reelboard-accel

clean:
	mkdir -p build
	rm -rf build/**

build-arduino-colorlamp:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/arduino-colorlamp.hex -target arduino /src/arduino-colorlamp/main.go

flash-arduino-colorlamp:
	avrdude -c arduino -p atmega328p -P /dev/ttyACM0 -U flash:w:build/arduino-colorlamp.hex

arduino-colorlamp:
	make clean
	make build-arduino-colorlamp
	make flash-arduino-colorlamp

build-microbit-blink:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/microbit-blink.hex -target microbit /src/microbit-blink/main.go

flash-microbit-blink:
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-blink.hex reset exit'

microbit-blink:
	make clean
	make build-microbit-blink
	make flash-microbit-blink

build-microbit-pixelbuttons:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/microbit-pixelbuttons.hex -target microbit /src/microbit-pixelbuttons/main.go

flash-microbit-pixelbuttons:
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-pixelbuttons.hex reset exit'

microbit-pixelbuttons:
	make clean
	make build-microbit-pixelbuttons
	make flash-microbit-pixelbuttons

build-microbit-images:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.6.1 tinygo build -o /src/build/microbit-images.hex -target microbit /src/microbit-images/main.go

flash-microbit-images:
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-images.hex reset exit'

microbit-images:
	make clean
	make build-microbit-pixelbuttons
	make flash-microbit-pixelbuttons

build-microbit-accel:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/microbit-accel.hex -target microbit /src/accel/main.go

flash-microbit-accel:
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/microbit-accel.hex reset exit'

microbit-accel:
	make clean
	make build-microbit-accel
	make flash-microbit-accel

build-reelboard-accel:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/reelboard-accel.hex -target reelboard /src/accel/main.go

flash-reelboard-accel:
	openocd -f interface/cmsis-dap.cfg -f target/nrf51.cfg -c 'program build/reelboard-accel.hex reset exit'

reelboard-accel:
	make clean
	make build-reelboard-accel
	make flash-reelboard-accel

build-nrf-colorlamp:
	docker run --rm -v "$(PWD):/src" -v "$(GOPATH):/gohost" -e "GOPATH=$(GOPATH):/gohost" tinygo/tinygo:0.7.1 tinygo build -o /src/build/nrf-colorlamp.hex -target pca10040 /src/nrf-colorlamp/main.go

flash-nrf-colorlamp:
	nrfjprog -f nrf52 --sectorerase --program build/nrf-colorlamp.hex --reset

nrf-colorlamp:
	make clean
	make build-nrf-colorlamp
	make flash-nrf-colorlamp
