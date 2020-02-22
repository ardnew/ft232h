package main

import (
	"time"

	se "github.com/ardnew/gompsse"
)

type rotation byte

const (
	down  rotation = 0
	right rotation = 1
	up    rotation = 2
	left  rotation = 3
)

func (r rotation) data() byte {
	switch r {
	case down:
		return 0x40 | 0x08
	case right:
		return 0x40 | 0x80 | 0x20 | 0x08
	case up:
		return 0x80 | 0x08
	case left:
		return 0x20 | 0x08
	default:
		return 0x00
	}
}

type ili9341 struct {
	mpsse  *se.MPSSE
	pinCS  se.DPin
	pinDC  se.CPin
	pinRST se.CPin
	rot    rotation
}

func NewILI9341(m *se.MPSSE, cs se.DPin, dc se.CPin, rst se.CPin, rot rotation) (*ili9341, error) {

	// configure at 30MHz, 1ms latency, CS active low, SPI mode 0
	if err := m.SPI.SetConfig(30000000, 1, cs, true, 0); nil != err {
		return nil, err
	}
	if err := m.SPI.Init(); nil != err {
		return nil, err
	}

	return &ili9341{
		mpsse:  m,
		pinCS:  cs,
		pinDC:  dc,
		pinRST: rst,
		rot:    rot,
	}, nil
}

func (lcd *ili9341) init() error {

	// command list is based on https://github.com/martnak/STM32-ILI9341

	if err := lcd.mpsse.GPIO.Set(lcd.pinRST, false); nil != err {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	if err := lcd.mpsse.GPIO.Set(lcd.pinRST, true); nil != err {
		return err
	}

	// SOFTWARE RESET
	if err := lcd.command(0x01); nil != err {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	// POWER CONTROL A
	if err := lcd.commandData(0xCB, []uint8{0x39, 0x2C, 0x00, 0x34, 0x02}); nil != err {
		return err
	}

	// POWER CONTROL B
	if err := lcd.commandData(0xCF, []uint8{0x00, 0xC1, 0x30}); nil != err {
		return err
	}

	// DRIVER TIMING CONTROL A
	if err := lcd.commandData(0xE8, []uint8{0x85, 0x00, 0x78}); nil != err {
		return err
	}

	// DRIVER TIMING CONTROL B
	if err := lcd.commandData(0xEA, []uint8{0x00, 0x00}); nil != err {
		return err
	}

	// POWER ON SEQUENCE CONTROL
	if err := lcd.commandData(0xED, []uint8{0x64, 0x03, 0x12, 0x81}); nil != err {
		return err
	}

	// PUMP RATIO CONTROL
	if err := lcd.commandData(0xF7, []uint8{0x20}); nil != err {
		return err
	}

	// POWER CONTROL,VRH[5:0]
	if err := lcd.commandData(0xC0, []uint8{0x23}); nil != err {
		return err
	}

	// POWER CONTROL,SAP[2:0];BT[3:0]
	if err := lcd.commandData(0xC1, []uint8{0x10}); nil != err {
		return err
	}

	// VCM CONTROL
	if err := lcd.commandData(0xC5, []uint8{0x3E, 0x28}); nil != err {
		return err
	}

	// VCM CONTROL 2
	if err := lcd.commandData(0xC7, []uint8{0x86}); nil != err {
		return err
	}

	// MEMORY ACCESS CONTROL
	if err := lcd.commandData(0x36, []uint8{0x48}); nil != err {
		return err
	}

	// PIXEL FORMAT
	if err := lcd.commandData(0x3A, []uint8{0x55}); nil != err {
		return err
	}

	// FRAME RATIO CONTROL, STANDARD RGB COLOR
	if err := lcd.commandData(0xB1, []uint8{0x00, 0x18}); nil != err {
		return err
	}

	// DISPLAY FUNCTION CONTROL
	if err := lcd.commandData(0xB6, []uint8{0x08, 0x82, 0x27}); nil != err {
		return err
	}

	// 3GAMMA FUNCTION DISABLE
	if err := lcd.commandData(0xF2, []uint8{0x00}); nil != err {
		return err
	}

	// GAMMA CURVE SELECTED
	if err := lcd.commandData(0x26, []uint8{0x01}); nil != err {
		return err
	}

	// POSITIVE GAMMA CORRECTION
	if err := lcd.commandData(0xE0, []uint8{0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08,
		0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00}); nil != err {
		return err
	}

	// NEGATIVE GAMMA CORRECTION
	if err := lcd.commandData(0xE1, []uint8{0x00, 0x0E, 0x14, 0x03, 0x11, 0x07,
		0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F}); nil != err {
		return err
	}

	// EXIT SLEEP
	if err := lcd.command(0x11); nil != err {
		return err
	}

	time.Sleep(120 * time.Millisecond)

	// TURN ON DISPLAY
	if err := lcd.command(0x29); nil != err {
		return err
	}

	// MADCTL
	if err := lcd.commandData(0x36, []uint8{lcd.rot.data()}); nil != err {
		return err
	}

	return nil
}

func (lcd *ili9341) command(cmd uint8) error {

	if err := lcd.mpsse.GPIO.Set(lcd.pinDC, false); nil != err {
		return err
	}
	if _, err := lcd.mpsse.SPI.Write([]uint8{cmd}, true, true); nil != err {
		return err
	}
	return nil
}

func (lcd *ili9341) data(data []uint8) error {

	if err := lcd.mpsse.GPIO.Set(lcd.pinDC, true); nil != err {
		return err
	}
	if _, err := lcd.mpsse.SPI.Write(data, true, true); nil != err {
		return err
	}
	return nil
}

func (lcd *ili9341) commandData(cmd uint8, data []uint8) error {
	if err := lcd.command(cmd); nil != err {
		return err
	}
	if err := lcd.data(data); nil != err {
		return err
	}
	return nil
}

func (lcd *ili9341) window(x0 uint16, y0 uint16, x1 uint16, y1 uint16) error {

	// column address set
	if err := lcd.commandData(0x2A, []uint8{uint8((x0 >> 8) & 0xFF), uint8(x0 & 0xFF),
		uint8((x1 >> 8) & 0xFF), uint8(x1 & 0xFF)}); nil != err {
		return err
	}

	// row address set
	if err := lcd.commandData(0x2B, []uint8{uint8((y0 >> 8) & 0xFF), uint8(y0 & 0xFF),
		uint8((y1 >> 8) & 0xFF), uint8(y1 & 0xFF)}); nil != err {
		return err
	}

	// write to RAM
	if err := lcd.command(0x2C); nil != err {
		return err
	}

	return nil
}

func (lcd *ili9341) fillScreen(color uint16) error {

	if err := lcd.window(0, 0, 240, 320); nil != err {
		return err
	}

	numPixels := 240 * 320
	block := make([]uint8, 2*numPixels)
	for i := 0; i < numPixels; i++ {
		block[2*i] = uint8((color >> 8) & 0xFF)
		block[2*i+1] = uint8((color >> 0) & 0xFF)
	}

	if err := lcd.data(block); nil != err {
		return err
	}

	return nil
}
