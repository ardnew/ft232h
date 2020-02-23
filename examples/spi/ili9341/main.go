package main

import (
	"log"

	"github.com/ardnew/ft232h"
	"github.com/ardnew/ft232h/drv/ili9341"
)

func main() {

	var (
		ft  *ft232h.FT232H
		lcd *ili9341.ILI9341
		err error
	)

	for _, desc := range []string{"FT232H-C"} {

		ft, err = ft232h.NewFT232HWithDesc(desc)
		if nil != err {
			log.Printf("NewFT232HWithDesc(): %+v", err)
			continue
		}
		defer ft.Close()

		lcd, err = ili9341.InitILI9341(ft, &ili9341.Config{
			PinCS:  ft232h.D3,
			PinDC:  ft232h.C0,
			PinRST: ft232h.C4,
			Rotate: ili9341.RotDown,
		})
		if err != nil {
			log.Printf("InitILI9341(): %+v", err)
			continue
		}

		color := []ili9341.RGB{ili9341.Red, ili9341.Green, ili9341.Blue}
		index := 0
		for {
			if err := lcd.FillScreen(color[index]); nil != err {
				log.Printf("FillScreen(): %+v", err)
				continue
			}
			index = (index + 1) % len(color)
			//time.Sleep(100 * time.Millisecond)
		}

		log.Printf("%s", ft)
	}
}
