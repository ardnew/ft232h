package ft232h

import (
	"fmt"
)

// GPIO stores interface configuration settings for the GPIO ("C" port) and
// provides methods for reading and writing to GPIO pins.
// The GPIO interface is always initialized and available in any mode.
type GPIO struct {
	device *FT232H
	config *GPIOConfig
}

// GPIOConfig stores the most-recently read/written pin levels and directions.
type GPIOConfig struct {
	Dir uint8
	Val uint8
}

func (c *GPIOConfig) String() string {
	sym := func(i uint) rune {
		out := (c.Dir & (1 << i)) > 0
		hi := (c.Val & (1 << i)) > 0
		if out {
			if hi {
				return '^'
			} else {
				return '_'
			}
		} else {
			if hi {
				return '1'
			} else {
				return '0'
			}
		}
	}
	str := make([]rune, NumCPins)
	for i := range str {
		str[NumCPins-i-1] = sym(uint(i))
	}
	return string(str)
}

func (gpio *GPIO) String() string {
	return fmt.Sprintf("{ FT232H: %p, Config: %q }", gpio.device, gpio.config)
}

// GPIOConfigDefault returns the default pin levels and directions for the GPIO
// interface. All pins are configured as inputs at logic level LOW by default.
func GPIOConfigDefault() *GPIOConfig {
	return &GPIOConfig{
		Dir: 0x00, // each bit clear, all pins INPUT by default
		Val: 0x00, // each bit clear, all pins LOW by default
	}
}

// Write changes all pin direction and value configurations.
// This does not transfer any changes to the GPIO interface.
func (cfg *GPIOConfig) Write(dir uint8, val uint8) {
	cfg.Dir, cfg.Val = dir, val
}

// Set changes the pin direction and value configuration.
// This does not transfer any changes to the GPIO interface.
func (cfg *GPIOConfig) Set(pin CPin, dir Dir, val bool) error {
	if !pin.Valid() {
		return fmt.Errorf("invalid pin: %v", pin)
	}
	switch dir {
	case Output:
		cfg.Dir |= pin.Mask()
	case Input:
		cfg.Dir &= ^pin.Mask()
	}
	if val {
		cfg.Val |= pin.Mask()
	} else {
		cfg.Val &= ^pin.Mask()
	}
	return nil
}

// Init resets all GPIO pin directions and values using the most recently read
// or written configuration, returning a non-nil error if unsuccessful.
func (gpio *GPIO) Init() error {
	return gpio.Config(gpio.config)
}

// Config configures all GPIO pin directions and values to the settings defined
// in the given cfg, returning a non-nil error if unsuccessful.
func (gpio *GPIO) Config(cfg *GPIOConfig) error {
	gpio.config.Write(cfg.Dir, cfg.Val)
	return gpio.Write(cfg.Val)
}

// ConfigPin configures the given GPIO pin direction and value.
// The direction and value of all other pins is set based on the most recently
// read or written configuration determined prior to this call, and are all
// updated during this call.
// If you need more fine-grained control, use Read()/Write() directly.
func (gpio *GPIO) ConfigPin(pin CPin, dir Dir, val bool) error {
	if err := gpio.config.Set(pin, dir, val); nil != err {
		return err
	}
	return gpio.Write(gpio.config.Val)
}

// Write sets the value of all output pins at once using the given bitmask val,
// returning a non-nil error if unsuccessful.
func (gpio *GPIO) Write(val uint8) error {

	dir := gpio.config.Dir
	val &= dir // set only the pins configured as OUTPUT
	err := _FT_WriteGPIO(gpio, dir, val)
	if nil != err {
		return err
	}
	gpio.config.Val = val
	return nil
}

// Read returns the current value of all GPIO pins, returning 0 and a non-nil
// error if unsuccessful.
func (gpio *GPIO) Read() (uint8, error) {

	val, err := _FT_ReadGPIO(gpio)
	if nil != err {
		return 0, err
	}
	gpio.config.Val = val
	return val, nil
}

// Set sets the given pin to output with the given val.
// See ConfigPin() for other semantics.
func (gpio *GPIO) Set(pin CPin, val bool) error {
	return gpio.ConfigPin(pin, Output, val)
}

// Get reads the current value of the given pin.
func (gpio *GPIO) Get(pin CPin) (bool, error) {
	set, err := gpio.Read()
	if nil != err {
		return false, err
	}
	return (set & pin.Mask()) > 0, nil
}

// Chdir changes the GPIO direction of the given pin.
// Use ConfigPin() to change both direction and value, or Config() to change all
// pin directions (and values).
func (gpio *GPIO) Chdir(pin CPin, dir Dir) error {
	return gpio.ConfigPin(pin, dir, (gpio.config.Val&pin.Mask()) > 0)
}
