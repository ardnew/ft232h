package ili9341

import (
	"fmt"
	"time"

	"github.com/ardnew/ft232h"
)

type ILI9341 struct {
	device *ft232h.FT232H
	config *Config
}

type Config struct {
	PinCS  ft232h.Pin
	PinDC  ft232h.CPin
	PinRST ft232h.CPin
	Rotate Rotation
}

func NewILI9341(ft *ft232h.FT232H, config *Config) (*ILI9341, error) {

	if 0 == config.PinCS.Mask() {
		return nil, fmt.Errorf("chip-select pin not provided")
	}

	if 0 == config.PinDC.Mask() {
		return nil, fmt.Errorf("data/command pin not provided")
	}

	c := &ft232h.SPIConfig{
		SPIOption: &ft232h.SPIOption{
			CS:        config.PinCS,
			ActiveLow: true,
			Mode:      0,
		},
		Clock:   30000000, // 30 MHz
		Latency: ft232h.SPILatencyDefault,
	}

	if err := ft.SPI.Config(c); nil != err {
		return nil, err
	}

	lcd := &ILI9341{device: ft, config: config}

	if err := lcd.Init(); nil != err {
		return nil, err
	}

	return lcd, nil
}

type Point struct {
	X int
	Y int
}

func MakePoint(x int, y int) Point {
	return Point{X: x, Y: y}
}

type Size struct {
	Width  int
	Height int
}

func MakeSize(w int, h int) Size {
	return Size{Width: w, Height: h}
}

type Frame struct {
	Origin Point
	Size   Size
}

func MakeFrame(x int, y int, w int, h int) Frame {
	return Frame{
		Origin: Point{X: x, Y: y},
		Size:   Size{Width: w, Height: h},
	}
}

func MakeFrameRect(x0 int, y0 int, x1 int, y1 int) Frame {
	return MakeFrame(x0, y0, x1-x0, y1-y0)
}

func (f *Frame) colAddress() []uint8 {
	x0 := f.Origin.X
	x1 := f.Origin.X + f.Size.Width
	return []uint8{
		uint8((x0 >> 8) & 0xFF),
		uint8(x0 & 0xFF),
		uint8((x1 >> 8) & 0xFF),
		uint8(x1 & 0xFF),
	}
}

func (f *Frame) rowAddress() []uint8 {
	y0 := f.Origin.Y
	y1 := f.Origin.Y + f.Size.Height
	return []uint8{
		uint8((y0 >> 8) & 0xFF),
		uint8(y0 & 0xFF),
		uint8((y1 >> 8) & 0xFF),
		uint8(y1 & 0xFF),
	}
}

type Rotation byte

const (
	// rotation indicates position of board pins when looking at the screen
	RotDown    Rotation = 0
	RotLeft    Rotation = 1
	RotUp      Rotation = 2
	RotRight   Rotation = 3
	RotDefault Rotation = RotDown
)

func (r Rotation) MADCTL() byte {
	switch r {
	case RotDown:
		return 0x40 | 0x08
	case RotLeft:
		return 0x40 | 0x80 | 0x20 | 0x08
	case RotUp:
		return 0x80 | 0x08
	case RotRight:
		return 0x20 | 0x08
	default:
		return RotDefault.MADCTL()
	}
}

const (
	numHorzPixels = 240
	numVertPixels = 320
	NumPixels     = numHorzPixels * numVertPixels
)

func (r Rotation) Size() Size {
	switch r {
	case RotDown:
		return Size{Width: numHorzPixels, Height: numVertPixels}
	case RotLeft:
		return Size{Width: numVertPixels, Height: numHorzPixels}
	case RotUp:
		return Size{Width: numHorzPixels, Height: numVertPixels}
	case RotRight:
		return Size{Width: numVertPixels, Height: numHorzPixels}
	default:
		return RotDefault.Size()
	}
}

type RGB16 uint16 // RGB 5-6-5 format

type RGB struct {
	R int16
	G int16
	B int16
}

var (
	Black       = RGB{R: 0x00, G: 0x00, B: 0x00}
	Navy        = RGB{R: 0x00, G: 0x00, B: 0x0F}
	DarkGreen   = RGB{R: 0x00, G: 0x1F, B: 0x00}
	DarkCyan    = RGB{R: 0x00, G: 0x1F, B: 0x0F}
	Maroon      = RGB{R: 0x0F, G: 0x00, B: 0x00}
	Purple      = RGB{R: 0x0F, G: 0x00, B: 0x0F}
	Olive       = RGB{R: 0x0F, G: 0x1F, B: 0x00}
	LightGrey   = RGB{R: 0x18, G: 0x30, B: 0x18}
	DarkGrey    = RGB{R: 0x0F, G: 0x1F, B: 0x0F}
	Blue        = RGB{R: 0x00, G: 0x00, B: 0x1F}
	Green       = RGB{R: 0x00, G: 0x3F, B: 0x00}
	Cyan        = RGB{R: 0x00, G: 0x3F, B: 0x1F}
	Red         = RGB{R: 0x1F, G: 0x00, B: 0x00}
	Magenta     = RGB{R: 0x1F, G: 0x00, B: 0x1F}
	Yellow      = RGB{R: 0x1F, G: 0x3F, B: 0x00}
	White       = RGB{R: 0x1F, G: 0x3F, B: 0x1F}
	Orange      = RGB{R: 0x1F, G: 0x29, B: 0x00}
	GreenYellow = RGB{R: 0x15, G: 0x3F, B: 0x05}
	Pink        = RGB{R: 0x1F, G: 0x00, B: 0x1F}
)

func (c RGB16) Unpack() RGB {
	return RGB{
		R: int16((uint16(c) >> 11) & 0x1F),
		G: int16((uint16(c) >> 5) & 0x3F),
		B: int16((uint16(c) >> 0) & 0x1F),
	}
}

func (c *RGB) Pack() RGB16 {
	return RGB16(((uint16(c.R) & 0x1F) << 11) |
		((uint16(c.G) & 0x3F) << 5) |
		(uint16(c.B) & 0x1F))
}

func (c RGB16) MSB() uint8 {
	return uint8((c >> 8) & 0xFF)
}

func (c *RGB) MSB() uint8 {
	return c.Pack().MSB()
}

func (c RGB16) LSB() uint8 {
	return uint8(c & 0xFF)
}

func (c *RGB) LSB() uint8 {
	return c.Pack().LSB()
}

func (c RGB16) Buffer(n uint) []uint8 {
	// construct a buffer of 16-bit color (ordered MSB-first, suitable for passing
	// to ILI9341.SendData()). returned slice will have length 2*n bytes.
	buff := make([]uint8, 2*n)
	msb := c.MSB()
	lsb := c.LSB()
	for i := uint(0); i < 2*n; i += 2 {
		buff[i] = msb
		buff[i+1] = lsb
	}
	return buff
}

func (c RGB) Buffer(n uint) []uint8 {
	return c.Pack().Buffer(n)
}

func Wheel() func() RGB {
	var pos uint8
	return func() RGB {
		p := 0xFF - pos
		pos++
		if p < 0x55 {
			return RGB{int16(p * 0x03), int16(0xFF - p*0x03), int16(0x00)}
		} else if p < 0xAA {
			p -= 0x55
			return RGB{int16(0xFF - p*0x03), int16(0x00), int16(p * 0x03)}
		} else {
			p -= 0xAA
			return RGB{int16(0x00), int16(p * 0x03), int16(0xFF - p*0x03)}
		}
	}
}

func (lcd *ILI9341) setPinRST(set bool) error {
	if 0 == lcd.config.PinRST.Mask() {
		return fmt.Errorf("reset pin undefined")
	}
	if err := lcd.device.GPIO.Set(lcd.config.PinRST, set); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) setPinDC(set bool) error {
	if 0 == lcd.config.PinDC.Mask() {
		return fmt.Errorf("data/command pin undefined")
	}
	// clear DC line to indicate command on MOSI, set DC to indicate data
	if err := lcd.device.GPIO.Set(lcd.config.PinDC, set); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) Reset() error {
	// hardware reset
	if err := lcd.setPinRST(false); nil != err {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	if err := lcd.setPinRST(true); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) SendCommand(cmd uint8) error {
	// clear DC line to indicate command on MOSI
	if err := lcd.setPinDC(false); nil != err {
		return err
	}
	// write command using auto CS-assertion
	if _, err := lcd.device.SPI.Write([]uint8{cmd}, true, true); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) SendData(data []uint8) error {
	// write data using CS auto-assertion.
	if err := lcd.WriteData(data, true, true); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) WriteData(data []uint8, start bool, stop bool) error {
	// only assert DC line if we are starting a transfer
	if start {
		// set DC line to indicate data on MOSI
		if err := lcd.setPinDC(true); nil != err {
			return err
		}
	}

	// write data using optional CS auto-assertion. if the start flag is not true,
	// the CS line will not be asserted before starting transfer, and if stop flag
	// is not true, the line will not be de-asserted after transfer. this is used
	// in case your writes need to be broken up across multiple calls.
	if _, err := lcd.device.SPI.Write(data, start, stop); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) SendCommandData(cmd uint8, data []uint8) error {
	if err := lcd.SendCommand(cmd); nil != err {
		return err
	}
	if err := lcd.SendData(data); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) Init() error {

	if err := lcd.Reset(); nil != err {
		return err
	}

	for _, s := range []struct {
		cmd   uint8
		data  []uint8
		delay time.Duration
	}{
		{0x01, nil, 200 * time.Millisecond},              // software reset
		{0xCB, []uint8{0x39, 0x2C, 0x00, 0x34, 0x02}, 0}, // power control A
		{0xCF, []uint8{0x00, 0xC1, 0x30}, 0},             // power control B
		{0xE8, []uint8{0x85, 0x00, 0x78}, 0},             // driver timing control A
		{0xEA, []uint8{0x00, 0x00}, 0},                   // driver timing control B
		{0xED, []uint8{0x64, 0x03, 0x12, 0x81}, 0},       // power-on sequence control
		{0xF7, []uint8{0x20}, 0},                         // pump ratio control
		{0xC0, []uint8{0x23}, 0},                         // power control VRH[5:0]
		{0xC1, []uint8{0x10}, 0},                         // power control SAP[2:0] BT[3:0]
		{0xC5, []uint8{0x3E, 0x28}, 0},                   // VCM control
		{0xC7, []uint8{0x86}, 0},                         // VCM control 2
		{0x36, []uint8{0x48}, 0},                         // memory access control
		{0x3A, []uint8{0x55}, 0},                         // pixel format
		{0xB1, []uint8{0x00, 0x18}, 0},                   // frame ratio control, standard RGB color
		{0xB6, []uint8{0x08, 0x82, 0x27}, 0},             // display function control
		{0xF2, []uint8{0x00}, 0},                         // 3-Gamma function disable
		{0x26, []uint8{0x01}, 0},                         // gamma curve selected
		{0xE0, []uint8{ // positive gamma correction
			0x0F, 0x31, 0x2B, 0x0C, 0x0E,
			0x08, 0x4E, 0xF1, 0x37, 0x07,
			0x10, 0x03, 0x0E, 0x09, 0x00}, 0},
		{0xE1, []uint8{ // negative gamma correction
			0x00, 0x0E, 0x14, 0x03, 0x11,
			0x07, 0x31, 0xC1, 0x48, 0x08,
			0x0F, 0x0C, 0x31, 0x36, 0x0F}, 0},
		{0x11, nil, 120 * time.Millisecond},            // exit sleep
		{0x29, nil, 0},                                 // turn on display
		{0x36, []uint8{lcd.config.Rotate.MADCTL()}, 0}, // MADCTL
	} {
		if err := func() error {
			if nil == s.data {
				return lcd.SendCommand(s.cmd)
			} else {
				return lcd.SendCommandData(s.cmd, s.data)
			}
		}(); nil != err {
			return err
		}
		if s.delay > 0 {
			time.Sleep(s.delay)
		}
	}

	return nil
}

func (lcd *ILI9341) SetFrame(frame Frame) error {

	// column-, row-address set, write to RAM
	var caset, raset, ramwr uint8 = 0x2A, 0x2B, 0x2C

	vis := lcd.Normalize(frame)

	if err := lcd.SendCommandData(caset, vis.colAddress()); nil != err {
		return err
	}
	if err := lcd.SendCommandData(raset, vis.rowAddress()); nil != err {
		return err
	}
	if err := lcd.SendCommand(ramwr); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) SetFrameRect(x0 int, y0 int, x1 int, y1 int) error {
	return lcd.SetFrame(MakeFrameRect(x0, y0, x1, y1))
}

func (lcd *ILI9341) Size() Size {
	return lcd.config.Rotate.Size()
}

func (lcd *ILI9341) Clip(p Point) Point {
	s := lcd.config.Rotate.Size()
	if p.X < 0 {
		p.X = 0
	} else {
		if p.X > s.Width {
			p.X = s.Width
		}
	}
	if p.Y < 0 {
		p.Y = 0
	} else {
		if p.Y > s.Height {
			p.Y = s.Height
		}
	}
	return p
}

func (lcd *ILI9341) Normalize(f Frame) Frame {

	// if size is negative in either dimension, adjust the origin to exist on the
	// lesser axis, and flip the sign of the size dimension
	if f.Size.Width < 0 {
		f.Origin.X += f.Size.Width
		f.Size.Width = -f.Size.Width
	}
	if f.Size.Height < 0 {
		f.Origin.Y += f.Size.Height
		f.Size.Height = -f.Size.Height
	}

	// ---------------------------------------------------------------------------
	// if origin is less than screen bounds, it will be clipped. adjust the size
	// by the amount being clipped so that the greater axis stays in place. if
	// origin is greater than screen bounds, then our resulting frame will have
	// zero units in that dimension after clipping, so it doesn't matter. this is
	// all because we can guarantee size is a positive value, >0 per above.
	// ---------------------------------------------------------------------------

	// this will temporarily malform our rectangle, but is necessary. the result
	// after clipping below will yield the correct region.
	if f.Origin.X < 0 {
		f.Size.Width += f.Origin.X
	}
	if f.Origin.Y < 0 {
		f.Size.Height += f.Origin.Y
	}

	// now make sure both points exist in the visible screen area
	p1 := lcd.Clip(f.Origin)
	p2 := lcd.Clip(Point{
		X: p1.X + f.Size.Width,
		Y: p1.Y + f.Size.Height,
	})

	// reconstruct the resultive visible frame
	return MakeFrameRect(p1.X, p1.Y, p2.X, p2.Y)
}

func (lcd *ILI9341) FillScreen(color RGB) error {

	sz := lcd.config.Rotate.Size()
	if err := lcd.SetFrame(MakeFrame(0, 0, sz.Width, sz.Height)); nil != err {
		return err
	}
	if err := lcd.SendData(color.Buffer(NumPixels)); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) FillFrame(color RGB, frame Frame) error {

	fr := lcd.Normalize(frame)
	if 0 == fr.Size.Width || 0 == fr.Size.Height {
		return nil
	}
	if err := lcd.SetFrame(fr); nil != err {
		return err
	}
	px := fr.Size.Width * fr.Size.Height
	if err := lcd.SendData(color.Buffer(uint(px))); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) FillFrameRect(color RGB, x int, y int, w int, h int) error {
	return lcd.FillFrame(color, MakeFrame(x, y, w, h))
}

func (lcd *ILI9341) DrawPixel(color RGB, x int, y int) error {

	pt := lcd.Clip(MakePoint(x, y))
	if err := lcd.SetFrame(MakeFrame(pt.X, pt.Y, 1, 1)); nil != err {
		return err
	}
	if err := lcd.SendData(color.Buffer(1)); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) DrawBitmap1BPP(fg RGB, bg RGB, frame Frame, bmp []uint8) error {

	fr := lcd.Normalize(frame)
	if 0 == fr.Size.Width || 0 == fr.Size.Height {
		return nil
	}

	w, h := fr.Size.Width, fr.Size.Height

	fm, fl, bm, bl :=
		fg.MSB(), fg.LSB(), bg.MSB(), bg.LSB()

	wordWidth, bit := (w+7)/8, uint8(0)
	data, n := make([]uint8, 2*w*h), 0

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if (i & 7) > 0 {
				bit <<= 1
			} else {
				bit = bmp[j*wordWidth+i/8]
			}
			if (bit & 0x80) > 0 {
				data[n] = fm
				data[n+1] = fl
			} else {
				data[n] = bm
				data[n+1] = bl
			}
			n += 2
		}
	}
	if err := lcd.SetFrame(fr); nil != err {
		return err
	}
	if err := lcd.SendData(data); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) DrawBitmapRect1BPP(fg RGB, bg RGB, x int, y int, w int, h int, bmp []uint8) error {
	return lcd.DrawBitmap1BPP(fg, bg, MakeFrame(x, y, w, h), bmp)
}

func (lcd *ILI9341) DrawBitmap16BPP(frame Frame, bmp []uint16) error {

	fr := lcd.Normalize(frame)
	if 0 == fr.Size.Width || 0 == fr.Size.Height {
		return nil
	}

	w, h := fr.Size.Width, fr.Size.Height

	numPx := w * h
	if numPx > len(bmp) {
		return fmt.Errorf("not enough data to fill drawing area")
	}

	// re-order color data, MSB-first
	data := make([]uint8, 2*numPx)
	for i, rgb := range bmp[:numPx] {
		data[2*i] = uint8(rgb>>8) & 0xFF
		data[2*i+1] = uint8(rgb>>0) & 0xFF
	}

	if err := lcd.SetFrame(fr); nil != err {
		return err
	}
	if err := lcd.SendData(data); nil != err {
		return err
	}
	return nil
}

func (lcd *ILI9341) DrawBitmapRect16BPP(x int, y int, w int, h int, bmp []uint16) error {
	return lcd.DrawBitmap16BPP(MakeFrame(x, y, w, h), bmp)
}
