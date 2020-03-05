package ft232h

import (
	"fmt"
	"log"
)

// I2C stores interface configuration settings for an I²C master and provides
// methods for reading and writing to I²C slave devices.
// The interface must be initialized by calling either Init or Config (not both)
// before use.
type I2C struct {
	device *FT232H
	config *i2cConfig
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
	lastNACK  bool
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
		lastNACK:  i2cLastNACKDefault,
	}
}

// I2CConfig constructs an I2C configuration struct using the settings stored in
// the private configuration field of an instance of I2C.
func (c *i2cConfig) I2CConfig() *I2CConfig {
	return &I2CConfig{
		I2COption: &I2COption{
			BreakOnNACK:  c.breakNACK,
			LastByteNACK: c.lastNACK,
		},
		Clock:        c.clockRate,
		Latency:      c.latency,
		Clock3Phase:  c.options.clock3Phase(),
		LowDriveOnly: c.options.lowDriveOnly(),
	}
}

// I2CClockRate holds one of the supported I²C clock rate constants.
type I2CClockRate uint32

// Constants defining the supported I²C clock rates.
const (
	I2CClockStandardMode  I2CClockRate = 100000  // 100 kb/sec
	I2CClockFastMode      I2CClockRate = 400000  // 400 kb/sec
	I2CClockFastModePlus  I2CClockRate = 1000000 // 1000 kb/sec
	I2CClockHighSpeedMode I2CClockRate = 3400000 // 3.4 Mb/sec
)

// I2COption holds all of the dynamic configuration settings that can be changed
// while an I²C interface is open.
type I2COption struct {
	BreakOnNACK  bool // do not continue reading/writing stream on slave NACK
	LastByteNACK bool // send NACK after last byte read from I²C slave
}

// i2cOption stores the various I²C configuration options as a 32-bit bitmap.
type i2cOption uint32

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
	i2cNACKLastByte i2cXferOption = 0x00000008

	// no address phase, no USB interframe delays
	i2cFastTransferBytes i2cXferOption = 0x00000010
	i2cFastTransferBits  i2cXferOption = 0x00000020
	i2cFastTransfer      i2cXferOption = 0x00000030

	// if i2cFastTransfer is set then setting this bit would mean that the address
	// field should be ignored. The address is either a part of the data or this
	// is a special I²C frame that doesn't require an address
	i2cNoAddress i2cXferOption = 0x00000040

	// default read/write options
	i2cXferDefault = i2cFastTransfer | i2cFastTransferBytes

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
	i2c.config.lastNACK = opt.LastByteNACK

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
// The given address is the unshifted 7-bit I²C slave address to read from.
// There is no maximum length for the number of bytes to read.
// If start is true, an I²C start condition is generated before transfer.
// If stop is true, an I²C stop condition is generated after transfer.
// Returns the slice of bytes successfully read and a non-nil error if there was
// an error.
func (i2c *I2C) Read(addr uint, count uint, start bool, stop bool) ([]uint8, error) {

	opt := i2cXferDefault

	if start {
		opt |= i2cStartBit
	}

	if stop {
		opt |= i2cStopBit
	}

	// following flags are not compatible when fast transfer is enabled (deffault)
	// with start/stop condition generation
	if !start && !stop {
		if i2c.config.lastNACK {
			opt |= i2cNACKLastByte
		}
		if i2c.config.breakNACK {
			opt |= i2cBreakOnNACK
		}
	}

	if addr > 0x7F {
		return nil, fmt.Errorf("invalid slave address (0x00-0x7F): 0x%02X", addr)
	}

	log.Printf("<<<<<< [%02X, {%+v}, {%032b}]", addr, count, opt)
	return _I2C_Read(i2c, addr, count, opt)
}

// Write writes the given byte slice data to the I²C interface.
// The given address is the unshifted 7-bit I²C slave address to write to.
// There is no maximum length for the data slice.
// If start is true, an I²C start condition is generated before transfer.
// If stop is true, an I²C stop condition is generated after transfer.
// Returns the slice of bytes successfully written and a non-nil error if there
// was an error.
func (i2c *I2C) Write(addr uint, data []uint8, start bool, stop bool) (uint, error) {

	opt := i2cXferDefault

	if start {
		opt |= i2cStartBit
	}

	if stop {
		opt |= i2cStopBit
	}

	// following flag is not compatible when fast transfer is enabled (deffault)
	// with start/stop condition generation
	if !start && !stop {
		if i2c.config.breakNACK {
			opt |= i2cBreakOnNACK
		}
	}

	if addr > 0x7F {
		return 0, fmt.Errorf("invalid slave address (0x00-0x7F): 0x%02X", addr)
	}

	log.Printf(">>>>>> [%02X, {%+v}, {%032b}]", addr, data, opt)
	return _I2C_Write(i2c, addr, data, opt)
}
