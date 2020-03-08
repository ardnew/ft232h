[docimg]:https://godoc.org/github.com/ardnew/ft232h?status.svg
[docurl]:https://godoc.org/github.com/ardnew/ft232h
[cciimg]:https://circleci.com/gh/ardnew/ft232h.svg?style=shield
[cciurl]:https://circleci.com/gh/ardnew/ft232h
[repimg]:https://goreportcard.com/badge/github.com/ardnew/ft232h
[repurl]:https://goreportcard.com/report/github.com/ardnew/ft232h

# ft232h
### Go module for [FTDI FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm) USB to GPIO/SPI/I²C/JTAG/UART protocol converter

[![GoDoc][docimg]][docurl] [![CircleCI][cciimg]][cciurl] [![Go Report Card][repimg]][repurl]

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
If you are not [using Go modules](https://blog.golang.org/using-go-modules) for your application (or are unsure), use the built-in `go` package manager:
```sh
go get -u -v github.com/ardnew/ft232h
```
Otherwise, you are using Go modules, either use the same command above (sans `-u`), or simply add the import statement to your source code and the module will be installed automatically:
```go
import (
	// ... other imports ...
	"github.com/ardnew/ft232h"
)
```
No other files or configuration to your build process are necessary.

#### Linux
Many Linux distributions ship with the FTDI Virtual COM Port (VCP) driver pre-installed (as a kernel module, usually `ftdi_sio`). However, [according to FTDI](http://www.ftdichip.com/Support/Documents/ProgramGuides/D2XX_Programmer's_Guide(FT_000071).pdf):
> For Linux, Mac OS X (10.4 and later) and Windows CE (4.2
> and later) the D2XX driver and VCP driver are mutually
> exclusive options as only one driver type may be installed
> at a given time for a given device ID.

There are [a lot of ways](https://www.google.com/search?q=d2xx+ftdi_sio) to resolve the issue, including [fancy udev rules to swap out modules when (un)plugging devices](https://stackoverflow.com/a/43514662/1054397), but I don't personally use the VCP driver.

On Ubuntu, you can simply prevent the VCP module from being auto-loaded at bootup by blacklisting the module. For example, create a new file `/etc/modprobe.d/blacklist-ftdi.conf` with a single directive:
```sh
# the official FTDI driver D2XX is incompatible with the VCP driver,
# preventing communication with FT232H breakouts
blacklist ftdi_sio
```
Be sure to unload the module if it was already loaded:
```sh
sudo rmmod ftdi_sio
```

#### macOS
Despite FTDI's [own quote from the `D2XX Programmer's Guide`](http://www.ftdichip.com/Support/Documents/ProgramGuides/D2XX_Programmer's_Guide(FT_000071).pdf) above, I've found that the current versions of macOS (10.13 and later, personal experience) have no problem co-existing with the `D2XX` driver included with this `ft232h` Go module. It _Just Works_ and no configuration is necessary.

## Usage
> The obligatory ~~useless~~basic example

Simply import the module and open the device:
```go
import (
	"log"
	"github.com/ardnew/ft232h"
)

func main() {
	// open the first device found, matching all command-line flags if provided
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close() // be sure to close device

	// do stuff
	log.Printf("doing stuff with device: %s", ft)
}
```
I'm sure that was very helpful.

#### Design
The design of this module API was intended to marry the following principles, ordered by importance:
1. **Simple / Concise** – Clear, consistent, conventional behavior
   - Transparent infrastructure, automatic configuration to service methods
   - Minimize subtleties and unapparent side-effects (or document those that exist, are inherent)
   - Prefer Go native or simple composite types
   - Modular protocol-based architecture
   - Integrated command-line flag support
2. **Robust / General Purpose** – Maximize peripheral device support
   - Clear, descriptive error messages for invalid configurations
   - Low-level device methods, flexible high-level protocol methods
   - Integration test every revision automatically
   - Correct behavior at edge cases
3. **Performant / Efficient** – Utilize host and device resources effectively
   - Minimize USB transactions and HID interframe delays
   - Maximize throughput of serial protocol transfers

## Peripheral devices

See [`drv/`](https://github.com/ardnew/ft232h/drv) for a list of existing device drivers.

See the [package documentation][docurl] for info on developing device drivers and the API for using this module directly.

## Examples

Demo applications using this module and its device drivers can be found in [`examples/`](https://github.com/ardnew/ft232h/examples).

## Notes

#### Where to get one
[Adafruit sells a very nice breakout with a bunch of extras](https://www.adafruit.com/product/2264):
- USB-C and Stemma QT/Qwiic I²C connectors (with a little switch to short the chip's two awkward `SDA` pins!)
- On-board EEPROM (for storing chip configuration)
- 5V (`VBUS`) and 3.3V (on-board regulator, up to 500mA draw) outputs
