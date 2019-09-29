### Bootloader TODO

1. Figure out why the blinker example has bad bytes in its read-write data. Once
this is fixed, return "byte mismatch" to being a fatal error.

2. Figure out how to clear out the buffer (buffers?) at startup so there are
not "left over" lines in the buffers when the host side starts up again.

3. Determine what state the rpi3 gets into when it stops responding to DATA
requests and we are forced to just do a reset.
