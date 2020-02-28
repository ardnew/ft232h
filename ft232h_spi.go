package ft232h

import "fmt"

// SPI stores interface configuration settings for an SPI master and provides
// methods for reading and writing to SPI slave devices.
type SPI struct {
	device *FT232H
	config *spiConfig
}

// SPIConfig holds all of the configuration settings for initializing an SPI
// interface.
type SPIConfig struct {
	*SPIOption
	Clock   uint32 // valid range: 0-30000000 (30 MHz)
	Latency byte   // 1-255 USB HiSpeed, 2-255 USB FullSpeed
}

// NewSPIConfig constructs an SPIConfig struct with the given configuration.
func NewSPIConfig(cs Pin, activeLow bool, mode byte, clock uint32, latency byte) *SPIConfig {
	return &SPIConfig{
		SPIOption: NewSPIOption(cs, activeLow, mode),
		Clock:     clock,
		Latency:   latency,
	}
}

// SPIConfigDefault returns the default configuration settings for an SPI
// interface.
func SPIConfigDefault() *SPIConfig {
	return spiConfigDefault().GetConfig()
}

// SPIConfig returns the current configuration settings of the SPI receiver.
func (spi *SPI) SPIConfig() *SPIConfig {
	return spi.config.GetConfig()
}

// Constants related to SPI interface initialization.
const (
	SPIClockMaximum   uint32 = 30000000
	SPIClockDefault   uint32 = 12000000
	SPILatencyDefault byte   = 16
)

// spiConfig holds all of the configuration settings for an SPI channel stored
// privately in each instance of SPI.
type spiConfig struct {
	clockRate  uint32 // in Hertz
	latency    uint8  // in ms
	options    spiOption
	pin        uint32 // port D pins ("low byte lines of MPSSE")
	chipSelect Pin    // may be DPin (MPSSE low byte) or CPin (GPIO)
}

// spiConfigDefault returns an spiConfig struct stored in the private
// configuration field of an SPI instance with the default settings for all
// fields.
func spiConfigDefault() *spiConfig {
	return &spiConfig{
		clockRate:  SPIClockDefault,
		latency:    SPILatencyDefault,
		options:    spiOptionDefault,
		pin:        spiPinConfigDefault(),
		chipSelect: spiCSDefault.cs(),
	}
}

// GetConfig constructs an SPI configuration struct using the settings stored in
// the private configuration field of an instance of SPI.
func (c *spiConfig) GetConfig() *SPIConfig {
	return &SPIConfig{
		SPIOption: NewSPIOption(c.options.cs(), c.options.activeLow(), c.options.mode()),
		Clock:     c.clockRate,
		Latency:   c.latency,
	}
}

// SPIOption holds all of the dynamic configuration settings that can be changed
// while an SPI interface is open.
type SPIOption struct {
	CS        Pin  // CS pin to assert when writing (can be DPin or CPin (GPIO))
	ActiveLow bool // CS asserted "active" by driving pin LOW or HIGH
	Mode      byte // SPI operating mode (mode 0 and 2 support only)
}

// NewSPIOption constructs an SPIOption struct with the given options.
func NewSPIOption(cs Pin, activeLow bool, mode byte) *SPIOption {
	return &SPIOption{
		CS:        cs,
		ActiveLow: activeLow,
		Mode:      mode,
	}
}

// spiOption stores the various SPI configuration options as a 32-bit bitmap.
type spiOption uint32

// Constants defining SPI operating modes (supports mode 0 and 2 only (CPHA=2))
const (
	spiMode0       spiOption = 0x00000000 // capture on RISE, propagate on FALL
	spiMode1       spiOption = 0x00000001 // capture on FALL, propagate on RISE
	spiMode2       spiOption = 0x00000002 // capture on FALL, propagate on RISE
	spiMode3       spiOption = 0x00000003 // capture on RISE, propagate on FALL
	spiModeMask    spiOption = 0x00000003
	spiModeDefault           = spiMode0
)

// Constants defining CS pins capable of using auto-assertion (DBUS pins only)
const (
	spiCSD3      spiOption = 0x00000000 // SPI CS on D3
	spiCSD4      spiOption = 0x00000004 // SPI CS on D4
	spiCSD5      spiOption = 0x00000008 // SPI CS on D5
	spiCSD6      spiOption = 0x0000000C // SPI CS on D6
	spiCSD7      spiOption = 0x00000010 // SPI CS on D7
	spiCSMask    spiOption = 0x0000001C
	spiCSDefault           = spiCSD3
)

// Constants defining the polarity of CS assertion
const (
	spiCSActiveLow     spiOption = 0x00000020 // drive pin low to assert CS
	spiCSActiveHigh    spiOption = 0x00000000 // drive pin high to assert CS
	spiCSActiveMask    spiOption = 0x00000020
	spiCSActiveDefault           = spiCSActiveLow
)

// Constants with values shared by fields of the SPI configuration.
const (
	spiOptionInvalid spiOption = 0xAAAAAAAA
	spiOptionDefault           = spiCSActiveDefault | spiCSDefault | spiModeDefault
)

// Valid verifies the spiOption receiver opt isnt equal to the sentinel value
// for invalid SPI options.
func (opt spiOption) Valid() bool { return opt != spiOptionInvalid }

// mode reads the SPI mode in the spiOption receiver opt and returns its value
// as a byte (0 = mode 0, ..., 3 = mode 3)
func (opt spiOption) mode() byte {
	return byte(opt & spiModeMask)
}

// cs reads the chip-select mask in the spiOption receiver opt and returns its
// corresponding DPin as type Pin.
func (opt spiOption) cs() Pin {
	switch opt & spiCSMask {
	case spiCSD3, spiCSD4, spiCSD5, spiCSD6, spiCSD7:
		return D(int(opt>>2) + 3)
	default:
		return DPin(0) // invalid pin
	}
}

// activeLow reads the active-low/high flag in the spiOption receiver opt and
// returns true if CS is asserted by driving pin LOW, false if pin HIGH.
func (opt spiOption) activeLow() bool {
	return spiCSActiveLow == (opt & spiCSActiveMask)
}

// spiOptionCS translates a DPin p to its corresponding chip-select mask for the
// option field of an SPI configuration struct.
func (p DPin) spiOptionCS() spiOption {
	if p.Valid() && p.Pos() >= 3 {
		return (spiOption(p.Pos()-3) << 2) & spiCSMask
	} else {
		return spiOptionInvalid
	}
}

// spiPinConfig represents the default direction and value for DBUS pins when
// MPSSE is configured for SPI operation.
type spiPinConfig struct {
	initDir  byte // direction of lines after SPI channel initialization
	initVal  byte // value of lines after SPI channel initialization
	closeDir byte // direction of lines after SPI channel is closed
	closeVal byte // value of lines after SPI channel is closed
}

// spiPin creates a bitmask from the DPin p and spiPinConfig cfg for the pin
// field of an SPI configuration struct.
func (p DPin) spiPin(cfg *spiPinConfig) uint32 {
	if p.Valid() {
		pos := p.Pos()
		return 0 | // <- for formatting
			(uint32(cfg.initDir) << (pos + 0)) |
			(uint32(cfg.initVal) << (pos + 8)) |
			(uint32(cfg.closeDir) << (pos + 16)) |
			(uint32(cfg.closeVal) << (pos + 24))
	} else {
		return 0
	}
}

// spiPinConfigDefault defines the default spiPinConfig value for all DBUS pins.
// all output pins are configured LOW except for the default CS pin (D3) since
// we also have spiCSActiveLow by default. this means we won't activate the
// default slave line until intended. it also means SCLK idles LOW (change
// initVal to PinHI to idle HIGH).
func spiPinConfigDefault() uint32 {
	var pin uint32
	for i, cfg := range [NumDPins]*spiPinConfig{
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D0 SCLK
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D1 MOSI
		&spiPinConfig{initDir: PinIN, initVal: PinLO, closeDir: PinIN, closeVal: PinLO}, // D2 MISO
		&spiPinConfig{initDir: PinOT, initVal: PinHI, closeDir: PinOT, closeVal: PinHI}, // D3 CS
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D4 GPIO
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D5 GPIO
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D6 GPIO
		&spiPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D7 GPIO
	} {
		pin |= D(i).spiPin(cfg)
	}
	return pin
}

// spiXferOption stores the various SPI transfer options as a 32-bit bitmap.
type spiXferOption uint32

// Constants defining the various SPI transfer options.
const (
	spiXferBytes spiXferOption = 0x00000000 // size is provided in bytes
	spiXferBits  spiXferOption = 0x00000001 // size is provided in bits

	spiCSManual   spiXferOption = 0x00000000
	spiCSAssert   spiXferOption = 0x00000002 // assert CS before start
	spiCSDeAssert spiXferOption = 0x00000004 // deassert CS after end

	// default transfer options
	spiXferDefault = spiXferBytes | spiCSManual
)

// Change changes the currently configured CS pin.
// It can be called while the SPI interface is open without having to first
// close and reopen the device.
// The CS pin can be either a DBUS or CBUS (GPIO) pin, see the documentation on
// Write for details.
func (spi *SPI) Change(cs Pin) error {

	// clear current CS selection
	spi.config.options &= ^(spiCSMask)

	if cs.IsMPSSE() {
		csOpt := cs.(DPin).spiOptionCS()
		if !csOpt.Valid() {
			return fmt.Errorf("invalid CS pin: %s [%08b][%d]", cs, csOpt, csOpt)
		}
		spi.config.options |= csOpt
	} else {
		// no changes necessary for CS on GPIO pin
	}

	spi.config.chipSelect = cs // update only if we didnt return early on error

	// only invoke the driver if we have an active SPI channel. otherwise, these
	// options get set on next Init().
	if ModeSPI == spi.device.mode {
		if err := _SPI_Change(spi); nil != err {
			return err
		}
	}

	return nil
}

// Option changes the dynamic configuration parameters of the SPI interface.
// It can be called while the SPI interface is open without having to first
// close and reopen the device.
func (spi *SPI) Option(opt *SPIOption) error {

	activeOpt := spiCSActiveHigh
	if opt.ActiveLow {
		activeOpt = spiCSActiveLow
	}

	modeOpt := spiOption(opt.Mode)
	if modeOpt > spiModeMask {
		return fmt.Errorf("invalid SPI mode: Mode %d", opt.Mode)
	}

	spi.config.options = activeOpt | modeOpt

	return spi.Change(opt.CS)
}

// Config initializes the SPI interface with the given configuration to a state
// ready for read/write.
// If the given configuration is nil, the default configuration is used (see
// SPIConfigDefault).
// It is not necessary to call Init after calling Config.
// See documentation of Init for other semantics.
func (spi *SPI) Config(cfg *SPIConfig) error {

	if nil == cfg {
		cfg = SPIConfigDefault()
	}

	if 0 == cfg.Clock {
		spi.config.clockRate = SPIClockDefault
	} else {
		if cfg.Clock <= SPIClockMaximum {
			spi.config.clockRate = cfg.Clock
		} else {
			return fmt.Errorf("invalid clock rate: %d", cfg.Clock)
		}
	}

	if 0 == cfg.Latency {
		spi.config.latency = SPILatencyDefault
	} else {
		spi.config.latency = cfg.Latency
	}

	if err := spi.Option(cfg.SPIOption); nil != err {
		return err
	}

	return spi.Init()
}

// Init initializes the SPI interface to a state ready for read/write.
// If Config has not been called, the default configuration is used (see
// SPIConfigDefault).
// If the interface is already initialized, it is first closed before
// initializing the interface.
func (spi *SPI) Init() error {

	if err := _SPI_InitChannel(spi); nil != err {
		return err
	}

	spi.device.mode = ModeSPI

	return spi.device.GPIO.Init() // reset GPIO
}

// Close closes both the SPI interface and the connection to the FT232H device.
func (spi *SPI) Close() error {
	return spi.device.Close()
}

// Write writes the given data to the SPI interface, optionally asserting the CS
// line before writing and/or de-asserting the CS line after writing.
// There is no maximum length for the data slice.
//
// The CS pin may be either a DBUS pin or CBUS (GPIO) pin. If it is a DBUS pin,
// then the MPSSE engine automatically handles CS assertion before and after
// writing, depending on the given flags start and stop. If it is a CBUS pin,
// then the GPIO pin is automatically set and cleared depending on the given
// flags start and stop. In both cases, the current value of the ActiveLow flag
// in the SPI configuration options determines if the CS line driven LOW
// (ActiveLow true, DEFAULT) or HIGH (ActiveLow false) when asserting and then
// de-asserting.
func (spi *SPI) Write(data []uint8, start bool, stop bool) (uint32, error) {

	cs := spi.config.chipSelect
	opt := spiXferDefault
	ass := 0 == uint32(spiCSActiveLow&spi.config.options)

	if start {
		if cs.IsMPSSE() {
			opt |= spiCSAssert
		} else {
			if err := spi.device.GPIO.Set(cs.(CPin), ass); nil != err {
				return 0, err
			}
		}
	}

	if stop {
		if cs.IsMPSSE() {
			opt |= spiCSDeAssert
		} else {
			// deassert on return
			defer func() { spi.device.GPIO.Set(cs.(CPin), !ass) }()
		}
	}

	return _SPI_Write(spi, data, opt)
}

// Write writes the given data to the SPI interface, optionally asserting the CS
// line before writing and/or de-asserting the CS line after writing using the
// given CS pin.
// There is no maximum length for the data slice.
// If the given CS pin is not the same as the currently configured CS pin, the
// CS configuration is changed and persists after writing.
// The CS pin can be either a DBUS or CBUS (GPIO) pin, see the documentation on
// Write for details.
func (spi *SPI) WriteTo(cs Pin, data []uint8, start bool, stop bool) (uint32, error) {

	if (start || stop) && !cs.Equals(spi.config.chipSelect) {
		// change if we are writing to a slave different than currently configured
		if err := spi.Change(cs); nil != err {
			return 0, err
		}
	}
	return spi.Write(data, start, stop)
}
