package ft232h

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

// GPIOConfigDefault returns the default pin levels and directions for the GPIO
// interface. All pins are configured as outputs at logic level LOW by default.
func GPIOConfigDefault() *GPIOConfig {
	return &GPIOConfig{
		Dir: 0xFF, // each bit set, all pins OUTPUT by default
		Val: 0x00, // each bit clear, all pins LOW by default
	}
}

// Init resets all GPIO pin directions and values using the most recently read
// or written configuration, returning a non-nil error if unsuccessful.
func (gpio *GPIO) Init() error {
	return gpio.Config(gpio.config)
}

// Config configures all GPIO pin directions and values to the settings defined
// in the given cfg, returning a non-nil error if unsuccessful.
func (gpio *GPIO) Config(cfg *GPIOConfig) error {
	gpio.config.Dir = cfg.Dir
	return gpio.Write(cfg.Val)
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

// Set sets the given pin to output with the given val, returning a non-nil
// error if unsuccessful.
// The direction and value of all other pins is set based on the most recently
// read or written configuration determined prior to calling Set, and are all
// updated during the call to Set (AFTER writing the configuration). If you need
// more fine-grained control, use Read/Write instead.
func (gpio *GPIO) Set(pin CPin, val bool) error {

	dir := gpio.config.Dir | uint8(pin)
	set := gpio.config.Val

	if val {
		set |= uint8(pin)
	} else {
		set &= ^uint8(pin)
	}

	return gpio.Config(&GPIOConfig{dir, set})
}

// Get reads the value of the given pin, returning a non-nil error if
// unsuccessful.
func (gpio *GPIO) Get(pin CPin) (bool, error) {

	set, err := gpio.Read()
	if nil != err {
		return false, err
	}

	return (set & uint8(pin)) > 0, nil
}
