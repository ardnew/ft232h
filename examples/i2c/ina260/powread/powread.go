package main

import (
	"log"

	"github.com/ardnew/ft232h"
)

func main() {

	ft, err := ft232h.NewFT232HWithDesc("FT232H-C")
	if nil != err {
		log.Fatalf("NewFT232HWithDesc(): %s", err)
	}
	defer ft.Close()
	log.Printf("%s", ft)

	if err := ft.I2C.Init(); nil != err {
		log.Fatalf("I2C.Init(): %v", err)
	}

	// set position at voltage register 0x02, default INA260 address 0x40
	if _, err := ft.I2C.Write(0x40, []uint8{0xFF}, true, true); nil != err {
		log.Fatalf("I2C.Write(): %v", err)
	}

	for {
		if id, err := ft.I2C.Read(0x40, 2, true, true); nil != err {
			log.Fatalf("I2C.Read(): %v", err)
		} else {
			log.Printf("reply: %+v", id)
		}
	}
}
