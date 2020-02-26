package ft232h

import "fmt"

type SPI struct {
	device *FT232H
	config *spiConfig
}

// spiConfig holds all of the configuration settings for an SPI channel.
type spiConfig struct {
	clockRate  uint32 // in Hertz
	latency    uint8  // in ms
	options    spiOption
	pin        uint32 // port D pins ("low byte lines of MPSSE")
	chipSelect Pin    // may be DPin (MPSSE low byte) or CPin (GPIO)
}

func spiConfigDefault() *spiConfig {
	return &spiConfig{
		clockRate:  SPIClockDefault,
		latency:    SPILatencyDefault,
		options:    spiCSActiveDefault | spiCSDefault | spiModeDefault,
		pin:        spiPinConfigDefault(),
		chipSelect: spiCSDefault.pin(),
	}
}

type spiOption uint32

// Constants defining the available options in the SPI configuration struct.
const (
	spiOptionInvalid spiOption = 0xAAAAAAAA

	// Known SPI operating modes
	//   LIMITATION: libMPSSE only supports mode 0 and mode 2 (CPHA==2).
	SPIMode0       spiOption = 0x00000000 // capture on RISE, propagate on FALL
	SPIMode1       spiOption = 0x00000001 // capture on FALL, propagate on RISE
	SPIMode2       spiOption = 0x00000002 // capture on FALL, propagate on RISE
	SPIMode3       spiOption = 0x00000003 // capture on RISE, propagate on FALL
	spiModeMask    spiOption = 0x00000003
	spiModeDefault spiOption = SPIMode0

	// DPins available for chip-select operation
	spiCSD3      spiOption = 0x00000000 // SPI CS on D3
	spiCSD4      spiOption = 0x00000004 // SPI CS on D4
	spiCSD5      spiOption = 0x00000008 // SPI CS on D5
	spiCSD6      spiOption = 0x0000000C // SPI CS on D6
	spiCSD7      spiOption = 0x00000010 // SPI CS on D7
	spiCSMask    spiOption = 0x0000001C
	spiCSDefault spiOption = spiCSD3

	// Other options
	spiCSActiveLow     spiOption = 0x00000020 // drive pin low to assert CS
	spiCSActiveHigh    spiOption = 0x00000000 // drive pin high to assert CS
	spiCSActiveDefault spiOption = spiCSActiveLow
)

// Valid verifies opt isn't equal to the constant for invalid SPI options
func (opt spiOption) Valid() bool { return opt != spiOptionInvalid }

// pin translates a chip-select mask opt from an SPI configuration struct option
// to its corresponding DPin.
func (opt spiOption) pin() DPin {
	if !opt.Valid() {
		return DPin(0)
	}
	switch opt {
	case spiCSD3, spiCSD4, spiCSD5, spiCSD6, spiCSD7:
		return D(int(opt>>2) + 3)
	default:
		return DPin(0)
	}
}

// spiOption translates a DPin p to its corresponding chip-select mask for an
// SPI configuration struct option.
func (p DPin) spiOption() spiOption {
	if p.Valid() && p >= 3 {
		return (spiOption(p.Pos()) << 2) & spiCSMask
	} else {
		return spiOptionInvalid
	}
}

// Constants related to board pins when MPSSE operating in SPI mode
const (
	SPIClockMaximum   uint32 = 30000000
	SPIClockDefault   uint32 = 12000000 // valid range: 0-30000000 (30 MHz)
	SPILatencyDefault byte   = 16       // 1-255 USB HiSpeed, 2-255 USB FullSpeed
)

// spiPinConfig represents the default direction and value for pins associated
// with the lower byte lines of MPSSE, reserved for serial functions SPI/I²C
// (or port "D" on FT232H), but has a few GPIO pins as well.
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

// spiPinConfigDefault defines the initial spiPinConfig value for all pins
// represented by this type. all output pins are configured LOW except for the
// default CS pin (D3) since we also have spiCSActiveLow by default. this means
// we won't activate the default slave line until intended. it also means SCLK
// idles LOW (change initVal to PinHI to idle HIGH).
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

type spiXferOption uint32

// Constants controlling the supported SPI transfer options
const (
	spiXferBytes spiXferOption = 0x00000000 // size is provided in bytes
	spiXferBits  spiXferOption = 0x00000001 // size is provided in bits

	spiCSManual   spiXferOption = 0x00000000
	spiCSAssert   spiXferOption = 0x00000002 // assert CS before start
	spiCSDeAssert spiXferOption = 0x00000004 // deassert CS after end

	// default transfer options
	spiXferDefault
)

//func SPIWarningCSGPIO(pin Pin) error {
//	return fmt.Errorf("SPI chip-select on GPIO pin %s", pin)
//}

func (spi *SPI) ChangeCS(cs Pin) error {

	if (cs.IsMPSSE() == spi.config.chipSelect.IsMPSSE()) &&
		(cs.Mask() == spi.config.chipSelect.Mask()) {
		return nil // no change, provided pin is already CS
	}

	// clear current CS selection
	spi.config.options &= ^(spiCSMask)

	if cs.IsMPSSE() {
		csOpt := cs.(DPin).spiOption()
		if !csOpt.Valid() {
			return fmt.Errorf("invalid CS pin: %d", cs)
		}
		spi.config.options |= csOpt
		// only invoke the driver if we have an active SPI channel. otherwise, these
		// options get set on next Init().
		if ModeSPI == spi.device.mode {
			if err := _SPI_ChangeCS(spi); nil != err {
				return err
			}
		}
	} else {
		// no changes necessary for CS on GPIO pin
	}

	spi.config.chipSelect = cs // update only if we didnt return early on error

	return nil
}

func (spi *SPI) SetOptions(cs Pin, activeLow bool, mode byte) error {

	var (
		activeOpt spiOption
		modeOpt   spiOption
	)

	if activeLow {
		activeOpt = spiCSActiveLow
	} else {
		activeOpt = spiCSActiveHigh
	}

	if spiOption(mode) > spiModeMask {
		return fmt.Errorf("invalid SPI mode: Mode %d", mode)
	} else {
		modeOpt = spiOption(mode)
	}

	spi.config.options = activeOpt | modeOpt

	return spi.ChangeCS(cs)
}

func (spi *SPI) SetConfig(clock uint32, latency byte, cs Pin, activeLow bool, mode byte) error {

	if 0 == clock {
		spi.config.clockRate = SPIClockDefault
	} else {
		if clock <= SPIClockMaximum {
			spi.config.clockRate = clock
		} else {
			return fmt.Errorf("invalid clock rate: %d", clock)
		}
	}

	if 0 == latency {
		spi.config.latency = SPILatencyDefault
	} else {
		spi.config.latency = latency
	}

	return spi.SetOptions(cs, activeLow, mode)
}

func (spi *SPI) Init() error {

	if err := _SPI_InitChannel(spi); nil != err {
		return err
	}

	spi.device.mode = ModeSPI

	return spi.device.GPIO.Init() // reset GPIO
}

func (spi *SPI) Close() error {
	return spi.device.Close()
}

func (spi *SPI) Write(data []uint8, start bool, stop bool) (uint32, error) {

	cs := spi.config.chipSelect
	opt := spiXferBytes
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

func (spi *SPI) WriteTo(cs Pin, data []uint8, start bool, stop bool) (uint32, error) {

	if start || stop {
		if err := spi.ChangeCS(cs); nil != err {
			return 0, err
		}
	}
	return spi.Write(data, start, stop)
}
