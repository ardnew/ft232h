## **This is a brief summary. See [`READMORE`](READMORE.md) for the complete overview.**

# ft232h
##### Go module for the [FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm) USB to SPI/I²C/UART Protocol Converter with GPIO

_This is a work-in-progress and not at all stable_

## Features
- [x] **TBD** (WIP)

## Installation
Installation is conventional, just use the Go built-in package manager:
```sh
go get -v github.com/ardnew/ft232h
```
No other libraries or configuration is required. 

###### Common issues
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
	// open the first MPSSE-capable USB device found
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close() // be sure to close device
	log.Printf("%s", ft)
}
```

## Peripherals
#### Getting started
The great thing about the FT232H is being able to communicate over regular USB with peripheral devices using GPIO and several different serial protocols – SPI and I²C in particular – and being able to use standard Go straight from your PC makes that even greater. No need to stumble around programming a microcontroller just to mediate the communication.

Adding support for a peripheral device is straight-forward. You can either create a new driver package under [`drv/`](drv), or you can simply interact with the interface directly from your application. 

To demonstrate this, a basic [driver package](drv/ili9341) was created to drive an ILI9341 320x240 TFT LCD using the `ft232h.SPI` interface with methods to draw pixels, rectangles, and 16-bit RGB bitmaps. 

These methods alone in [the driver package](drv/ili9341) were sufficient to implement the other half of this demonstration – an [example application](examples/spi/ili9341/boing) using the ILI9341 driver package. This application was a port of the [tinygo project](https://tinygo.org/)'s ILI9341 device driver example [`pyportal_boing`](https://github.com/tinygo-org/drivers/tree/master/examples/ili9341/pyportal_boing), which was in turn a port of Adafruit's [original Arduino demo](https://github.com/adafruit/Adafruit_ILI9341/tree/master/examples/pyportal_boing) released for their PyPortal.

So to get started, please review the [driver package](drv/ili9341) and companion application [`boing`](examples/spi/ili9341/boing) for details on how this `ft232h` go module is intended for use, which should also help you become familiar with the API and general architecture. The design was intended to be as concise and general-purpose as possible, to not litter the namespace with subtleties, yet low-level enough to wield some _Real Power_.

In the mean-time, **_[hold on to your butts](https://www.youtube.com/watch?v=-W6as8oVcuM)_**, and watch the above mentioned `ili9341` driver package with `boing` companion application in full glorious 320x240 16-bit RGB @ 42.67 FPS over 30 MHz SPI – all written in standard Go running straight from my desktop!

<p align=center>
	<a href="https://www.youtube.com/watch?v=H-9oN2VmrUw">
		ILI9341 with FT232H and golang<br/>
		<img src="https://img.youtube.com/vi/H-9oN2VmrUw/0.jpg" alt="ILI9341 with FT232H and golang"><br/>
		https://www.youtube.com/watch?v=H-9oN2VmrUw<br/>
	</a>
</p>
