package main

import (
	"log"
	"time"

	"github.com/ardnew/ft232h"
)

func main() {

	ft, err := ft232h.NewFT232H()
	if nil != err {
		log.Fatalf("NewFT232H(): %s", err)
	}
	defer ft.Close()

	ft.GPIO.Config(&ft232h.GPIOConfig{Dir: 0x55, Val: 0x00})

	pin := ft232h.C(6)
	setAll := false
	val := uint8(0)

	for {
		if setAll {
			if err := ft.GPIO.Write(val); nil != err {
				log.Printf("GPIO.Write(): %s", err)
				break
			}
			val = ^val
		} else {
			if err := ft.GPIO.Set(pin, (val&pin.Mask()) > 0); nil != err {
				log.Printf("GPIO.Set(): %s", err)
				break
			}
		}
		log.Printf("GPIO: %s", ft.GPIO)
		setAll = !setAll
		time.Sleep(time.Second)
	}
}
