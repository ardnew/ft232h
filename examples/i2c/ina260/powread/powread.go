package main

import (
	"log"

	"github.com/ardnew/ft232h"
)

const (
	slave = 0x40       // default INA260 slave address
	order = ft232h.MSB // byte order of INA260 data transfers

	// register addresses
	cAddr = 0x01 // voltage
	vAddr = 0x02 // current
	pAddr = 0x03 // power

	// precision of the INA260 each register's LSB
	cLSB = 1.25 // milliamps
	vLSB = 1.25 // millivolts
	pLSB = 10.0 // milliwatts

)

func main() {

	// open the FT232H
	ft, err := ft232h.New()
	if nil != err {
		log.Fatalf("New(): %s", err)
	}
	defer ft.Close()
	log.Printf("%s", ft)

	// initialize FT232H MPSSE engine in I²C mode
	if err := ft.I2C.Init(); nil != err {
		log.Fatalf("I2C.Init(): %v", err)
	}

	// create a voltage register object
	reg := ft.I2C.Reg(slave, vAddr, ft232h.Addr8Bit, order)

	// initialize the voltage register pointer
	voltage, e := reg.Reader(2)
	if nil != e {
		log.Fatalf("reg.Reader(): %v", e)
	}

	// repeatedly dump the voltage register
	for {
		if v, e := voltage(false); nil != e {
			log.Fatalf("voltage(): %v", e)
		} else {
			log.Printf("voltage = %.1f mV", float32(v)*vLSB)
		}
	}
}
