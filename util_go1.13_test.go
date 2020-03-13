// +build go1.13

package ft232h

import (
	"fmt"
	"testing"
)

func TestParseUint32_go1_13(t *testing.T) {

  // binary and octal literals not supported until go1.13
	for _, test := range []struct {
		str string
		exp uint32
		ok  bool // true if expected parse to succeed
	}{
		{
			str: "0b10101010101010101010101010101010",
			exp: 0b10101010101010101010101010101010,
			ok:  true,
		},
		{
			str: "0o25252525252",
			exp: 0o25252525252,
			ok:  true,
		},
		{
			str: "0b01010101010101010101010101010101",
			exp: 0b01010101010101010101010101010101,
			ok:  true,
		},
		{
			str: "0o12525252525",
			exp: 0o12525252525,
			ok:  true,
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
