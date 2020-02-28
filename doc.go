/*
Package ft232h provides a high-level interface to the FTDI FT232H USB to
SPI/I²C/UART/GPIO protocol converter.

Dependencies

FTDI uses a custom (vendor-defined) protocol to communicate with their USB
devices, and they release proprietary driver software `FTD2XX` (binary-only)
that application programmers (see: YOU) should use with FTDI USB devices.
These drivers are thus only available for systems officially supported by
FTDI. These drivers are fairly low-level, with barebones C source code header
and accompanying user guide as the only documentation resources available.

The following links contain useful documentation on FTDI USB device programming:

      FT232H Datasheet: https://github.com/ardnew/ft232h/blob/master/doc/FT232H-datasheet.pdf
         FTD2XX Manual: https://github.com/ardnew/ft232h/blob/master/doc/FTD2XX-user-guide.pdf
   libMPSSE I²C Manual: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-I2C-user-guide.pdf
   libMPSSE SPI Manual: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-SPI-user-guide.pdf
        FTD2XX Product: https://www.ftdichip.com/Drivers/D2XX.htm
      libMPSEE Product: https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm

*/
package ft232h
