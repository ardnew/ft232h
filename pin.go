package ft232h

import (
	"fmt"
	"math"
	"math/bits"
)

// Pin defines the methods required for representing an FT232H port pin.
type Pin interface {
	IsMPSSE() bool     // true if DPin (port "D"), false if CPin (GPIO/port "C")
	Mask() uint8       // the bitmask used to address the pin, equal to 1<<Pos()
	Pos() uint         // the ordinal pin number (0-7), equal to log2(Mask())
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
func (p DPin) Pos() uint { return uint(math.Log2(float64(p))) }

// Pos is the ordinal pin number (0-7) on port "C".
func (p CPin) Pos() uint { return uint(math.Log2(float64(p))) }

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

// Dir represents the direction of a GPIO pin
type Dir bool

// Constants of GPIO pin direction type Dir
const (
	Input  Dir = false // GPIO input pins (bit clear)
	Output Dir = true  // GPIO output pins (bit set)
)

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

	NumDPins = 8 // number of MPSSE low-byte line pins, port "D"
	NumCPins = 8 // number of MPSSE high-byte line pins, port "C"
)

// D returns a DPin bitmask with only the given bit at position pin set.
// If the given pin position is greater than 7, the invalid bitmask (0) is
// returned.
func D(pin uint) DPin {
	if pin >= 0 && pin < NumDPins {
		return DPin(1 << pin)
	} else {
		return DPin(0) // invalid DPin
	}
}

// C returns a CPin bitmask with only the given bit at position pin set.
// If the given pin position is greater than 7, the invalid bitmask (0) is
// returned.
func C(pin uint) CPin {
	if pin < NumCPins {
		return CPin(1 << pin)
	} else {
		return CPin(0) // invalid CPin
	}
}
