[docimg]:https://godoc.org/github.com/ardnew/ft232h/drv?status.svg
[docurl]:https://godoc.org/github.com/ardnew/ft232h/drv

# ft232h/drv

[![GoDoc][docimg]][docurl]

## Peripheral device drivers
The official device drivers required to control the FT232H are quite complex and require an understanding of microcontroller programming in C. This `ft232h` module greatly simplifies that programming interface, bridging many of those peripherals with the native Go ecosystem on a common PC.

Adding support for a peripheral device is straight-forward. There's no _required_ Go `interface` patterns to implement. Each serial (and GPIO) capability of the FT232H is exposed as a named member of the `type FT232H struct`. Each member has its own conventional methods associated with it (configure, read, write, etc.), which a device driver can interact with directly.

The idea is to create a single driver package at `drv/FOO` to implement the definitions and procedures provided by each supported peripheral device. This would allow a consuming application to reuse one driver without having to import the entire driver library.

Alternatively, you can just interact with the `ft232h` module directly from an application.

See the [godoc](https://godoc.org/github.com/ardnew/ft232h) documentation for API details.

# Supported devices
The following peripheral devices currently have driver support.

## GPIO
> TBD

## SPI
 - [x] **ILI9341** - [`github.com/ardnew/ft232h/drv/ili9341`](ili9341)
   - Driver for the ILI9341 320x240 TFT LCD chipset using `ft232h.SPI` and `ft232h.GPIO` interfaces – including methods to draw pixels, rectangles, and 16-bit RGB bitmaps.
   - [`boing`](../examples/spi/ili9341/boing) - demo application

## I²C
> TBD

## JTAG
> TBD

## UART
> TBD
