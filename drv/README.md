# ft232h
## Peripheral device drivers
The following device drivers are implemented as individual packages using the `github.com/ardnew/ft232h` Go module.

#### GPIO
> TBD

#### SPI

##### [`github.com/ardnew/ft232h/drv/ili9341`](https://github.com/ardnew/ft232h/drv/ili9341) 
Driver for the ILI9341 320x240 TFT LCD chipset using the `ft232h.SPI` and `ft232h.GPIO` interfaces – including methods to draw pixels, rectangles, and 16-bit RGB bitmaps.

These methods alone are sufficient to implement a demo application – [`boing`](https://github.com/ardnew/ft232h/examples/spi/ili9341/boing).

#### I²C
> TBD

#### JTAG
> TBD

#### UART
> TBD

