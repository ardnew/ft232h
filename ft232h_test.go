package ft232h

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	defIndex  int    = 0
	defVID    int    = 0x0403
	defPID    int    = 0x6014
	defSerial string = ""
	defDesc   string = ""

	alnum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	optIndex  int
	optVID    int
	optPID    int
	optSerial string
	optDesc   string

	rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func randString(length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = alnum[rng.Intn(len(alnum))]
	}
	return string(b)
}

func TestMain(m *testing.M) {

	flag.IntVar(&optIndex, "index", defIndex, "use device enumerated at `index`")
	flag.IntVar(&optVID, "vid", defVID, "use device with vendor ID `vid`")
	flag.IntVar(&optPID, "pid", defPID, "use device with product ID `pid`")
	flag.StringVar(&optSerial, "serial", defSerial, "use device with identifier `serial`")
	flag.StringVar(&optDesc, "desc", defDesc, "use device with description `desc`")

	flag.Parse()
	os.Exit(m.Run())
}

func TestNewFT232H(t *testing.T) {

	var (
		ft  *FT232H
		err error
	)

	ft, err = NewFT232H()
	if nil != err {
		t.Fatalf("could not open device: %v", err)
	}

	err = ft.Close()
	if nil != err {
		t.Fatalf("could not close device: %v", err)
	}
}

func TestNewFT232HWithIndex(t *testing.T) {

	for _, test := range []struct {
		name  string
		index int
		open  bool // true if expecting device to open successfully
	}{
		{name: "selected", index: optIndex, open: true},
		{name: "invalid", index: 0xBADC0DE, open: false},
	} {
		t.Run(fmt.Sprintf("%s={%d}", test.name, test.index),
			func(s *testing.T) {
				ft, err := NewFT232HWithIndex(test.index)
				if test.open {
					if nil == err { // success
						if err := ft.Close(); nil != err {
							s.Fatalf("could not close device at index=%d: %v", test.index, err)
						}
					} else { // fail
						s.Fatalf("could not open device at index=%d: %v", test.index, err)
					}
				} else {
					if nil == err { // fail
						ft.Close()
						s.Fatalf("opened device at index=%d", test.index)
					} else { // success
						// empty
					}
				}
			})
	}
}

func TestNewFT232HWithVIDPID(t *testing.T) {

	for _, test := range []struct {
		name string
		vid  int
		pid  int
		open bool // true if expecting device to open successfully
	}{
		{name: "selected", vid: optVID, pid: optPID, open: true},
		{name: "invalidVID", vid: 0xBAD, pid: optPID, open: false},
		{name: "invalidPID", vid: optVID, pid: 0xC0DE, open: false},
		{name: "invalid", vid: 0xBAD, pid: 0xC0DE, open: false},
	} {
		t.Run(fmt.Sprintf("%s={0x%X,0x%X}", test.name, test.vid, test.pid),
			func(s *testing.T) {
				ft, err := NewFT232HWithVIDPID(uint16(test.vid), uint16(test.pid))
				if test.open {
					if nil == err { // success
						if err := ft.Close(); nil != err {
							s.Fatalf("could not close device with VID=0x%X, PID=0x%X: %v",
								test.vid, test.pid, err)
						}
					} else { // fail
						s.Fatalf("could not open device with VID=0x%X, PID=0x%X: %v",
							test.vid, test.pid, err)
					}
				} else {
					if nil == err { // fail
						ft.Close()
						s.Fatalf("opened device with VID=0x%X, PID=0x%X",
							test.vid, test.pid)
					} else { // success
					}
				}
			})
	}
}

func TestNewFT232HWithSerial(t *testing.T) {

	for _, test := range []struct {
		name   string
		serial string
		open   bool // true if expecting device to open successfully
	}{
		{name: "selected", serial: optSerial, open: true},
		{name: "empty", serial: "", open: true},
		{name: "random", serial: randString(16), open: false},
	} {
		t.Run(fmt.Sprintf("%s=\"%s\"", test.name, test.serial),
			func(s *testing.T) {
				ft, err := NewFT232HWithSerial(test.serial)
				if test.open {
					if nil == err { // success
						if err := ft.Close(); nil != err {
							s.Fatalf("could not close device with serial=\"%s\": %v",
								test.serial, err)
						}
					} else { // fail
						s.Fatalf("could not open device with serial=\"%s\": %v",
							test.serial, err)
					}
				} else {
					if nil == err { // fail
						ft.Close()
						s.Fatalf("opened device with serial=\"%s\"", test.serial)
					} else { // success
						// empty
					}
				}
			})
	}
}

func TestNewFT232HWithDesc(t *testing.T) {

	for _, test := range []struct {
		name string
		desc string
		open bool // true if expecting device to open successfully
	}{
		{name: "selected", desc: optDesc, open: true},
		{name: "empty", desc: "", open: true},
		{name: "random", desc: randString(64), open: false},
	} {
		t.Run(fmt.Sprintf("%s=\"%s\"", test.name, test.desc),
			func(s *testing.T) {
				ft, err := NewFT232HWithDesc(test.desc)
				if test.open {
					if nil == err { // success
						if err := ft.Close(); nil != err {
							s.Fatalf("could not close device with description=\"%s\": %v",
								test.desc, err)
						}
					} else { // fail
						s.Fatalf("could not open device with description=\"%s\": %v",
							test.desc, err)
					}
				} else {
					if nil == err { // fail
						ft.Close()
						s.Fatalf("opened device with description=\"%s\"", test.desc)
					} else { // success
						// empty
					}
				}
			})
	}
}

func TestParseUint32(t *testing.T) {

	for _, test := range []struct {
		name string
		str  string
		exp  uint32
		ok   bool // true if expecting parse to succeed
	}{
		{name: "", str: "0b10101010101010101010101010101010", exp: 0b10101010101010101010101010101010, ok: true},
		{name: "", str: "025252525252", exp: 025252525252, ok: true},
		{name: "", str: "0o25252525252", exp: 0o25252525252, ok: true},
		{name: "", str: "2863311530", exp: 2863311530, ok: true},
		{name: "", str: "0xAaAaAaAa", exp: 0xAaAaAaAa, ok: true},
		{name: "", str: "AaAaAaAa", exp: 0xAaAaAaAa, ok: true},
		{name: "", str: "0b01010101010101010101010101010101", exp: 0b01010101010101010101010101010101, ok: true},
		{name: "", str: "012525252525", exp: 012525252525, ok: true},
		{name: "", str: "0o12525252525", exp: 0o12525252525, ok: true},
		{name: "", str: "1431655765", exp: 1431655765, ok: true},
		{name: "", str: "0X55555555", exp: 0x55555555, ok: true},
		{name: "", str: "55555555", exp: 55555555, ok: true},
		{name: "", str: "-0b10101010", exp: 0, ok: false},
		{name: "", str: "-02525", exp: 0, ok: false},
		{name: "", str: "-0o2525", exp: 0, ok: false},
		{name: "", str: "-170", exp: 0, ok: false},
		{name: "", str: "-0xAa", exp: 0, ok: false},
		{name: "", str: "-Aa", exp: 0, ok: false},
		{name: "", str: "4294967296", exp: 0, ok: false},
	} {
		t.Run(fmt.Sprintf("%s=\"%s\"", test.name, test.str),
			func(s *testing.T) {
				act, ok := parseUint32(test.str)
				if test.ok {
					if ok {
						if act == test.exp { // success
							// empty
						} else { // fail
							s.Fatalf("parsed uint32={%d} expected={%d} from string=\"%s\"",
								act, test.exp, test.str)
						}
					} else { // fail
						s.Fatalf("could not parse uint32 from string=\"%s\"", test.str)
					}
				} else {
					if ok { // fail
						s.Fatalf("parsed uint32={%d} from string=\"%s\"", act, test.str)
					} else { // success
						// empty
					}
				}
			})
	}
}

func TestPin(t *testing.T) {

	for i := -1; i <= 8; i++ {

		ok := i > -1 && i < 8

		d := D(i)
		c := C(i)

		if ok {

			if !d.Valid() {
				t.Fatalf("expected D(%d) to be valid", i)
			}
			if !c.Valid() {
				t.Fatalf("expected C(%d) to be valid", i)
			}

			if (1<<i) != d.Mask() {
				t.Fatalf("D(%d) mask={%08b}, expected={%08b}", i, d.Mask(), 1<<i)
			}
			if (1<<i) != c.Mask() {
				t.Fatalf("C(%d) mask={%08b}, expected={%08b}", i, c.Mask(), 1<<i)
			}

			if i != d.Pos() {
				t.Fatalf("D(%d) pos={%d}, expected={%d}", i, d.Pos(), i)
			}
			if i != c.Pos() {
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
