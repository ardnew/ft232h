# ft232h
##### Go module for the [FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm) USB to SPI/I²C/UART Protocol Converter with GPIO

_This is a work-in-progress and not at all stable_

## Features
- [x] Go **module** compatible with `go get` (see: [Installation](#installation))
  - No installation or configuration required
- [x] **No** dynamic library dependencies (`libMPSSE`, `FTD2XX`, etc.)
  - Go applications using the module need no additional libraries to be packaged or deployed with the compiled executable
  - Native drivers are statically linked with Go module, transparent to the consuming application
- [x] Support for multiple host OS (see: [Drivers](#drivers))
  - Linux 32-bit (`386`) and 64-bit (`amd64`, `arm64`) - includes Raspberry Pi models 3 and 4
  - macOS (`amd64`)
  - Windows not currently supported
- [x] **TBD** (WIP)

## Drivers
All communication with MPSSE-capable devices (including FT232H) is performed internally using FTDI's open-source driver [`libMPSSE`](https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm). That software however depends on FTDI's proprietary, binary-only driver [`FTD2XX`](https://www.ftdichip.com/Drivers/D2XX.htm) (based on [`libusb`](https://github.com/libusb/libusb)), which is only available for certain host platforms.

To make these library dependencies as transparent to the user as possible - so that no configuration, compilation, or installation is required - these libraries have been modified to support static linkage and have been [re-compiled into a single static library archive `libft232h.a`](#building-libft232h-optional) for each supported OS. The Go module uses [`cgo`](https://golang.org/cmd/cgo/) to automatically link against the correct static library based on the user's current OS and architecture. This all happens internally so that applications importing the `ft232h` Go module do not have to explicitly use or specify the native drivers to use. The module and native drivers are built directly into the resulting Go application.

#### File structure for native drivers
All of the C software related to the native drivers `libft232h`, `libMPSSE`, and `FTD2XX` is contained underneath [`native/`](native). This includes both the source code required for building and the compiled executable code used by the `ft232h` Go module at compile-time. The files are organized as follows:
```sh
└── native/
    ├── lib/ # Pre-compiled libft232h.a libraries ...
    │   └── `${GOOS}_${GOARCH}`/ # .. separated by platform
    ├── inc/  # libMPSSE and FTD2XX C header APIs needed by cgo
    └── src/  # libMPSSE C source code, GNU Makefile
        └── `${GOOS}_${GOARCH}`/ # build outputs and proprietary FTD2XX library
```

#### Building `libft232h` (optional)
The static library `libft232h.a` can easily be rebuilt if necessary to support other platforms or if any changes to `libMPSSE` or `FTD2XX` are required.

A [GNU Makefile](native/src/Makefile) was created to simplify the build and installation process. It performs the following tasks:
1. Compiles the `libMPSSE` C source code into object files (.o)
2. Extracts the object files from the proprietary `FTD2XX` static library
3. Archives all of the `libMPSSE` and `FTD2XX` object files into a single static library `libft232h.a`
4. Copies the `libft232h.a` static library to the necessary `ft232h` Go module subdirectory based on target OS and architecture.

Running `make` without any arguments will print the current build configuration and available `make` targets. You will need to define the `os` and `arch` variables at the top of the `Makefile` for your target system. If you are cross-compiling, you will also want to define the `cross` variable based on the path-to or prefix-of your cross-compiler. For example, several prefixes are already defined (and commented out) for the supported platforms. Once configured, simply running `make clean && make reinstall` should be all you need for the `ft232h` Go module to use the rebuilt library.

#### Building `libft232h` for other platforms (optional)
To support other platforms, you will need to make sure FTDI releases `FTD2XX` for that platform. Once downloaded, you will want to copy the included static library to a subdirectory of `native/src/<$(os)>-$(arch)` and update the `$(ftd2xx-*)` version/path definitions as well as any necessary `CFLAGS` and `LDFLAGS` in the `Makefile`.

After compiling and installing the `libft232h.a` static library, you will also need to update the `ft232h` Go module source file [`native_bridge.go`](native_bridge.go). The `cgo` preamble at the top of this file needs to include a valid, build-constrained, `-L<path>` option in `LDFLAGS` pointing to the path of your target's compiled `libft232h.a` static library. See the other supported targets in that file for examples.

## Notes

#### Where to get one
[Adafruit sells a very nice breakout with a bunch of extras](https://www.adafruit.com/product/2264):
- USB-C and Stemma QT/Qwiic I²C connectors (with a little switch to short the chip's two awkward `SDA` pins!)
- On-board EEPROM (for storing chip configuration)
- 5V (`VBUS`) and 3.3V (on-board regulator, up to 500mA draw) outputs