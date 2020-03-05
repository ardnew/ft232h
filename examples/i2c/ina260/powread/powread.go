package main

import (
	"log"

	"github.com/ardnew/ft232h"
)

const (
	slave = 0x40 // default INA260 slave address
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

	// set register pointer at 0x02 (voltage)
	if _, err := ft.I2C.Write(slave, []uint8{0x02}, true, true); nil != err {
		log.Fatalf("I2C.Write(): %v", err)
	}

	for {
		if id, err := ft.I2C.Read(slave, 2, true, true); nil != err {
			log.Fatalf("I2C.Read(): %v", err)
		} else {
			log.Printf("reply: %+v", id)
		}
	}
}
