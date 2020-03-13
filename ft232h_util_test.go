package ft232h

import (
	"fmt"
	// "log"
	"testing"
)

func TestParseUint32(t *testing.T) {

	// minGoVersion := func(v string) bool {
	// 	ok, err := validateGoVersion(">= " + v)
	// 	if !ok {
	// 		log.Printf("validate Go version: %+v", err)
	// 	}
	// 	return ok
	// }

	for _, test := range []struct {
		str string
		exp uint32
		ok  bool // true if expected parse to succeed
	}{
		// {
		// 	str: "0b10101010101010101010101010101010",
		// 	exp: 0xAAAAAAAA,
		// 	ok:  minGoVersion("1.13"),
		// },
		{
			str: "025252525252",
			exp: 0xAAAAAAAA,
			ok:  true,
		},
		// {
		// 	str: "0o25252525252",
		// 	exp: 0xAAAAAAAA,
		// 	ok:  minGoVersion("1.13"),
		// },
		{
			str: "2863311530",
			exp: 0xAAAAAAAA,
			ok:  true,
		},
		{
			str: "0xAaAaAaAa",
			exp: 0xAAAAAAAA,
			ok:  true,
		},
		{
			str: "AaAaAaAa",
			exp: 0xAAAAAAAA,
			ok:  true,
		},
		// {
		// 	str: "0b01010101010101010101010101010101",
		// 	exp: 0x55555555,
		// 	ok:  minGoVersion("1.13"),
		// },
		{
			str: "012525252525",
			exp: 0x55555555,
			ok:  true,
		},
		// {
		// 	str: "0o12525252525",
		// 	exp: 0x55555555,
		// 	ok:  minGoVersion("1.13"),
		// },
		{
			str: "1431655765",
			exp: 0x55555555,
			ok:  true,
		},
		{
			str: "0X55555555",
			exp: 0x55555555,
			ok:  true,
		},
		{
			str: "55555555",
			exp: 55555555,
			ok:  true,
		},
		{
			str: "-0b10101010",
			exp: 0,
			ok:  false,
		},
		{
			str: "-02525",
			exp: 0,
			ok:  false,
		},
		{
			str: "-0o2525",
			exp: 0,
			ok:  false,
		},
		{
			str: "-170",
			exp: 0,
			ok:  false,
		},
		{
			str: "-0xAa",
			exp: 0,
			ok:  false,
		},
		{
			str: "-Aa",
			exp: 0,
			ok:  false,
		},
		{
			str: "4294967296",
			exp: 0,
			ok:  false,
		},
	} {
		t.Run(fmt.Sprintf("%q", test.str),
			func(s *testing.T) {
				act, ok := parseUint32(test.str)
				if test.ok {
					if ok {
						if act == test.exp { // success
							// empty
						} else { // fail
							s.Fatalf("parsed uint32={%d} expected={%d} from string=%q",
								act, test.exp, test.str)
						}
					} else { // fail
						s.Fatalf("could not parse uint32 from string=%q", test.str)
					}
				} else {
					if ok { // fail
						s.Fatalf("parsed uint32={%d} from string=%q", act, test.str)
					} else { // success
						// empty
					}
				}
			})
	}
}

func TestAddrSpace(t *testing.T) {

	for _, test := range []struct {
		addr  AddrSpace
		bits  uint
		bytes uint
	}{
		{
			addr:  AddrSpace(0),
			bits:  0,
			bytes: 0,
		},
		{
			addr:  Addr8Bit,
			bits:  8,
			bytes: 1,
		},
		{
			addr:  Addr16Bit,
			bits:  16,
			bytes: 2,
		},
		{
			addr:  Addr32Bit,
			bits:  32,
			bytes: 4,
		},
		{
			addr:  Addr64Bit,
			bits:  64,
			bytes: 8,
		},
		{
			addr:  AddrSpace(1 << 4),
			bits:  0,
			bytes: 0,
		},
	} {
		t.Run(fmt.Sprintf("%d={0b%08b}", test.addr, test.addr),
			func(s *testing.T) {

				if test.addr.Bits() != test.bits {
					s.Fatalf("address space={%d} bits={%d}, expected={%d}",
						test.addr, test.addr.Bits(), test.bits)
				}

				if test.addr.Bytes() != test.bytes {
					s.Fatalf("address space={%d} bytes={%d}, expected={%d}",
						test.addr, test.addr.Bytes(), test.bytes)
				}
			})
	}

}
