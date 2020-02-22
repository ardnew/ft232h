package main

import (
	"log"

	se "github.com/ardnew/gompsse"
)

const (
	TFTCS  = se.D3
	TFTDC  = se.C0
	TFTRST = se.C4
)

func main() {

	for _, desc := range []string{"FT232H-C"} {
		m, err := se.NewMPSSEWithDesc(desc)
		if nil != err {
			log.Printf("NewMPSSEWithDesc(): %+v", err)
			continue
		}
		defer m.Close()

		lcd, lerr := NewILI9341(m, TFTCS, TFTDC, TFTRST, down)
		if lerr != nil {
			log.Printf("NewILI9341(): %+v", lerr)
			continue
		}

		if err := lcd.init(); nil != err {
			log.Printf("init(): %+v", err)
			continue
		}

		if err := lcd.fillScreen(0x100F); nil != err {
			log.Printf("fillScreen(): %+v", err)
			continue
		}

		log.Printf("%s", m)
	}
}
