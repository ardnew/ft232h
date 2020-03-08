package ft232h_test

import (
	"log"

	"github.com/ardnew/ft232h"
)

func Example_command_line() {

	// Call NewFT232H() to open an FT232H device from a command line-oriented
	// application to help select which FTDI device to use (by parsing predefined
	// command line flags) if more than one is connected to the system.
	//
	// If no flags are provided, the first MPSSE-capable USB device found is used.
	// Use -h to see all available flags.
	//
	// See the NewFT232H() godoc for other semantics related to the flag package.
	//
	// To open a specific device without using command line flags, use one of the
	// functions of form NewFT232HWith*(). In particular, NewFT232HWithMask(nil)
	// will open the first compatible device found.

	// Open first device that matches all command line flags (if any provided)
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close() // be sure to close device

	// At this point, you can call Init() or Config() on one of the interface
	// fields GPIO, SPI, I2C, ...
	log.Printf("using: %s", ft) // FT232H implements String() descriptively
}
