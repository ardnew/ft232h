[docimg]:https://godoc.org/github.com/ardnew/ft232h/native?status.svg
[docurl]:https://godoc.org/github.com/ardnew/ft232h/native

# ft232h/native

[![GoDoc][docimg]][docurl]

## Native FTDI Drivers
All communication with MPSSE-capable devices (including FT232H) is performed internally using FTDI's open-source driver [`libMPSSE`](https://www.ftdichip.com/Support/SoftwareExamples/MPSSE.htm). That software however depends on FTDI's proprietary, binary-only driver [`D2XX`](https://www.ftdichip.com/Drivers/D2XX.htm) (based on [`libusb`](https://github.com/libusb/libusb)), which is only available for certain host platforms.

To make these library dependencies as transparent to the user as possible - so that no configuration, compilation, or installation is required - these libraries have been modified to support static linkage and have been [re-compiled into a single static library archive `libft232h.a`](#building-libft232h-optional) for each supported OS.

The Go module uses [`cgo`](https://golang.org/cmd/cgo/) to automatically link against the correct static library based on the user's current OS and architecture.

This all happens internally so that applications importing the `ft232h` Go module do not have to explicitly use or specify the native drivers to use. The module and native drivers are built directly into the resulting Go application.

#### File structure for native drivers
All of the C software related to the native drivers `libft232h`, `libMPSSE`, and `D2XX` is contained underneath [`github.com/ardnew/ft232h/native/`](https://github.com/ardnew/ft232h/native). This includes both the source code required for building and the compiled executable code used by the `ft232h` Go module at compile-time. The files are organized as follows:
```sh
└── native/
    ├── lib/ # Pre-compiled libft232h.a libraries ...
    │   └── `${GOOS}_${GOARCH}`/ # .. separated by platform
    ├── inc/  # libMPSSE and D2XX C header APIs needed by cgo
    └── src/  # libMPSSE C source code, GNU Makefile
        └── `${GOOS}_${GOARCH}`/ # build outputs and proprietary D2XX library
```

#### Building `libft232h` (optional)
The static library `libft232h.a` can be rebuilt if necessary to support other platforms or if any changes to `libMPSSE` or `D2XX` are required.

Building the native drivers is done with a [GNU Makefile](src/Makefile). It performs the following tasks:
1. Compiles the `libMPSSE` C source code into object files (.o)
2. Extracts the object files from the proprietary `D2XX` static library
3. Archives all of the `libMPSSE` and `D2XX` object files into a single static library `libft232h.a`
4. Copies the `libft232h.a` static library to the necessary `ft232h` Go module subdirectory based on target OS and architecture.

Running `make` without any arguments will build `libft232h.a` for the default platform (`linux-amd64`) with the native `gcc` toolchain. To build for a different platform and/or use a cross-compiler, you must define the `platform` and/or `cross` variables when invoking `make`. See the `help` target (i.e. `$ make help`) for examples and all recognized values of `platform`.

The default target (`build`) also copies the compiled `libft232h.a` to the appropriate directory required by `cgo` on success, and should be all you need for the `ft232h` Go module to use the rebuilt library.

###### Building `libft232h` for other platforms
To support other platforms, you will need to make sure FTDI releases `D2XX` for that platform. You can view and download [official releases from **here**](https://www.ftdichip.com/Drivers/D2XX.htm). Once downloaded, you will want to copy the included static library to a subdirectory of `native/src/$(os)-$(arch)/libftd2xx/$(version)`, following existing convention, and add/update the various `$(*ftd2xx*)` definitions – as well as any necessary `CFLAGS` and `LDFLAGS` – in the `Makefile`.

After compiling and installing the `libft232h.a` static library, you will also need to update the `ft232h` Go module source file [`native_bridge.go`](https://github.com/ardnew/ft232h/native_bridge.go). The `cgo` preamble at the top of this file needs to include a valid, build-constrained, `-L<path>` option in `LDFLAGS` pointing to the path of your target's compiled `libft232h.a` static library. See the other supported targets in that file for examples.

