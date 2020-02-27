# ft232h (**`WIP`**)
##### Go module for [FTDI FT232H](https://www.ftdichip.com/Products/ICs/FT232H.htm) USB to GPIO/SPI/I²C/JTAG/UART protocol converter

_This is a work-in-progress and not at all stable_

## Features
- [x] `GPIO` - read/write
   - 8 dedicated pins available in any mode
   - 8-bit parallel, and 1-bit serial read/write operations
- [x] `SPI` - read/write
   - SPI `Mode0` and `Mode2` only, i.e. `CPHA=1`
   - configurable clock rate up to 30 MHz
   - chip/slave-select `CS` on both ports, pins `D3—D7`, `C0—C7`, including:
     - automatic assert-on-write/read, configurable polarity
     - multi-slave support with independent clocks `SCLK`, SPI modes, `CPOL`, etc.
   - unlimited effective transfer time/size
     - USB uses 64 KiB packets internally
- [ ] `I2C` - _not yet implementented_
- [ ] `JTAG` - _not yet implementented_
- [ ] `UART` - _not yet implementented_
- [x] **TBD** (WIP)

## Installation
Installing the module can be done with the Go built-in package manager. If you are not [using Go modules](https://blog.golang.org/using-go-modules) for your application (or are unsure), use the following:
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
# the official FTDI driver FTD2XX is incompatible with the VCP driver,
# preventing communication with FT232H breakouts
blacklist ftdi_sio
```
Be sure to unload the module if it was already loaded:
```sh
sudo rmmod ftdi_sio
```

#### macOS
Despite FTDI's [own quote from the `D2XX Programmer's Guide`](http://www.ftdichip.com/Support/Documents/ProgramGuides/D2XX_Programmer's_Guide(FT_000071).pdf) above, I've found that the current versions of macOS (10.13 and later, personal experience) have no problem co-existing with the `FTD2XX` driver included with this `ft232h` Go module. It _Just Works_ and no configuration is necessary.

## Usage
The obligatory ~~useless~~basic example:
Simply import the module and open the device:
```go
import (
  "log"
  "github.com/ardnew/ft232h"
)

func main() {
  // open the fist MPSSr-capable USB device found
  ft, err := ft232h.NewFT232H()
  if nil != err {
    log.Fatalf("NewFT232H(): %s", err)
  }
  defer ft.Close() // be sure to close device

  // do stuff
  log.Printf("ᵈᵒⁱⁿᵍ ˢᵗᵘᶠᶠ ᴅᴏɪɴɢ sᴛᴜғғ DOING STUFF: %s", ft)
}
```
I'm sure that was very helpful.

## Peripheral devices
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

## Drivers
All communication with MPSSE-capable devices (including FT232H) is performed internally using FTDI's open-source driver [`libMPSSE`](https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm). That software however depends on FTDI's proprietary, binary-only driver [`FTD2XX`](https://www.ftdichip.com/Drivers/D2XX.htm) (based on [`libusb`](https://github.com/libusb/libusb)), which is only available for certain host platforms.

To make these library dependencies as transparent to the user as possible - so that no configuration, compilation, or installation is required - these libraries have been modified to support static linkage and have been [re-compiled into a single static library archive `libft232h.a`](#building-libft232h-optional) for each supported OS.

The Go module uses [`cgo`](https://golang.org/cmd/cgo/) to automatically link against the correct static library based on the user's current OS and architecture.

This all happens internally so that applications importing the `ft232h` Go module do not have to explicitly use or specify the native drivers to use. The module and native drivers are built directly into the resulting Go application.

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
The static library `libft232h.a` can be rebuilt if necessary to support other platforms or if any changes to `libMPSSE` or `FTD2XX` are required.

Building the native drivers is done with a [GNU Makefile](native/src/Makefile). It performs the following tasks:
1. Compiles the `libMPSSE` C source code into object files (.o)
2. Extracts the object files from the proprietary `FTD2XX` static library
3. Archives all of the `libMPSSE` and `FTD2XX` object files into a single static library `libft232h.a`
4. Copies the `libft232h.a` static library to the necessary `ft232h` Go module subdirectory based on target OS and architecture.

Running `make` without any arguments will print the current build configuration and available `make` targets. You will need to define the `os` and `arch` variables at the top of the `Makefile` for your target system.

If you are cross-compiling, you will also want to define the `cross` variable based on the path-to or prefix-of your cross-compiler. For example, several prefixes are already defined (and commented out) for the supported platforms.

Once configured, simply running `make clean && make reinstall` should be all you need for the `ft232h` Go module to use the rebuilt library.

###### Building `libft232h` for other platforms
To support other platforms, you will need to make sure FTDI releases `FTD2XX` for that platform. You can view and download [official releases from **here**](https://www.ftdichip.com/Drivers/D2XX.htm). Once downloaded, you will want to copy the included static library to a subdirectory of `native/src/<$(os)>-$(arch)` and update the `$(ftd2xx-*)` version/path definitions as well as any necessary `CFLAGS` and `LDFLAGS` in the `Makefile`.

After compiling and installing the `libft232h.a` static library, you will also need to update the `ft232h` Go module source file [`native_bridge.go`](native_bridge.go). The `cgo` preamble at the top of this file needs to include a valid, build-constrained, `-L<path>` option in `LDFLAGS` pointing to the path of your target's compiled `libft232h.a` static library. See the other supported targets in that file for examples.

## Notes

#### Where to get one
[Adafruit sells a very nice breakout with a bunch of extras](https://www.adafruit.com/product/2264):
- USB-C and Stemma QT/Qwiic I²C connectors (with a little switch to short the chip's two awkward `SDA` pins!)
- On-board EEPROM (for storing chip configuration)
- 5V (`VBUS`) and 3.3V (on-board regulator, up to 500mA draw) outputs
