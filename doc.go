// Package ft232h provides a high-level interface to the FTDI FT232H USB to
// SPI/I²C/UART/GPIO protocol converter.
//
// Dependencies
//
// FTDI uses a custom (vendor-defined) protocol to communicate with their USB
// devices, and they release proprietary driver software `FTD2XX` (binary-only)
// that application programmers (see: YOU) should use with FTDI USB devices.
// These drivers are thus only available for systems officially supported by
// FTDI. These drivers are fairly low-level, with barebones C source code header
// and accompanying user guide as the only documentation resources available.
//
// Luckily, FTDI also develops a wrapper library `LibMPSSE` that greatly
// simplifies usage in the case of SPI, I²C, JTAG, and certain GPIO pins. This
// software is also open-source but is not guaranteed or supported by FTDI.
// However (of course), this software is a "wrapper" in the sense that it still
// depends on the proprietary (binary-only) `FTD2XX` driver software.
//
// Pre-compiled libraries - both `LibMPSSE` and `FTD2XX` - are included in this
// package for Linux (x86, AMD64, ARMv8) and macOS (AMD64). Source code and a
// GNU Makefile project has also been included to easily rebuild `LibMPSSE` for
// your target platform if you choose. However, the `FTD2XX` software cannot be
// rebuilt and must be downloaded from FTDI's Web site if you prefer to fetch
// your own copy.
//
// The following are links to documentation referenced above or are otherwise
// useful for FTDI USB device programming:
//
//                 Datasheet: https://github.com/ardnew/ft232h/blob/master/doc/FT232H-datasheet.pdf
//         FTD2XX User Guide: https://github.com/ardnew/ft232h/blob/master/doc/FTD2XX-user-guide.pdf
//   LibMPSSE I²C User Guide: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-I2C-user-guide.pdf
//   LibMPSSE SPI User Guide: https://github.com/ardnew/ft232h/blob/master/doc/LibMPSSE-SPI-user-guide.pdf
//
//       FTD2XX Product Page: https://www.ftdichip.com/Drivers/D2XX.htm
//     LibMPSEE Product Page: https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm
//
package ft232h
