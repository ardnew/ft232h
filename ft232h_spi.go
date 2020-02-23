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
	chipSelect Pin
}

func spiConfigDefault() *spiConfig {
	return &spiConfig{
		clockRate:  SPIClockDefault,
		latency:    SPILatencyDefault,
		options:    spiCSActiveDefault | spiCSDefault | spiModeDefault,
		pin:        spiDPinConfigDefault(),
		chipSelect: spiCSDefault.pin(),
	}
}

type spiOption uint32

// Constants defining the available options in the SPI configuration struct.
const (
	// Known SPI operating modes
	//   LIMITATION: libMPSSE only supports mode 0 and mode 2 (CPHA==2).
	SPIMode0       spiOption = 0x00000000 // capture on RISE, propagate on FALL
	SPIMode1       spiOption = 0x00000001 // capture on FALL, propagate on RISE
	SPIMode2       spiOption = 0x00000002 // capture on FALL, propagate on RISE
	SPIMode3       spiOption = 0x00000003 // capture on RISE, propagate on FALL
	spiModeMask    spiOption = 0x00000003
	spiModeInvalid spiOption = 0x000000FF
	spiModeDefault spiOption = SPIMode0

	// DPins available for chip-select operation
	spiCSD3      spiOption = 0x00000000 // SPI CS on D3
	spiCSD4      spiOption = 0x00000004 // SPI CS on D4
	spiCSD5      spiOption = 0x00000008 // SPI CS on D5
	spiCSD6      spiOption = 0x0000000C // SPI CS on D6
	spiCSD7      spiOption = 0x00000010 // SPI CS on D7
	spiCSMask    spiOption = 0x0000001C
	spiCSInvalid spiOption = 0x000000FF
	spiCSDefault spiOption = spiCSD3

	// Other options
	spiCSActiveLow     spiOption = 0x00000020 // drive pin low to assert CS
	spiCSActiveHigh    spiOption = 0x00000000 // drive pin high to assert CS
	spiCSActiveDefault spiOption = spiCSActiveLow
)

func (opt spiOption) pin() DPin {
	switch opt {
	case spiCSD3:
		return D3
	case spiCSD4:
		return D4
	case spiCSD5:
		return D5
	case spiCSD6:
		return D6
	case spiCSD7:
		return D7
	default:
		return spiCSDefault.pin()
	}
}

// Constants related to board pins when MPSSE operating in SPI mode
const (
	SPIClockMaximum   uint32 = 30000000
	SPIClockDefault   uint32 = 12000000 // valid range: 0-30000000 (30 MHz)
	SPILatencyDefault byte   = 16       // 1-255 USB HiSpeed, 2-255 USB FullSpeed
)

// spiCSPin translates a DPin value to its corresponding chip-select mask for
// the SPI configuration struct option.
var spiCSPin = map[DPin]spiOption{
	D0: spiCSInvalid,
	D1: spiCSInvalid,
	D2: spiCSInvalid,
	D3: spiCSD3,
	D4: spiCSD4,
	D5: spiCSD5,
	D6: spiCSD6,
	D7: spiCSD7,
}

// spiDPinConfig represents the default direction and value for pins associated
// with the lower byte lines of MPSSE, reserved for serial functions SPI/I²C
// (or port "D" on FT232H), but has a few GPIO pins as well.
type spiDPinConfig struct {
	initDir  byte // direction of lines after SPI channel initialization
	initVal  byte // value of lines after SPI channel initialization
	closeDir byte // direction of lines after SPI channel is closed
	closeVal byte // value of lines after SPI channel is closed
}

// spiDPinConfigDefault defines the initial spiDPinConfig value for all pins
// represented by this type. all output pins are configured LOW except for the
// default CS pin (D3) since we also have spiCSActiveLow by default. this means
// we won't activate the default slave line until intended. it also means SCLK
// idles LOW (change initVal to PinHI to idle HIGH).
func spiDPinConfigDefault() uint32 {
	return spiDPin([NumDPins]*spiDPinConfig{
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D0 SCLK
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D1 MOSI
		&spiDPinConfig{initDir: PinIN, initVal: PinLO, closeDir: PinIN, closeVal: PinLO}, // D2 MISO
		&spiDPinConfig{initDir: PinOT, initVal: PinHI, closeDir: PinOT, closeVal: PinHI}, // D3 CS
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D4 GPIO
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D5 GPIO
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D6 GPIO
		&spiDPinConfig{initDir: PinOT, initVal: PinLO, closeDir: PinOT, closeVal: PinLO}, // D7 GPIO
	})
}

// spiDPin constructs the 32-bit field pin of the spiConfig struct from the
// provided spiDPinConfig slice cfg for each pin (identified by its index in the
// given slice).
func spiDPin(cfg [NumDPins]*spiDPinConfig) uint32 {
	var pin uint32
	for i, c := range cfg {
		pin |= (uint32(c.initDir) << i) | (uint32(c.initVal) << (8 + i)) |
			(uint32(c.closeDir) << (16 + i)) | (uint32(c.closeVal) << (24 + i))
	}
	return pin
}

type spiXferOption uint32

// Constants controlling the supported SPI transfer options
const (
	spiXferBytes spiXferOption = 0x00000000 // size is provided in bytes
	spiXferBits  spiXferOption = 0x00000001 // size is provided in bits

	spiCSAssert   spiXferOption = 0x00000002 // assert CS before start
	spiCSDeAssert spiXferOption = 0x00000004 // deassert CS after end
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
		var csOpt spiOption
		if csOpt, ok := spiCSPin[cs.(DPin)]; !ok || (spiCSInvalid == csOpt) {
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
			defer func() { _ = spi.device.GPIO.Set(cs.(CPin), !ass) }()
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
