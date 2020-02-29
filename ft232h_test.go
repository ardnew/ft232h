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
