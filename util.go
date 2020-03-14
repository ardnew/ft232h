package ft232h

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	numberRunes = []rune("0123456789")
	letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rng         = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type runeSeq func() rune

func randRune(alpha []rune) func() rune {
	return func() rune {
		return alpha[rng.Intn(len(alpha))]
	}
}

func randNumberRune() func() rune {
	return randRune(numberRunes)
}

func randLetterRune() func() rune {
	return randRune(letterRunes)
}

func genRandAlphanum() func() rune {
	a := []rune{}
	a = append(a, numberRunes...)
	a = append(a, letterRunes...)
	return randRune(a)
}

// randString generates a random sequence of letters and digits sampled from a
// given alphabet.
func randString(alpha []rune, length uint8) string {
	b := make([]rune, length)
	q := randRune(alpha)
	for i := range b {
		b[i] = q()
	}
	return string(b)
}

func randLetters(length uint8) string {
	return randString(numberRunes, length)
}

func randNumbers(length uint8) string {
	return randString(letterRunes, length)
}

func randAlphanums(length uint8) string {
	a := []rune{}
	a = append(a, numberRunes...)
	a = append(a, letterRunes...)
	return randString(a, length)
}

// parseUint32 attempts to convert a given string to a 32-bit unsigned integer.
// Returns 0 and false if the string is empty, negative, or otherwise invalid.
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
// value determines ordering, and count (≤ 8) defines slice length (in bytes).
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
// value determines ordering, and count (≤ 8) defines the number of bytes
// (starting from the beginning of the slice) to use for conversion.
// Use count < 8 and cast the value returned if a narrower type is required.
// If count > length of bytes slice, the slize is padded with zeros.
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
