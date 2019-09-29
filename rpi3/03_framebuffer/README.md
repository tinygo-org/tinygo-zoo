This program demonstrates three things:

* Putting a font into your binary so you can use it with bare metal.
* Rendering the font on screen on the pi
* Scrolling so you have something that works like a console

The scrolling is GPU assisted.  If you don't have the maximum height of your
framebuffer correctly set, this will fail-- probably about 15 lines after
the line marked with '<-------' or you will get an exception.

You need to add this line to your config.txt file on your SD card (this is really
setting a BIOS option)

```
max_framebuffer_height=1536
```

Here is my complete config.txt if you want to duplicate it exactly:

```
# Force the monitor to HDMI mode so that sound will be sent over HDMI cable
hdmi_drive=2
# Set monitor mode to DMT
hdmi_group=2
# Set monitor resolution to 1024x768 XGA 60Hz (HDMI_DMT_XGA_60)
hdmi_mode=16
# don't show rainbow screen
disable_splash=1
# more GPU space, on 1GB machines, use 128MB for GPU
gpu_mem_1024=128
# set FB height for our fancy scrolling
max_framebuffer_height=1536
```
