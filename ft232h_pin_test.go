package ft232h

import "testing"

func TestPin(t *testing.T) {

	for i := -1; i <= 8; i++ {

		ok := i > -1 && i < 8

		d := D(uint(i))
		c := C(uint(i))

		if ok {

			if !d.Valid() {
				t.Fatalf("expected D(%d) to be valid", i)
			}
			if !c.Valid() {
				t.Fatalf("expected C(%d) to be valid", i)
			}

			if (1 << i) != d.Mask() {
				t.Fatalf("D(%d) mask={%08b}, expected={%08b}", i, d.Mask(), 1<<i)
			}
			if (1 << i) != c.Mask() {
				t.Fatalf("C(%d) mask={%08b}, expected={%08b}", i, c.Mask(), 1<<i)
			}

			if i != int(d.Pos()) {
				t.Fatalf("D(%d) pos={%d}, expected={%d}", i, d.Pos(), i)
			}
			if i != int(c.Pos()) {
				t.Fatalf("C(%d) pos={%d}, expected={%d}", i, c.Pos(), i)
			}

			if d.Equals(c) || c.Equals(d) {
				t.Fatalf("D(%d) equals C(%d)", i, i)
			}

		} else {
			if d.Valid() {
				t.Fatalf("expected D(%d) to be invalid", i)
			}
			if c.Valid() {
				t.Fatalf("expected C(%d) to be invalid", i)
			}
		}
	}
}
