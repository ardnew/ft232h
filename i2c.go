package ft232h

import (
	"fmt"
	"math/bits"
)

// I2C stores interface configuration settings for an I²C master and provides
// methods for reading and writing to I²C slave devices.
// The interface must be initialized by calling either Init or Config (not both)
// before use.
type I2C struct {
	device *FT232H
	config *i2cConfig
}

// String returns a descriptive string of an I²C interface.
func (i2c *I2C) String() string {
	return fmt.Sprintf("{ FT232H: %p, Config: %s }", i2c.device, i2c.config)
}

// I2CConfig holds all of the configuration settings for initializing an I²C
// interface.
type I2CConfig struct {
	*I2COption
	Clock        I2CClockRate // 100000 (100 kb/s) - 3400000 (3.4 Mb/s)
	Latency      byte         // 1-255 USB HiSpeed, 2-255 USB FullSpeed
	Clock3Phase  bool         // I²C 3-phase clocking enabled=true/disabled=false
	LowDriveOnly bool         // float HIGH (pullup) if true, drive HIGH if false
}

// I2CConfigDefault returns the default configuration settings for an I²C
// interface.
func I2CConfigDefault() *I2CConfig {
	return i2cConfigDefault().I2CConfig()
}

// I2CConfig returns the current configuration settings of the I2C receiver.
func (i2c *I2C) GetConfig() *I2CConfig {
	return i2c.config.I2CConfig()
}

// Constants defining legal 7-bit I²C slave addresses
const (
	I2CSlaveAddressMin = 0x08
	I2CSlaveAddressMax = 0x77
)

// Constants related to I²C interface initialization.
const (
	I2CClockMaximum   I2CClockRate = I2CClockHighSpeedMode
	I2CClockDefault   I2CClockRate = I2CClockFastMode
	I2CLatencyDefault byte         = 2
)

// i2cConfig holds all of the configuration settings for an I²C channel stored
// privately in each instance of I2C.
type i2cConfig struct {
	clockRate I2CClockRate
	latency   uint8
	options   i2cOption
	breakNACK bool
	readNACK  bool
	noDelay   bool
}

// String returns a descriptive string of an i2cConfig.
func (c i2cConfig) String() string {
	return fmt.Sprintf("{ Clock: %q, Latency: \"%d ms\", Options: %s, "+
		"BreakOnNACK: %t, NACKAfterRead: %t, NoUSBDelay: %t }",
		c.clockRate, c.latency, c.options, c.breakNACK, c.readNACK, c.noDelay)
}

// i2cConfigDefault returns an i2cConfig struct stored in the private
// configuration field of an I2C instance with the default settings for all
// fields.
func i2cConfigDefault() *i2cConfig {
	return &i2cConfig{
		clockRate: I2CClockDefault,
		latency:   I2CLatencyDefault,
		options:   i2cOptionDefault,
		breakNACK: i2cBreakNACKDefault,
		readNACK:  i2cLastNACKDefault,
		noDelay:   i2cNoDelayDefault,
	}
}

// I2CConfig constructs an I2C configuration struct using the settings stored in
// the private configuration field of an instance of I2C.
func (c *i2cConfig) I2CConfig() *I2CConfig {
	return &I2CConfig{
		I2COption: &I2COption{
			BreakOnNACK:  c.breakNACK,
			LastReadNACK: c.readNACK,
			NoUSBDelay:   c.noDelay,
		},
		Clock:        c.clockRate,
		Latency:      c.latency,
		Clock3Phase:  c.options.clock3Phase(),
		LowDriveOnly: c.options.lowDriveOnly(),
	}
}

// I2CClockRate holds one of the supported I²C clock rate constants.
type I2CClockRate uint32

// Constants defining the supported I²C clock rates.%d ms
const (
	I2CClockStandardMode  I2CClockRate = 100000  // 100 kb/sec
	I2CClockFastMode      I2CClockRate = 400000  // 400 kb/sec
	I2CClockFastModePlus  I2CClockRate = 1000000 // 1000 kb/sec
	I2CClockHighSpeedMode I2CClockRate = 3400000 // 3.4 Mb/sec
)

// String returns a descriptive string of an I2CClockRate.
func (c I2CClockRate) String() string {
	switch c {
	case I2CClockStandardMode:
		return "Standard mode (100 KHz)"
	case I2CClockFastMode:
		return "Fast mode (400 KHz)"
	case I2CClockFastModePlus:
		return "Fast mode plus (1000 KHz)"
	case I2CClockHighSpeedMode:
		return "High-speed mode (3.4 MHz)"
	default:
		return fmt.Sprintf("Unsupported mode (%d Hz)", c)
	}

}

// I2COption holds all of the dynamic configuration settings that can be changed
// while an I²C interface is open.
type I2COption struct {
	BreakOnNACK  bool // do not continue reading/writing stream on slave NACK
	LastReadNACK bool // send NACK after last byte read from I²C slave
	NoUSBDelay   bool // pack all I²C data into the fewest number of USB packets
}

// i2cOption stores the various I²C configuration options as a 32-bit bitmap.
type i2cOption uint32

// String returns a descriptive string of an i2cOption.
func (o i2cOption) String() string {
	return fmt.Sprintf("{ 3-PhaseClock: %t, LowDriveOnly: %t }",
		(o&i2cClock3PhaseMask) == i2cClock3PhaseEnable,
		(o&i2cLowDriveOnlyMask) == i2cLowDriveOnlyEnable,
	)
}

// Constants defining the options for I²C 3-phase clocking (enabled by default).
const (
	i2cClock3PhaseEnable  i2cOption = 0x00000000
	i2cClock3PhaseDisable i2cOption = 0x00000001
	i2cClock3PhaseMask    i2cOption = 0x00000001
	i2cClock3PhaseDefault i2cOption = i2cClock3PhaseEnable
)

// Constants defining how the I²C SDA pin is driven. If LOW drive only is
// enabled (default), the SDA line is driven only when outputting LOW, and is
// otherwise floating (hence you will need a pullup). If disabled, the SDA line
// is driven LOW and HIGH (via internal 75K pullup, I believe).
const (
	i2cLowDriveOnlyDisable i2cOption = 0x00000000
	i2cLowDriveOnlyEnable  i2cOption = 0x00000002
	i2cLowDriveOnlyMask    i2cOption = 0x00000002
	i2cLowDriveOnlyDefault i2cOption = i2cLowDriveOnlyEnable
)

// Constants with values shared by fields of the I²C configuration.
const (
	i2cOptionInvalid i2cOption = 0xAAAAAAAA
	i2cOptionDefault           = i2cLowDriveOnlyDefault | i2cClock3PhaseDefault
)

// Constants related to the dynamic configuration options of an I²C interface.
const (
	i2cBreakNACKDefault = false
	i2cLastNACKDefault  = false
	i2cNoDelayDefault   = true
)

// Valid verifies the i2cOption receiver opt isnt equal to the sentinel value
// for invalid I²C options.
func (opt i2cOption) Valid() bool { return opt != i2cOptionInvalid }

// clock3Phase returns the 3-phase clocking setting (true=enabled) of the
// receiver i2cOption.
func (opt i2cOption) clock3Phase() bool {
	return (opt & i2cClock3PhaseMask) == i2cClock3PhaseEnable
}

// clock3Phase returns the LOW drive only setting (true=enabled) of the
// receiver i2cOption.
func (opt i2cOption) lowDriveOnly() bool {
	return (opt & i2cLowDriveOnlyMask) == i2cLowDriveOnlyEnable
}

// i2cXferOption stores the various I²C transfer options as a 32-bit bitmap.
type i2cXferOption uint32

// Constants controlling the various I²C communication options
const (
	// Generate start condition before transmitting
	i2cStartBit i2cXferOption = 0x00000001

	// Generate stop condition before transmitting
	i2cStopBit i2cXferOption = 0x00000002

	// Continue transmitting data in bulk without caring about ACK or NACK from
	// device if this bit is not set. If this bit is set then stop transitting the
	// data in the buffer when the device NACKs
	i2cBreakOnNACK i2cXferOption = 0x00000004

	// libMPSSE generates an ACKs for every byte read. Some I²C slaves require the
	// I²C master to generate a NACK for the last data byte read. Setting this bit
	// enables working with such I²C slaves
	i2cLastReadNACK i2cXferOption = 0x00000008

	// no address phase, no USB interframe delays
	i2cFastTransferBytes i2cXferOption = 0x00000010
	i2cFastTransferBits  i2cXferOption = 0x00000020
	i2cFastTransfer      i2cXferOption = 0x00000030

	// if i2cFastTransfer is set then setting this bit would mean that the address
	// field should be ignored. The address is either a part of the data or this
	// is a special I²C frame that doesn't require an address
	i2cNoAddress i2cXferOption = 0x00000040

	// default read/write options
	i2cXferDefault i2cXferOption = 0x00000000

	// TBD
	// i2cCmdGetdeviceidRD = 0xF9
	// i2cCmdGetdeviceidWR = 0xF8
	// i2cGiveACK  = 1
	// i2cGiveNACK = 0
)

// Option changes the dynamic configuration parameters of the I²C interface.
// It can be called while the I²C interface is open without having to first
// close and reopen the device.
func (i2c *I2C) Option(opt *I2COption) error {

	i2c.config.breakNACK = opt.BreakOnNACK
	i2c.config.readNACK = opt.LastReadNACK
	i2c.config.noDelay = opt.NoUSBDelay

	return nil
}

// Config initializes the I²C interface with the given configuration to a state
// ready for read/write.
// If the given configuration is nil, the default configuration is used (see
// I2CConfigDefault).
// It is not necessary to call Init after calling Config.
// See documentation of Init for other semantics.
func (i2c *I2C) Config(cfg *I2CConfig) error {

	if nil == cfg {
		cfg = I2CConfigDefault()
	}

	if 0 == cfg.Clock {
		i2c.config.clockRate = I2CClockDefault
	} else {
		if cfg.Clock <= I2CClockMaximum {
			i2c.config.clockRate = cfg.Clock
		} else {
			return fmt.Errorf("invalid clock rate: %d", cfg.Clock)
		}
	}

	if 0 == cfg.Latency {
		i2c.config.latency = I2CLatencyDefault
	} else {
		i2c.config.latency = cfg.Latency
	}

	phaseOpt := i2cClock3PhaseDisable
	if cfg.Clock3Phase {
		phaseOpt = i2cClock3PhaseEnable
	}

	driveOpt := i2cLowDriveOnlyDisable
	if cfg.LowDriveOnly {
		driveOpt = i2cLowDriveOnlyEnable
	}

	i2c.config.options = driveOpt | phaseOpt

	if err := i2c.Option(cfg.I2COption); nil != err {
		return err
	}

	return i2c.Init()
}

// Init initializes the I²C interface to a state ready for read/write.
// If Config has not been called, the default configuration is used (see
// I2CConfigDefault).
// If the interface is already initialized, it is first closed before
// initializing the interface.
func (i2c *I2C) Init() error {

	if err := _I2C_InitChannel(i2c); nil != err {
		return err
	}

	i2c.device.mode = ModeI2C

	return i2c.device.GPIO.Init() // reset GPIO
}

// Close closes both the I²C interface and the connection to the FT232H device.
func (i2c *I2C) Close() error {
	return i2c.device.Close()
}

// Read reads the given count number of bytes from the I²C interface.
// The given slave is the unshifted 7-bit I²C slave address to read from.
// There is no maximum length for the number of bytes to read.
// If start is true, an I²C start condition is generated before transfer.
// If stop is true, an I²C stop condition is generated after transfer.
// Returns the slice of bytes successfully read and a non-nil error if there was
// an error.
func (i2c *I2C) Read(slave uint, count uint, start bool, stop bool) ([]uint8, error) {

	if !(slave >= I2CSlaveAddressMin && slave <= I2CSlaveAddressMax) {
		return nil, fmt.Errorf("invalid slave address (0x%02X-0x%02X): 0x%02X",
			I2CSlaveAddressMin, I2CSlaveAddressMax, slave)
	}

	opt := i2cXferDefault

	if start {
		opt |= i2cStartBit
	}

	if stop {
		opt |= i2cStopBit
	}

	if i2c.config.noDelay {
		opt |= i2cFastTransfer | i2cFastTransferBytes
	}

	// these flags are not supported when fast transfer (I2COption.NoUSBDelay)
	// is enabled with start/stop condition generation
	if !(i2c.config.noDelay && (start || stop)) {
		if i2c.config.readNACK {
			opt |= i2cLastReadNACK
		}
		if i2c.config.breakNACK {
			opt |= i2cBreakOnNACK
		}
	}

	return _I2C_Read(i2c, slave, count, opt)
}

// Write writes the given byte slice data to the I²C interface.
// The given slave is the unshifted 7-bit I²C slave address to write to.
// There is no maximum length for the data slice.
// If start is true, an I²C start condition is generated before transfer.
// If stop is true, an I²C stop condition is generated after transfer.
// Returns the slice of bytes successfully written and a non-nil error if there
// was an error.
func (i2c *I2C) Write(slave uint, data []uint8, start bool, stop bool) (uint, error) {

	if !(slave >= I2CSlaveAddressMin && slave <= I2CSlaveAddressMax) {
		return 0, fmt.Errorf("invalid slave address (0x%02X-0x%02X): 0x%02X",
			I2CSlaveAddressMin, I2CSlaveAddressMax, slave)
	}

	opt := i2cXferDefault

	if start {
		opt |= i2cStartBit
	}

	if stop {
		opt |= i2cStopBit
	}

	if i2c.config.noDelay {
		opt |= i2cFastTransfer | i2cFastTransferBytes
	}

	// this flag is not supported when fast transfer (I2COption.NoUSBDelay) is
	// enabled with start/stop condition generation
	if !(i2c.config.noDelay && (start || stop)) {
		if i2c.config.breakNACK {
			opt |= i2cBreakOnNACK
		}
	}

	return _I2C_Write(i2c, slave, data, opt)
}

// I2CReg represents a read-write register of an I²C slave device.
type I2CReg struct {
	i2c   *I2C      // the I²C interface to use
	slave uint      // unshifted 7-bit I²C slave address
	addr  uint      // register sub-address to read/write
	space AddrSpace // sub-address space used to format register in data payload
	order ByteOrder // byte order used to format register+data in data payload
}

// I2CRegReader represents a method for reading I²C slave device registers.
type I2CRegReader func(rewrite bool) (uint64, error)

// validate checks that fields of an I2CReg pass basic sanity requirements.
// Returns the byte-ordered register sub-address to send in read/write payloads.
func (reg *I2CReg) validate() ([]uint8, error) {

	if nil == reg {
		return nil, fmt.Errorf("invalid receiver (nil)")
	}

	if !(reg.slave >= I2CSlaveAddressMin && reg.slave <= I2CSlaveAddressMax) {
		return nil, fmt.Errorf("invalid I²C slave address (0x%02X-0x%02X): 0x%02X",
			I2CSlaveAddressMin, I2CSlaveAddressMax, reg.slave)
	}

	b := reg.space.Bytes()
	if 0 == b {
		return nil, fmt.Errorf("invalid sub-address space: %d", reg.space)
	}

	if bits.Len(reg.addr) > bits.Len(1<<(reg.space.Bits()-1)) {
		return nil, fmt.Errorf("register sub-address outside %s address space: 0x%02X",
			reg.space, reg.addr)
	}

	r := reg.order.Bytes(b, uint64(reg.addr))
	if 0 == len(r) {
		return nil, fmt.Errorf("invalid byte order: %d", reg.order)
	}

	return r, nil
}

// Reg constructs a new I2CReg for conveniently reading and writing data in I²C
// slave device registers.
func (i2c *I2C) Reg(slave uint, addr uint, space AddrSpace, order ByteOrder) *I2CReg {
	return &I2CReg{
		i2c:   i2c,
		slave: slave,
		addr:  addr,
		space: space,
		order: order,
	}
}

// Reader returns a closure that can be used to repeatedly read from a register.
// The size argument defines the number of bytes to read, i.e. the size of the
// data read from the register, often 2 (16-bit); it is also used along with the
// I2CReg byte order to format the value returned by the closure.
//
// The returned closure accepts a single argument rewrite that if true will
// reposition the register pointer to the receiver's register address, which is
// necessary when reading/writing multiple registers; otherwise, false, reads
// will occur much faster without having to reposition every time.
// The byte order given to the Reg() constructor is used to format the value
// returned when calling the closure.
func (reg *I2CReg) Reader(size uint) (I2CRegReader, error) {

	addr, err := reg.validate()
	if nil != err {
		return nil, err
	}

	if _, err := reg.i2c.Write(reg.slave, addr, true, false); nil != err {
		return nil, err
	}

	return func(rewrite bool) (uint64, error) {

		if rewrite {
			if _, err := reg.i2c.Write(reg.slave, addr, true, false); nil != err {
				return 0, err
			}
		}

		if dat, err := reg.i2c.Read(reg.slave, size, true, true); nil != err {
			return 0, err
		} else {
			return reg.order.Uint(size, dat), nil
		}

	}, nil
}
