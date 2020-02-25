# ft232h
##### Go module for the FT232H USB to SPI/I²C/UART Protocol Converter with GPIO

_This is a work-in-progress and not at all stable_

## Features
- [x] Go module compatible (see `go mod`)
- [x] Designed for and tested with [FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm)
  - [Adafruit sells a very nice breakout with a bunch of extras](https://www.adafruit.com/product/2264):
    - USB-C and Stemma QT/Qwiic I²C connectors (with a little switch to short the chip's two awkward `SDA` pins!)
    - On-board EEPROM (for storing chip configuration)
    - 5V (`VBUS`) and 3.3V (on-board regulator, up to 500mA draw) outputs
- [x] Includes re-compilable native FTDI drivers for multiple host OS
  - Linux 32-bit (`386`) and 64-bit (`amd64`, `arm64`) - includes Raspberry Pi models 3 and 4
  - macOS (`amd64`)
  - Windows not currently supported
- [x] **TBD** (WIP)

## Drivers
All communication with MPSSE-capable devices (including FT232H) is performed with FTDI's open-source driver [`LibMPSSE`](https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm). That software however depends on FTDI's proprietary driver [`FTD2XX`](https://www.ftdichip.com/Drivers/D2XX.htm) (based on [`libusb`](https://github.com/libusb/libusb)), which is only available for certain host platforms.

Contained in this project are all of the necessary C source files required to build `LibMPSSE`, as well as the required `FTD2XX` library (binary-only), with modifications to support being built as a single statically-linked library. A simple GNU Makefile (shared among all supported OS) has also been created to simplify building (see: [Building LibMPSSE](#building-libmpsse-optional)). The result is a single library file `libMPSSE.a` containing _all_ of the necessary FTDI driver dependencies with which the `ft232h` Go module can be linked.

A pre-compiled `libMPSSE.a` library is already included with this package for each supported OS, so no special configuration or installation is required.

Under [`native`](native), you will find the headers needed by the `ft232h` Go module to communicate with the C library (using [`cgo`](https://golang.org/cmd/cgo/)), the source code for `LibMPSSE`, the pre-compiled `FTD2XX` libraries, and the pre-compiled `libMPSSE.a` libraries separated for each supported platform:

```sh
└── native/
    ├── lib/  # Pre-compiled libMPSSE.a library needed by Go software using this ft232h module
    │   └── `${GOOS}_${GOARCH}`/ # separated by platform for cgo library path resolution
    ├── inc/  # LibMPSSE C APIs and FTD2XX C source code headers needed by cgo
    └── src/  # LibMPSSE C source code, GNU Makefile
        └── `${GOOS}_${GOARCH}`/ # build outputs and proprietary FTD2XX library
```

#### Building LibMPSSE (optional)
**TBD** (WIP)
