package ft232h_test

import (
	"log"

	"github.com/ardnew/ft232h"
)

// Example demonstrates opening an FT232H with optional
// command line flags that specify which FTDI device to use
// if more than one is connected to the system.
//
// If no flags are provided, the first MPSSE-capable device
// found is used. Use -h to see all available flags.
//
// See the NewFT232H() godoc for other semantics related to
// the flag package.
//
// To open a specific device without using command line
// flags, use one of the functions of form NewFT232HWith*()
// In particular, NewFT232HWithMask(nil) will open the
// first compatible device found.
func ExampleFt232h() {
	// open the first device found that matches all
	// command line flags (if any provided)
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close() // be sure to close device

	// at this point you can call Init() or Config() on
	// one of its fields GPIO, SPI, I2C, ...
	doStuff(ft)
}

// doStuff does stuff with an FT232H
func doStuff(ft *ft232h.FT232H) {
	// FT232H implements String() descriptively
	log.Printf("using: %s", ft)
}
