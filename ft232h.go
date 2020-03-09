package ft232h

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// FT232H is the primary type for interacting with the device, holding the USB
// device file descriptor configuration/status and individual communication
// interfaces.
// Open a connection with an FT232H by calling the NewFT232H() constructor.
// If more than one FTDI device (any FTDI device, not just FT232H) is present on
// the system, there are several constructor variations of form NewFT232HWith*()
// to help distinguish which device to open. The default constructor NewFT232H()
// will attempt to parse command line flags to select a specific device.
// The only interface that is initialized by default is GPIO. You must call an
// initialization method of one of the other interfaces before using it.
type FT232H struct {
	info *deviceInfo
	mode Mode
	open *OpenFlag
	I2C  *I2C
	SPI  *SPI
	GPIO *GPIO
}

// String constructs a string representation of an FT232H device.
func (m *FT232H) String() string {
	return fmt.Sprintf("{ Index: %s, Mode: %s, Open: %+v, I2C: %+v, SPI: %+v, GPIO: %s }",
		m.info, m.mode, m.open, m.I2C, m.SPI, m.GPIO)
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

// OpenFlag contains the attributes used to distinguish which FT232H device to
// open from a command-line-style string slice.
type OpenFlag struct {
	flag   *flag.FlagSet
	index  *int
	vid    *int
	pid    *int
	serial *string
	desc   *string
}

// String returns a descriptive string of all flags successfully parsed.
// func (o *OpenFlag) String() string {
// 	if o.flag.NFlag() > 0 {
// 		s := []string{}
// 		o.flag.Visit(func(f *flag.Flag) {
// 			s = append(s, fmt.Sprintf("-%s=%q", f.Name, f.Value))
// 		})
// 		return fmt.Sprintf("{ %s }", strings.Join(s, " "))
// 	} else {
// 		return "(none)"
// 	}
// }

// NewFT232H attempts to open a connection with the first MPSSE-capable USB
// device matching flags given at the command line. Use -h to see all of the
// supported flags.
func NewFT232H() (*FT232H, error) {
	return NewFT232HWithFlag(os.Args[1:], true)
}

// NewFT232HWithIndex attempts to open a connection with the MPSSE-capable USB
// device enumerated at index (starting at 0). Returns non-nil error if
// unsuccessful. A negative index is equivalent to 0.
func NewFT232HWithIndex(index int) (*FT232H, error) {
	if index < 0 {
		index = 0
	}
	return NewFT232HWithMask(&OpenMask{Index: fmt.Sprintf("%d", index)})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given vendor ID vid and product ID pid. Returns a non-nil
// error if unsuccessful.
func NewFT232HWithVIDPID(vid uint16, pid uint16) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{
		VID: fmt.Sprintf("%d", vid),
		PID: fmt.Sprintf("%d", pid),
	})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given serial no. Returns a non-nil error if unsuccessful.
// An empty string matches any serial number.
func NewFT232HWithSerial(serial string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Serial: serial})
}

// NewFT232HWithIndex attempts to open a connection with the first MPSSE-capable
// USB device with given description. Returns a non-nil error if unsuccessful.
// An empty string matches any description.
func NewFT232HWithDesc(desc string) (*FT232H, error) {
	return NewFT232HWithMask(&OpenMask{Desc: desc})
}

// NewFT232HWithFlag attempts to open a connection with the first MPSSE-capable
// USB device matching flags given in a command-line-style string slice.
// See type OpenFlag and func NewOpenFlag() for details.
func NewFT232HWithFlag(arg []string, fatal bool) (*FT232H, error) {
	o := NewOpenFlag(fatal)
	if len(arg) > 0 {
		if err := o.Parse(arg); nil != err {
		}
	}
	ft, err := NewFT232HWithMask(o.OpenMask())
	if nil != err {
		return nil, err
	}
	ft.open = o // keep a copy of the flagset for use/inspection by test suite
	return ft, nil
}

// NewFT232HWithMask attempts to open a connection with the first MPSSE-capable
// USB device matching all of the given attributes. Returns a non-nil error if
// unsuccessful. Uses the first device found if mask is nil or all attributes
// are empty strings.
//
// The attributes are each specified as strings, including the integers, so that
// any attribute not given (i.e. empty string) will never exclude a device. The
// integer attributes can be expressed in any base recognized by the Go grammar
// for numeric literals (e.g., "13", "0b1101", "0xD", and "D" are all valid and
// equivalent).
func NewFT232HWithMask(mask *OpenMask) (*FT232H, error) {
	m := &FT232H{info: nil, mode: ModeNone, open: nil, I2C: nil, SPI: nil, GPIO: nil}
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

// NewOpenFlag constructs a new FlagSet with fields to describe an FT232H.
// If fatal is true, the program will call os.Exit() if the flag parser fails on
// malformed input, unrecognized flags are provided, or the default help flag -h
// is received.
func NewOpenFlag(fatal bool) *OpenFlag {
	const (
		indexDefault  int    = 0
		vidDefault    int    = 0x0403
		pidDefault    int    = 0x6014
		serialDefault string = ""
		descDefault   string = ""
	)
	onError := flag.ContinueOnError
	if fatal {
		onError = flag.ExitOnError
	}
	f := flag.NewFlagSet(os.Args[0]+" open flags", onError)
	o := &OpenFlag{
		flag:   f,
		index:  f.Int("index", indexDefault, "open device enumerated at index `N` ≥ 0"),
		vid:    f.Int("vid", vidDefault, "open device with vendor ID"),
		pid:    f.Int("pid", pidDefault, "open device with product ID"),
		serial: f.String("serial", serialDefault, "open device with identifier"),
		desc:   f.String("desc", descDefault, "open device with description"),
	}
	return o
}

// BlessOpenFlag registers the flags in the flag package's default, top-level
// FlagSet var flag.CommandLine. This lets external packages (e.g. `go test`)
// inherit these flags and not call os.Exit() when these otherwise unexpected
// flags are received.
func BlessOpenFlag() {
	o := NewOpenFlag(false)
	o.flag.VisitAll(func(f *flag.Flag) {
		if nil == flag.Lookup(f.Name) {
			flag.Var(f.Value, f.Name, f.Usage)
		}
	})
}

// Parse parses flags from the given slice of strings arg into the fields of its
// receiver, and silently ignores any unexpected flags.
func (o *OpenFlag) Parse(arg []string) error {

	if o.flag.ErrorHandling() == flag.ContinueOnError {
		o.flag.SetOutput(ioutil.Discard)
	}

	// extract only the known flags from arg, so that we don't die when the user
	// provides unknown flags handled by other packages (e.g. `go test`)
	parse, keep := []string{}, false
	for _, a := range arg {

		// we set keep=true when the element being processed is the argument to a
		// known flag that was previously processed.
		// always add it to the slice to be parsed.
		if keep {
			parse, keep = append(parse, a), false
			continue // start processing next element
		}

		// ignore non-flag elements
		if !strings.HasPrefix(a, "-") {
			continue
		}

		// split any flags given as a single argument, i.e. "-flag=value" format,
		// on the first "=" found, subsequent "=" are preserved in s[1]
		s := strings.SplitN(a, "=", 2)

		// check if it is a known OpenFlag (but remove the flag prefix "-" first)
		if f := o.flag.Lookup(strings.TrimPrefix(s[0], "-")); nil != f {

			// this is a recognized OpenFlag. copy it to the slice to be parsed.
			parse = append(parse, a)

			// we need to keep the next element in arg if this is a non-bool flag and
			// its value was not already provided (i.e. using form "-flag=value").
			switch f.Value.(type) {
			case interface{ IsBoolFlag() }:
				// bool flags cannot be expressed with form "-flag" "value"
			default:
				keep = 1 == len(s)
			}
		}
	}

	if err := o.flag.Parse(parse); nil != err {
		return err
	}
	return nil
}

// OpenMask constructs an OpenMask using the parsed flags explicitly provided.
// If the OpenFlag has not yet been parsed, a zero OpenMask is returned that
// matches all devices.
func (o *OpenFlag) OpenMask() *OpenMask {
	m := &OpenMask{}
	o.flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "index":
			m.Index = f.Value.String()
		case "vid":
			m.VID = f.Value.String()
		case "pid":
			m.PID = f.Value.String()
		case "serial":
			m.Serial = f.Value.String()
		case "desc":
			m.Desc = f.Value.String()
		}
	})
	return m
}

// openDevice attempts to open the device matching all fields of a given mask.
// Returns a non-nil error if unsuccessful. The error SDeviceNotFound is
// returned if no device was found matching the given mask.
// See NewFT232HWithMask for semantics.
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

// Close closes the USB connection with an FT232H. Returns a non-nil error if
// unsuccessful.
func (m *FT232H) Close() error {
	if nil != m.info {
		return m.info.close()
	}
	m.mode = ModeNone
	return nil
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
	return fmt.Sprintf("%d:{ Open = %t, HiSpeed = %t, Chip = %q (0x%02X), "+
		"VID = 0x%04X, PID = 0x%04X, Location = %04X, "+
		"Serial = %q, Desc = %q, Handle = %p }",
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
// Returns a non-nil error if unsuccessful.
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

// devices queries all of the USB devices on the system using the D2XX bridge
// and returns a slice of deviceInfo pointers for all MPSSE-capable devices.
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
