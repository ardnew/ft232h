/*
Package ft232h is a Go module for FTDI FT232H USB to GPIO/SPI/I²C/JTAG/UART
protocol converter.

Dependencies

FTDI uses a custom (vendor-defined) protocol to communicate with their USB
devices and releases proprietary driver software `D2XX` (binary-only) that
application programmers (see: YOU) should use with FTDI USB devices.

The `D2XX` driver is thus only available for systems officially supported by
FTDI. These drivers are fairly low-level, containing only C source code headers
and a thin user guide as documentation.

FTDI Driver Documentation

The following links contain the relevant FTDI documentation for the software
versions used in the ft232h Go module:

     FT232H Datasheet: https://github.com/ardnew/ft232h/blob/master/doc/FT232H-datasheet.pdf
          D2XX Manual: https://github.com/ardnew/ft232h/blob/master/doc/D2XX-user-guide.pdf
  libMPSSE I²C Manual: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-I2C-user-guide.pdf
  libMPSSE SPI Manual: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-SPI-user-guide.pdf
         D2XX Product: https://www.ftdichip.com/Drivers/D2XX.htm
     libMPSEE Product: https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm

*/
package ft232h
