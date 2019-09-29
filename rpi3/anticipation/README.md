### Notice
This program runs on the host system (linux, darwin, etc) and requires
a "normal" copy of go.  The go compiler should be at least version 1.11.

See the `Makefile` for how to run this program with QEMU as a simulator.
See the `Makefile` for how to run this against real hardware.


### What you will see
You will see a sequence of log messages (lines preceded by timestamps like
2019/09/09 20:14:33) as the bootloading process goes on.  Once the program
is loaded and started, you'll see a log message that says
"starting terminal loop...."   After that, the terminal is connected to
your program and you'll notice the text sent from the device is not preceded
by a timestamp.

# Running with QEMU (copied from Makefile)
1. you have to use the right TTY that is created by anticipationbl when
it starts (via QEMU).  You should run anticipationbl first so you can see the output
of the device (tty) to use to connect to it.
2. You need to choose what you want the bootloader to load.  In the
example below, we are loading the kernel (must be an ELF file) created
by building in the 00_simple directory.

# Running on Hardware (copied from Makefile)
1. you need to know where your Serial-To-USB device is.  I have
listed the device below as it exists on my Mac.
2. You need to choose what you want the bootloader to load.  In the
example below, we are loading the kernel (must be an ELF file) created
by building in the 00_simple directory.
