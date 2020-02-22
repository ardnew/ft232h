package ft232h

type GPIO struct {
	device *FT232H
	config *gpioConfig
}

type gpioConfig struct {
	dir uint8
	val uint8
}

func gpioConfigDefault() *gpioConfig {
	return &gpioConfig{
		dir: 0xFF, // each bit set, all pins OUTPUT by default
		val: 0x00, // each bit clear, all pins LOW by default
	}
}

func (gpio *GPIO) Init() error {
	return gpio.Write(gpio.config.dir, gpio.config.val)
}

func (gpio *GPIO) Write(dir uint8, val uint8) error {

	val &= dir // only set output bits

	if err := _FT_WriteGPIO(gpio, dir, val); nil != err {
		return err
	}

	gpio.config.dir = dir
	gpio.config.val = val

	return nil
}

func (gpio *GPIO) Read() (uint8, error) {

	val, err := _FT_ReadGPIO(gpio)
	if nil != err {
		return 0, err
	}

	gpio.config.val = val

	return val, nil
}

func (gpio *GPIO) Set(pin CPin, val bool) error {

	dir := gpio.config.dir | uint8(pin)
	set := gpio.config.val

	if val {
		set |= uint8(pin)
	} else {
		set &= ^uint8(pin)
	}

	return gpio.Write(dir, set)
}

func (gpio *GPIO) Get(pin CPin) (bool, error) {

	set, err := gpio.Read()
	if nil != err {
		return false, err
	}

	return (set & uint8(pin)) > 0, nil
}
