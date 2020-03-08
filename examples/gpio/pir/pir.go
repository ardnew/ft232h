package main

import (
	"log"

	"github.com/ardnew/ft232h"
)

func main() {

	// open the FT232H
	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close()

	if err := ft.I2C.Init(); nil != err {
		log.Fatalf("I2C.Init(): %v", err)
	}

	last := true
	for {
		set, err := ft.GPIO.Get(ft232h.C(5))
		if nil != err {
			log.Fatalf("GPIO.Get(): %v", err)
		}
		if last != set {
			last = set
			log.Printf("change!")
		}
	}

}
