package ft232h

// I2C holds an active I2C channel through which all communication is performed.
type I2C struct {
	device *FT232H
	config *i2cConfig
}

// i2cConfig holds all of the configuration settings for an I2C channel
type i2cConfig struct {
	clockRate I2CClockRate
	latency   uint8
	options   i2cOption
}

func i2cConfigDefault() *i2cConfig {
	return &i2cConfig{
		clockRate: I2CClockDefault,
		latency:   I2CLatencyDefault,
		options:   i2cDriveOnlyZeroDefault | i2c3PhaseClockingDefault,
	}
}

type i2cOption uint32

const (
	// 3-phase clocking is enabled by default. Setting this bit in ConfigOptions
	// will disable it
	i2c3PhaseClockingEnable  i2cOption = 0x00000000
	i2c3PhaseClockingDisable i2cOption = 0x00000001
	i2c3PhaseClockingDefault i2cOption = i2c3PhaseClockingEnable

	// The I2C master should actually drive the SDA line only when the output is
	// LOW. It should tristate the SDA line when the output should be HIGH.
	i2cDriveOnlyZeroDisable i2cOption = 0x00000000
	i2cDriveOnlyZeroEnable  i2cOption = 0x00000002
	i2cDriveOnlyZeroDefault i2cOption = i2cDriveOnlyZeroEnable
)

// I2CClockRate holds one of the supported I2C clock rate constants
type I2CClockRate uint32

// Constants defining the supported I2C clock rates
const (
	I2CClockStandardMode  I2CClockRate = 100000  // 100 kb/sec
	I2CClockFastMode      I2CClockRate = 400000  // 400 kb/sec
	I2CClockFastModePlus  I2CClockRate = 1000000 // 1000 kb/sec
	I2CClockHighSpeedMode I2CClockRate = 3400000 // 3.4 Mb/sec
	I2CClockDefault       I2CClockRate = I2CClockStandardMode
)

const (
	I2CLatencyDefault = 16
)

type i2cXferOption uint32

// Constants controlling the various I2C communication options
const (
	// Generate start condition before transmitting
	i2cStartBit i2cXferOption = 0x00000001

	// Generate stop condition before transmitting
	i2cStopBit i2cXferOption = 0x00000002

	// Continue transmitting data in bulk without caring about Ack or nAck from
	// device if this bit is not set. If this bit is set then stop transitting the
	// data in the buffer when the device nAcks
	i2cBreakOnNACK i2cXferOption = 0x00000004

	// libMPSSE-I2C generates an ACKs for every byte read. Some I2C slaves require
	// the I2C master to generate a nACK for the last data byte read. Setting this
	// bit enables working with such I2C slaves
	i2cNACKLastByte i2cXferOption = 0x00000008

	// no address phase, no USB interframe delays
	i2cFastTransferBytes i2cXferOption = 0x00000010
	i2cFastTransferBits  i2cXferOption = 0x00000020
	i2cFastTransfer      i2cXferOption = 0x00000030

	// if i2cFastTransfer is set then setting this bit would mean that the address
	// field should be ignored. The address is either a part of the data or this
	// is a special I2C frame that doesn't require an address
	i2cNoAddress i2cXferOption = 0x00000040

	// TBD
	// i2cCmdGetdeviceidRD = 0xF9
	// i2cCmdGetdeviceidWR = 0xF8
	// i2cGiveACK  = 1
	// i2cGiveNACK = 0
)
