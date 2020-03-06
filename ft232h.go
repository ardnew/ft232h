package ft232h

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
)

// FT232H is the primary type for interacting with the device, holding the USB
// device file descriptor configuration/status and individual communication
// interfaces.
// Open a connection with an FT232H by calling the NewFT232H constructor.
// If more than one FTDI device (any FTDI device, not just FT232H) is present on
// the system, there are several constructor variations of form NewFT232HWith*
// to help distinguish which device to open.
// The only interface that is initialized by default is GPIO. You must call an
// initialization method of one of the other interfaces before using it.
type FT232H struct {
	info *deviceInfo
	mode Mode
	I2C  *I2C
	SPI  *SPI
	GPIO *GPIO
}

// String constructs a string representation of an FT232H device.
func (m *FT232H) String() string {
	return fmt.Sprintf("{ Index: %s, Mode: %s, I2C: %+v, SPI: %+v, GPIO: %s }",
		m.info, m.mode, m.I2C, m.SPI, m.GPIO)
}

// NewFT232H attempts to open a connection with the first MPSSE-capable USB
// device found, returning a non-nil error if unsuccessful.
func NewFT232H() (*FT232H, error) {
	return NewFT232HWithMask(nil) // first device found
}

// NewFT232HWithIndex attempts to open a connection with the MPSSE-capable USB
// device enumerated at index (starting at 0), returning a non-nil error if
// unsuccessful.
func NewFT232HWithIndex(index int) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Index: fmt.Sprintf("%d", index)})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given vendor ID vid and product ID pid, returning a non-nil
// error if unsuccessful.
func NewFT232HWithVIDPID(vid uint16, pid uint16) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{
		VID: fmt.Sprintf("%d", vid),
		PID: fmt.Sprintf("%d", pid),
	})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given serial no., returning a non-nil error if unsuccessful.
func NewFT232HWithSerial(serial string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Serial: serial})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given description, returning a non-nil error if unsuccessful.
func NewFT232HWithDesc(desc string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Desc: desc})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device matching all of the given attributes, returning a non-nil error if
// unsuccessful. Uses the first device found if mask is nil or all attributes
// are empty strings.
//
// The attributes are each specified as strings, including the integers, so that
// any attribute not given (i.e. empty string) will never exclude a device. The
// integer attributes can be expressed in any base recognized by the Go grammar
// for numeric literals (e.g., "13", "0b1101", "0xD", and "D" are all valid and
// equivalent).
func NewFT232HWithMask(mask *OpenMask) (*FT232H, error) {
	m := &FT232H{info: nil, mode: ModeNone, I2C: nil, SPI: nil}
	if err := m.openDevice(mask); nil != err {
		return nil, err
	}
	m.I2C = &I2C{device: m, config: i2cConfigDefault()}
	m.SPI = &SPI{device: m, config: spiConfigDefault()}
	m.GPIO = &GPIO{device: m, config: GPIOConfigDefault()}
	if err := m.GPIO.Init(); nil != err {
		return nil, err
	}
	return m, nil
}

// OpenMask contains strings for each of the supported attributes used to
// distinguish which FTDI device to open. See NewFT232HWithMask for semantics.
type OpenMask struct {
	Index  string
	VID    string
	PID    string
	Serial string
	Desc   string
}

// parseUint32 attempts to convert a given string to a 32-bit unsigned integer,
// returning zero and false if the string is empty, negative, or otherwise
// invalid.
// The string can be expressed in various bases, following the convention of
// Go's strconv.ParseUint with base = 0, bitSize = 32.
// The only exception is when the string contains hexadecimal chars and doesn't
// begin with the required prefix "0x". In this case, the "0x" prefix is added
// automatically.
func parseUint32(s string) (uint32, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, false
	}
	// ParseUint requires a leading "0x" for base 16
	if strings.ContainsAny(s, "abcdef") && !strings.HasPrefix(s, "0b") {
		s = "0x" + strings.TrimPrefix(s, "0x") // always prefix (but not twice!)
	}
	// now parse according to Go convention
	i, err := strconv.ParseInt(s, 0, 64)
	if nil != err || i < 0 || i > math.MaxUint32 {
		return 0, false
	} else {
		return uint32(i), true
	}
}

// openDevice attempts to open the device matching the given mask, returning
// a non-nil error if unsuccessful. The error SDeviceNotFound is returned if
// no device was found matching the given mask. See NewFT232HWithMask for
// semantics.
func (m *FT232H) openDevice(mask *OpenMask) error {

	var (
		dev []*deviceInfo
		sel *deviceInfo
		err error
	)

	u32Eq := func(i uint32, s string) bool {
		if u, ok := parseUint32(s); ok {
			return i == u
		}
		return false
	}

	if dev, err = devices(); nil != err {
		return err
	}

	for _, d := range dev {
		if nil == mask {
			sel = d
			break
		}
		if "" != mask.Index {
			if !u32Eq(uint32(d.index), mask.Index) {
				continue
			}
		}
		if "" != mask.VID {
			if !u32Eq(d.vid, mask.VID) {
				continue
			}
		}
		if "" != mask.PID {
			if !u32Eq(d.pid, mask.PID) {
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

// Close closes the connection with an FT232H, returning a non-nil error if
// unsuccessful.
func (m *FT232H) Close() error {
	if nil != m.info {
		return m.info.close()
	}
	m.mode = ModeNone
	return nil
}

// Pin defines the methods required for representing an FT232H port pin.
type Pin interface {
	IsMPSSE() bool     // true if DPin (port "D"), false if CPin (GPIO/port "C")
	Mask() uint8       // the bitmask used to address the pin, equal to 1<<Pos()
	Pos() int          // the ordinal pin number (0-7), equal to log2(Mask())
	String() string    // the string representation "D#" or "C#", with # = Pos()
	Valid() bool       // true IFF bitmask has exactly one bit set
	Equals(q Pin) bool // true IFF p and q have equal port and bitmask
}

// IsMPSSE is true for pins on FT232H port "D".
func (p DPin) IsMPSSE() bool { return true }

// IsMPSSE is false for pins on FT232H port "C".
func (p CPin) IsMPSSE() bool { return false }

// Mask is the bitmask used to address the pin on port "D".
func (p DPin) Mask() uint8 { return uint8(p) }

// Mask is the bitmask used to address the pin on port "C".
func (p CPin) Mask() uint8 { return uint8(p) }

// Pos is the ordinal pin number (0-7) on port "D".
func (p DPin) Pos() int { return int(math.Log2(float64(p))) }

// Pos is the ordinal pin number (0-7) on port "C".
func (p CPin) Pos() int { return int(math.Log2(float64(p))) }

// String is the string representation "D#" of the pin, with # equal to Pos.
func (p DPin) String() string { return fmt.Sprintf("D%d", p.Pos()) }

// String is the string representation "C#" of the pin, with # equal to Pos.
func (p CPin) String() string { return fmt.Sprintf("C%d", p.Pos()) }

// Valid is true if the pin bitmask has exactly one bit set, otherwise false.
func (p DPin) Valid() bool { return 1 == bits.OnesCount64(uint64(p)) }

// Valid is true if the pin bitmask has exactly one bit set, otherwise false.
func (p CPin) Valid() bool { return 1 == bits.OnesCount64(uint64(p)) }

// Equals is true if the given pin is on port "D" and has the same bitmask,
// otherwise false.
func (p DPin) Equals(q Pin) bool { return q.IsMPSSE() && p.Mask() == q.Mask() }

// Equals is true if the given pin is on port "C" and has the same bitmask,
// otherwise false.
func (p CPin) Equals(q Pin) bool { return !q.IsMPSSE() && p.Mask() == q.Mask() }

// Types representing individual port pins.
type (
	DPin uint8 // pin bitmask on MPSSE low-byte lines (port "D" of FT232H)
	CPin uint8 // pin bitmask on MPSSE high-byte lines (port "C" of FT232H)
)

// Constants related to GPIO pin configuration
const (
	PinLO byte = 0 // pin value clear
	PinHI byte = 1 // pin value set
	PinIN byte = 0 // pin direction input
	PinOT byte = 1 // pin direction output

	NumDPins = 8 // number of MPSSE low-byte line pins
	NumCPins = 8 // number of MPSSE high-byte line pins
)

// D returns a DPin bitmask with only the given bit at position pin set.
// If the given pin position is negative or greater than 7, the invalid bitmask
// (0) is returned.
func D(pin int) DPin {
	if pin >= 0 && pin < NumDPins {
		return DPin(1 << pin)
	} else {
		return DPin(0) // invalid DPin
	}
}

// C returns a CPin bitmask with only the given bit at position pin set.
// If the given pin position is negative or greater than 7, the invalid bitmask
// (0) is returned.
func C(pin int) CPin {
	if pin >= 0 && pin < NumCPins {
		return CPin(1 << pin)
	} else {
		return CPin(0) // invalid CPin
	}
}

// deviceInfo contains the USB device descriptor and attributes for a device
// managed by the D2XX driver.
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

// String constructs a readable string representation of the deviceInfo.
func (dev *deviceInfo) String() string {
	return fmt.Sprintf("%d:{ Open = %t, HiSpeed = %t, Chip = \"%s\" (0x%02X), "+
		"VID = 0x%04X, PID = 0x%04X, Location = %04X, "+
		"Serial = \"%s\", Desc = \"%s\", Handle = %p }",
		dev.index+1, dev.isOpen, dev.isHiSpeed, dev.chip, uint32(dev.chip),
		dev.vid, dev.pid, dev.locID, dev.serial, dev.desc, dev.handle)
}

// open attempts to open a raw USB interface through the D2XX bridge, returning
// a non-nil error if unsuccessful.
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

// close attempts to close a USB interface opened through the D2XX bridge,
// returning a non-nil error if unsuccessful.
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

// devices queries all of the USB devices on the system using the D2XX bridge,
// returning a slice of deviceInfo pointers for all MPSSE-capable devices.
// Returns a nil slice and non-nil error if the driver failed to obtain device
// information from the system.
// Returns an empty slice and nil error if no MPSSE-capable devices were found
// after successful communication with the system.
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

// AddrSpace represents the address space of a pointer. Intended to be used when
// specifying e.g. I²C register addresses and the like.
type AddrSpace uint8

// Constants defining various address spaces.
const (
	Addr8Bit  AddrSpace = 1 << iota // 8-bit addresses
	Addr16Bit                       // 16-bit addresses
	Addr32Bit                       // 32-bit addresses
	Addr64Bit                       // 64-bit addresses
)

// String returns a string representation of the address space.
func (s AddrSpace) String() string {
	return fmt.Sprintf("%d-bit", s.Bits())
}

// Bits returns the number of usable bits in an address space.
func (s AddrSpace) Bits() uint {
	switch s {
	case Addr8Bit, Addr16Bit, Addr32Bit, Addr64Bit:
		return s.Bytes() * 8
	}
	return 0
}

// Bytes returns the number of usable bytes in an address space.
func (s AddrSpace) Bytes() uint {
	switch s {
	case Addr8Bit, Addr16Bit, Addr32Bit, Addr64Bit:
		return uint(s)
	}
	return 0
}

// ByteOrder represents the byte order of a sequence of bytes.
type ByteOrder uint8

// Constants defining the supported byte orderings.
const (
	MSB ByteOrder = iota // most significant byte first (big endian)
	LSB                  // least significant byte first (little endian)
)

// String returns a string representation of the byte order.
func (o ByteOrder) String() string {
	switch o {
	case MSB:
		return "MSB"
	case LSB:
		return "LSB"
	default:
		return "(invalid byte order)"
	}
}

// Bytes converts the given value to an ordered slice of bytes. The receiver
// value determines ordering, and the count argument defines slice length.
func (o ByteOrder) Bytes(count uint, value uint64) []uint8 {

	if count > 8 {
		count = 8
	}

	b := make([]uint8, count)
	for i := range b {
		switch o {
		case MSB:
			b[i] = uint8((value >> ((count - uint(i) - 1) * 8)) & 0xFF)
		case LSB:
			b[i] = uint8((value >> (count * 8)) & 0xFF)
		}
	}
	return b
}

// Uint converts a given slice of bytes to an unsigned integer. The receiver
// value determines ordering, and the count argument defines the number of bytes
// (starting from the beginning of the slice) to use for conversion. Cast the
// value returned if a narrower type is required.
func (o ByteOrder) Uint(count uint, bytes []uint8) uint64 {

	if nil == bytes {
		return 0
	}
	if count > 8 {
		count = 8
	}
	if count > uint(len(bytes)) {
		bytes = append(bytes, make([]uint8, count-uint(len(bytes)))...)
	}

	n := uint64(0)
	for i, b := range bytes[:count] {
		switch o {
		case MSB:
			n |= uint64(b) << ((count - uint(i) - 1) * 8)
		case LSB:
			n |= uint64(b) << (uint(i) * 8)
		}
	}
	return n
}
