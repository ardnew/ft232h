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
Adding support for a peripheral SPI/I²C device is straight-forward. You can either create a new driver package under [`drv/`](drv), or you can simply interact with the interface directly from your application. 

To demonstrate this, a couple packages were created to act as reference and example. Please follow along in the `Peripherals` section in [`READMORE`](READMORE.md) for guidance.
