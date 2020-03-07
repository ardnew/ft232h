[docimg]:https://godoc.org/github.com/ardnew/ft232h?status.svg
[docurl]:https://godoc.org/github.com/ardnew/ft232h
[cciimg]:https://circleci.com/gh/ardnew/ft232h.svg?style=shield
[cciurl]:https://circleci.com/gh/ardnew/ft232h
[repimg]:https://goreportcard.com/badge/github.com/ardnew/ft232h
[repurl]:https://goreportcard.com/report/github.com/ardnew/ft232h

# ft232h
### Go module for [FTDI FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm) USB to GPIO/SPI/I²C/JTAG/UART protocol converter

[![GoDoc][docimg]][docurl] [![CircleCI][cciimg]][cciurl] [![Go Report Card][repimg]][repurl]

> **This is a brief summary. See [`READMORE`](READMORE.md) for the complete overview.**

## API features
#### This software is a work-in-progress (WIP) and not ready for use. The following features have been implemented, but their interfaces _may_ (will) change.
- [x] [**Documented**][docurl] and [**integration tested**][cciurl]
- [x] `GPIO` - read/write
   - 8 dedicated pins available in any mode
   - 8-bit parallel, and 1-bit serial read/write operations
- [x] `SPI` - read/write
   - SPI modes `0` and `2` only, i.e. `CPHA=1`
   - configurable clock rate up to 30 MHz
   - chip/slave-select `CS` on both ports (pins `D3—D7`, `C0—C7`), including:
     - automatic assert-on-write/read with configurable polarity
     - multi-slave support with independent clocks `SCLK`, SPI modes, `CPOL`, etc.
   - unlimited effective transfer time/size
     - USB uses 64 KiB packets internally
- [x] `I2C` - read/write
   - configurable clock rate up to high speed mode (3.4 Mb/s)
   - internal or external SDA pullup option
   - unlimited effective transfer time/size
     - USB uses 64 KiB packets internally
- [ ] `JTAG` - _not yet implementented_
- [ ] `UART` - _not yet implementented_
- [x] **TBD** (WIP)

## Installation
Installation is conventional, just use the Go built-in package manager:
```sh
go get -v github.com/ardnew/ft232h
```
No other libraries or configuration is required.

###### Linux
If you have trouble finding/opening your device in Linux, you probably have the incompatible module `ftdi_sio` loaded. See the Linux `Installation` section in [`READMORE`](READMORE.md) for details.

## Supported platforms
Internally, `ft232h` depends on some proprietary software from FTDI that is only available for a handful of platforms (binary-only). This would therefore be the only platforms supported by the `ft232h` Go module:
#### Linux
- [x] x86 (32-bit) `[386]`
- [x] x86_64 (64-bit) `[amd64]`
- [x] ARMv7 (32-bit) `[arm]` - includes Raspberry Pi models 3 and 4
- [x] ARMv8 (64-bit) `[arm64]` - includes Raspberry Pi model 4
#### macOS
- [x] x86_64 (64-bit) `[amd64]`
#### Windows
- [ ] x86 (32-bit) `[386]`
- [ ] x86_64 (64-bit) `[amd64]`
###### Windows compatibility
Windows support is possible – and in fact appears to be FTDI's preferred target – but drivers for this `ft232h` Go module have not been compiled or tested. The modifications made to `libMPSSE` to support static linkage would need to be verified or merged in for Windows. See the `Drivers` section in [`READMORE`](READMORE.md) for info.

## Usage
Simply import the module and open the device:
```go
import (
	"log"
	"github.com/ardnew/ft232h"
)

func main() {
	// open the fist MPSSE-capable USB device found
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close() // be sure to close device

	// do stuff
	log.Printf("doing stuff with device: %s", ft)
}
```

## Peripheral devices
Of course the FT232H isn't that useful without a device to interact with. You will still have to write drivers or adapters for your particular device. But it's a lot more fun developing device drivers with the full, native Go ecosystem on your PC at your disposal!

A basic [driver for the ILI9341 TFT LCD](drv/ili9341) using `ft232h.SPI` and `ft232h.GPIO` – along with a [demo application](examples/spi/ili9341/boing) drawing an animated bouncing ball – has been created to serve as a reference implementation.

For more details, be sure to read the `Peripheral devices` section  in [`READMORE`](READMORE.md) and, of course, the [godoc][docurl] for this `ft232h` module.
