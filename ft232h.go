// Package ft232h provides a high-level interface to the FTDI FT232H USB to
// SPI/I²C/UART/GPIO protocol converter.
package ft232h

import (
	"fmt"
	"math"
	"strings"
)

type FT232H struct {
	info *deviceInfo
	mode Mode
	I2C  *I2C
	SPI  *SPI
	GPIO *GPIO
}

func (m *FT232H) String() string {
	return fmt.Sprintf("{ Info: %s, Mode: %s, I2C: %+v, SPI: %+v, GPIO: %+v }",
		m.info, m.mode, m.I2C, m.SPI, m.GPIO)
}

func NewFT232H() (*FT232H, error) {
	return NewFT232HWithMask(nil) // first device found
}

func NewFT232HWithIndex(index uint) (*FT232H, error) {
	s := fmt.Sprintf("%d", index)
	return NewFT232HWithMask(&OpenMask{Index: s})
}

func NewFT232HWithVIDPID(vid uint16, pid uint16) (*FT232H, error) {
	v := fmt.Sprintf("%04x", vid)
	p := fmt.Sprintf("%04x", pid)
	return NewFT232HWithMask(&OpenMask{VID: v, PID: p})
}

func NewFT232HWithSerial(serial string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Serial: serial})
}

func NewFT232HWithDesc(desc string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Desc: desc})
}

func NewFT232HWithMask(mask *OpenMask) (*FT232H, error) {
	m := &FT232H{info: nil, mode: ModeNone, I2C: nil, SPI: nil}
	if err := m.openDevice(mask); nil != err {
		return nil, err
	}
	m.I2C = &I2C{device: m, config: i2cConfigDefault()}
	m.SPI = &SPI{device: m, config: spiConfigDefault()}
	m.GPIO = &GPIO{device: m, config: gpioConfigDefault()}
	if err := m.GPIO.Init(); nil != err {
		return nil, err
	}
	return m, nil
}

type OpenMask struct {
	Index  string
	VID    string
	PID    string
	Serial string
	Desc   string
}

func (m *FT232H) openDevice(mask *OpenMask) error {

	var (
		dev []*deviceInfo
		sel *deviceInfo
		err error
	)

	if dev, err = devices(); nil != err {
		return err
	}

	for _, d := range dev {
		if nil == mask {
			sel = d
			break
		}
		if "" != mask.Index {
			if mask.Index != fmt.Sprintf("%d", d.index) {
				continue
			}
		}
		if "" != mask.VID {
			ms := strings.ToLower(mask.VID)
			dx := fmt.Sprintf("%x", d.vid)
			dz := fmt.Sprintf("%04x", d.vid)
			if (ms != dx) && (ms != ("0x" + dx)) &&
				(ms != dz) && (ms != ("0x" + dz)) &&
				(ms != fmt.Sprintf("%d", d.vid)) {
				continue
			}
		}
		if "" != mask.PID {
			ms := strings.ToLower(mask.PID)
			dx := fmt.Sprintf("%x", d.pid)
			dz := fmt.Sprintf("%04x", d.pid)
			if (ms != dx) && (ms != ("0x" + dx)) &&
				(ms != dz) && (ms != ("0x" + dz)) &&
				(ms != fmt.Sprintf("%d", d.pid)) {
				continue
			}
		}
		if "" != mask.Serial {
			if strings.ToLower(mask.Serial) != strings.ToLower(d.serial) {
				continue
			}
		}
		if "" != mask.Desc {
			if strings.ToLower(mask.Desc) != strings.ToLower(d.desc) {
				continue
			}
		}
		sel = d
		break
	}

	if nil == sel {
		return SDeviceNotFound
	}

	if err = sel.open(); nil != err {
		return err
	}
	m.info = sel
	return nil
}

func (m *FT232H) Close() error {
	if nil != m.info {
		return m.info.close()
	}
	m.mode = ModeNone
	return nil
}

type Pin interface {
	IsMPSSE() bool // true if DPin (port "D"), false if GPIO CPin (port "C")
	Mask() uint8
	Pos() uint8
	String() string
}

// Types representing individual port pins.
type (
	DPin uint8 // pin on MPSSE low-byte lines (port "D" on FT232H)
	CPin uint8 // pin on MPSSE high-byte lines (port "C" on FT232H)
)

func (p DPin) IsMPSSE() bool  { return true }
func (p CPin) IsMPSSE() bool  { return false }
func (p DPin) Mask() uint8    { return uint8(p) }
func (p CPin) Mask() uint8    { return uint8(p) }
func (p DPin) Pos() uint8     { return uint8(math.Log2(float64(p))) }
func (p CPin) Pos() uint8     { return uint8(math.Log2(float64(p))) }
func (p DPin) String() string { return fmt.Sprintf("D%d", p.Pos()) }
func (p CPin) String() string { return fmt.Sprintf("C%d", p.Pos()) }

// Constants related to GPIO pin configuration
const (
	PinLO byte = 0 // pin value clear
	PinHI byte = 1 // pin value set
	PinIN byte = 0 // pin direction input
	PinOT byte = 1 // pin direction output

	NumDPins = 8 // number of MPSSE low-byte line pins
	NumCPins = 8 // number of MPSSE high-byte line pins
)

// Constants defining the available board pins on MPSSE low-byte lines
const (
	D0 DPin = 1 << iota
	D1
	D2
	D3
	D4
	D5
	D6
	D7
)

// Constants defining the available board pins on MPSSE high-byte lines
const (
	C0 CPin = 1 << iota
	C1
	C2
	C3
	C4
	C5
	C6
	C7
)

type deviceInfo struct {
	index     int
	isOpen    bool
	isHiSpeed bool
	chip      Chip
	vid       uint32
	pid       uint32
	locID     uint32
	serial    string
	desc      string
	handle    Handle
}

func (dev *deviceInfo) String() string {
	return fmt.Sprintf("%d:{ Open = %t, HiSpeed = %t, Chip = \"%s\" (0x%02X), "+
		"VID = 0x%04X, PID = 0x%04X, Location = %04X, "+
		"Serial = \"%s\", Desc = \"%s\", Handle = %p }",
		dev.index, dev.isOpen, dev.isHiSpeed, dev.chip, uint32(dev.chip),
		dev.vid, dev.pid, dev.locID, dev.serial, dev.desc, dev.handle)
}

func (dev *deviceInfo) open() error {
	if ce := dev.close(); nil != ce {
		return ce
	}
	if oe := _FT_Open(dev); nil != oe {
		return oe
	}
	dev.isOpen = true
	return nil
}

func (dev *deviceInfo) close() error {
	if !dev.isOpen {
		return nil
	}
	if ce := _FT_Close(dev); nil != ce {
		return ce
	}
	dev.isOpen = false
	return nil
}

func devices() ([]*deviceInfo, error) {

	n, ce := _FT_CreateDeviceInfoList()
	if nil != ce {
		return nil, ce
	}

	if 0 == n {
		return []*deviceInfo{}, nil
	}

	info, de := _FT_GetDeviceInfoList(n)
	if nil != de {
		return nil, de
	}

	return info, nil
}
