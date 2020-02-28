package ft232h

// #cgo darwin,amd64 LDFLAGS: -framework CoreFoundation -framework IOKit
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/native/lib/darwin-amd64
// #cgo linux,amd64  LDFLAGS: -L${SRCDIR}/native/lib/linux-amd64
// #cgo linux,arm64  LDFLAGS: -L${SRCDIR}/native/lib/linux-arm64
// #cgo linux,386    LDFLAGS: -L${SRCDIR}/native/lib/linux-386
// #cgo  CFLAGS: -I${SRCDIR}/native/inc
// #cgo LDFLAGS: -lft232h
// #include "libMPSSE_spi.h"
// #include "libMPSSE_i2c.h"
// #include "ftd2xx.h"
// #include "stdlib.h"
import "C"

// Type aliases for the native types needed by the C libraries.
type (
	Handle C.FT_HANDLE
	Status C.FT_STATUS
	Chip   C.FT_DEVICE
	Mode   C.int
)

// Constants related to device status
const (
	SOK                      Status = C.FT_OK
	SInvalidHandle           Status = C.FT_INVALID_HANDLE
	SDeviceNotFound          Status = C.FT_DEVICE_NOT_FOUND
	SDeviceNotOpened         Status = C.FT_DEVICE_NOT_OPENED
	SIOError                 Status = C.FT_IO_ERROR
	SInsufficientResources   Status = C.FT_INSUFFICIENT_RESOURCES
	SInvalidParameter        Status = C.FT_INVALID_PARAMETER
	SInvalidBaudRate         Status = C.FT_INVALID_BAUD_RATE
	SDeviceNotOpenedForErase Status = C.FT_DEVICE_NOT_OPENED_FOR_ERASE
	SDeviceNotOpenedForWrite Status = C.FT_DEVICE_NOT_OPENED_FOR_WRITE
	SFailedToWriteDevice     Status = C.FT_FAILED_TO_WRITE_DEVICE
	SEEPROMReadFailed        Status = C.FT_EEPROM_READ_FAILED
	SEEPROMWriteFailed       Status = C.FT_EEPROM_WRITE_FAILED
	SEEPROMEraseFailed       Status = C.FT_EEPROM_ERASE_FAILED
	SEEPROMNotPresent        Status = C.FT_EEPROM_NOT_PRESENT
	SEEPROMNotProgrammed     Status = C.FT_EEPROM_NOT_PROGRAMMED
	SInvalidArgs             Status = C.FT_INVALID_ARGS
	SNotSupported            Status = C.FT_NOT_SUPPORTED
	SOtherError              Status = C.FT_OTHER_ERROR
	SDeviceListNotReady      Status = C.FT_DEVICE_LIST_NOT_READY
)

// OK returns true if the status equals SOK, otherwise false.
func (s Status) OK() bool {
	return SOK == s
}

// Error returns the string representation of a status, required by the Go error
// interface. Returns the string "unknown error" is the status is invalid.
func (s Status) Error() string {
	switch s {
	case SOK:
		return "OK"
	case SInvalidHandle:
		return "invalid handle"
	case SDeviceNotFound:
		return "device not found"
	case SDeviceNotOpened:
		return "device not opened"
	case SIOError:
		return "IO error"
	case SInsufficientResources:
		return "insufficient resources"
	case SInvalidParameter:
		return "invalid parameter"
	case SInvalidBaudRate:
		return "invalid baud rate"
	case SDeviceNotOpenedForErase:
		return "device not opened for erase"
	case SDeviceNotOpenedForWrite:
		return "device not opened for write"
	case SFailedToWriteDevice:
		return "failed to write device"
	case SEEPROMReadFailed:
		return "EEPROM read failed"
	case SEEPROMWriteFailed:
		return "EEPROM write failed"
	case SEEPROMEraseFailed:
		return "EEPROM erase failed"
	case SEEPROMNotPresent:
		return "EEPROM not present"
	case SEEPROMNotProgrammed:
		return "EEPROM not programmed"
	case SInvalidArgs:
		return "invalid args"
	case SNotSupported:
		return "not supported"
	case SOtherError:
		return "other error"
	case SDeviceListNotReady:
		return "device list not ready"
	default:
		return "unknown error"
	}
}

// Constants defining the FTDI chip identifiers specified by FTDI.
const (
	CFTBM      Chip = C.FT_DEVICE_BM
	CFTAM      Chip = C.FT_DEVICE_AM
	CFT100AX   Chip = C.FT_DEVICE_100AX
	CFTUnknown Chip = C.FT_DEVICE_UNKNOWN
	CFT2232C   Chip = C.FT_DEVICE_2232C
	CFT232R    Chip = C.FT_DEVICE_232R
	CFT2232H   Chip = C.FT_DEVICE_2232H
	CFT4232H   Chip = C.FT_DEVICE_4232H
	CFT232H    Chip = C.FT_DEVICE_232H
	CFTX       Chip = C.FT_DEVICE_X_SERIES
	CFT4222H0  Chip = C.FT_DEVICE_4222H_0
	CFT4222H12 Chip = C.FT_DEVICE_4222H_1_2
	CFT4222H3  Chip = C.FT_DEVICE_4222H_3
	CFT4222P   Chip = C.FT_DEVICE_4222_PROG
	CFT900     Chip = C.FT_DEVICE_900
	CFT930     Chip = C.FT_DEVICE_930
	CUMFTPD3A  Chip = C.FT_DEVICE_UMFTPD3A
)

// String returns the descriptive string representation of an FTDI chip.
// Returns the string "invalid chip" if the chip is not defined.
func (c Chip) String() string {
	switch c {
	case CFTBM:
		return "FTBM"
	case CFTAM:
		return "FTAM"
	case CFT100AX:
		return "FT100AX"
	case CFTUnknown:
		return "FTUnknown"
	case CFT2232C:
		return "FT2232C"
	case CFT232R:
		return "FT232R"
	case CFT2232H:
		return "FT2232H"
	case CFT4232H:
		return "FT4232H"
	case CFT232H:
		return "FT232H"
	case CFTX:
		return "FTX"
	case CFT4222H0:
		return "FT4222H0"
	case CFT4222H12:
		return "FT4222H12"
	case CFT4222H3:
		return "FT4222H3"
	case CFT4222P:
		return "FT4222P"
	case CFT900:
		return "FT900"
	case CFT930:
		return "FT930"
	case CUMFTPD3A:
		return "UMFTPD3A"
	default:
		return "invalid chip"
	}
}

// Constants defining the legacy protocols supported by MPSSE.
const (
	ModeNone Mode = 0
	ModeSPI  Mode = 1
	ModeI2C  Mode = 2
)

// String returns a string describing the legacy protocol supported by MPSSE.
// Returns the string "unknown" if the mode is invalid.
func (m Mode) String() string {
	switch m {
	case ModeNone:
		return "None"
	case ModeSPI:
		return "SPI"
	case ModeI2C:
		return "I²C"
	default:
		return "unknown"
	}
}

// _FT_CreateDeviceInfoList requests the D2XX driver allocate and populate an
// internal list of MPSSE-capable USB devices connected to the system, returning
// the number of devices found if successful.
// Returns 0 and a non-nil error if the device list could not be created.
func _FT_CreateDeviceInfoList() (uint, error) {
	var n C.DWORD
	stat := Status(C.FT_CreateDeviceInfoList(&n))
	if !stat.OK() {
		return 0, stat
	}
	return uint(n), nil
}

// _FT_GetDeviceInfoList parses and returns a slice of deviceInfo pointers for
// all devices stored in the internal device list of the D2XX driver.
// Returns a nil slice and non-nil error if the device list could not be read.
// Returns an empty slice and nil error if no devices were found in the list.
func _FT_GetDeviceInfoList(n uint) ([]*deviceInfo, error) {
	ndev := C.DWORD(n)
	list := make([]C.FT_DEVICE_LIST_INFO_NODE, n)
	stat := Status(C.FT_GetDeviceInfoList(&list[0], &ndev))
	if !stat.OK() {
		return nil, stat
	}
	info := make([]*deviceInfo, n)
	for i, node := range list {
		// parse the C struct into our simpler Go definition
		info[i] = &deviceInfo{
			index:     i,
			isOpen:    1 == (node.Flags & 0x01),
			isHiSpeed: 2 == (node.Flags & 0x02),
			chip:      Chip(node.Type),
			vid:       (uint32(node.ID) >> 16) & 0xFFFF,
			pid:       (uint32(node.ID)) & 0xFFFF,
			locID:     uint32(node.LocId),
			serial:    C.GoString(&node.SerialNumber[0]),
			desc:      C.GoString(&node.Description[0]),
			handle:    Handle(node.ftHandle),
		}
	}
	return info, nil
}

// _FT_Open attempts to open a raw USB interface through the D2XX driver,
// returning a non-nil error if unsuccessful.
func _FT_Open(info *deviceInfo) error {
	stat := Status(C.FT_Open(C.int(info.index), (*C.PVOID)(&info.handle)))
	if !stat.OK() {
		return stat
	}
	return nil
}

// _FT_Close attempts to close a USB interface opened through the D2XX driver,
// returning a non-nil error if unsuccessful.
func _FT_Close(info *deviceInfo) error {
	stat := Status(C.FT_Close(C.PVOID(info.handle)))
	if !stat.OK() {
		return stat
	}
	return nil
}

// _FT_WriteGPIO sets the level val and direction dir for all pins on port "C"
// of the FT232H using the D2XX driver, returns a non-nil error if the driver
// could not set the pin configuration.
func _FT_WriteGPIO(gpio *GPIO, dir uint8, val uint8) error {
	stat := Status(C.FT_WriteGPIO(C.PVOID(gpio.device.info.handle), C.uint8(dir), C.uint8(val)))
	if !stat.OK() {
		return stat
	}
	return nil
}

// _FT_ReadGPIO reads the level of all pins on port "C" of the FT232H using the
// D2XX driver, returning 0 and a non-nil error if the pins could not be read.
func _FT_ReadGPIO(gpio *GPIO) (uint8, error) {
	var val C.uint8
	stat := Status(C.FT_ReadGPIO(C.PVOID(gpio.device.info.handle), &val))
	if !stat.OK() {
		return 0, stat
	}
	return uint8(val), nil
}

// _SPI_InitChannel initializes the MPSSE engine in SPI master mode with the
// configuration defined in the given spi using the libMPSSE driver.
// If the FT232H device is already opened in any mode (including SPI), the
// interface is first closed before re-opening with the new configuration.
// Returns a non-nil error if the interface could not be closed or (re)opened.
func _SPI_InitChannel(spi *SPI) error {

	// close any open channels before trying to init
	if err := spi.device.Close(); nil != err {
		return err
	}

	stat := Status(C.SPI_OpenChannel(C.uint32(spi.device.info.index),
		(*C.PVOID)(&spi.device.info.handle)))
	if !stat.OK() {
		return stat
	}

	config := C.SPI_ChannelConfig{
		ClockRate:     C.uint32(spi.config.clockRate),
		LatencyTimer:  C.uint8(spi.config.latency),
		configOptions: C.uint32(spi.config.options),
		Pin:           C.uint32(spi.config.pin),
		reserved:      C.uint16(0),
	}

	stat = Status(C.SPI_InitChannel(C.PVOID(spi.device.info.handle), &config))
	if !stat.OK() {
		return stat
	}

	return nil
}

// _SPI_Change reconfigures the dynamic interface parameters of an open SPI
// interface using the libMPSSE driver, returning a non-nil error if
// unsuccessful.
func _SPI_Change(spi *SPI) error {
	stat := Status(C.SPI_ChangeCS(C.PVOID(spi.device.info.handle),
		C.uint32(spi.config.options)))
	if !stat.OK() {
		return stat
	}
	return nil
}

// _SPI_Write performs an SPI write using the given open SPI interface and slice
// of uint8 data and transfer options using the libMPSSE driver, returning a
// non-nil error if unsuccessful.
// If the given data slice length is greater than UINT16_MAX (65536), multiple
// write requests are performed with the libMPSSE driver.
func _SPI_Write(spi *SPI, data []uint8, opt spiXferOption) (uint32, error) {

	// note that MPSSE has a limitation on the size of SPI transfers, since the
	// packet length has to fit into 16 bits, so the max transfer size is 65536.
	// we break up the buffer here to transmit as much as possible at once.
	const MaxTransferBytes = 65536
	var (
		total uint32
		sent  C.uint32
	)

	dataLen := len(data)
	for beg := 0; beg < dataLen; beg += MaxTransferBytes {
		end := beg + MaxTransferBytes
		if end > dataLen {
			end = dataLen
		}
		stat := Status(C.SPI_Write(C.PVOID(spi.device.info.handle),
			(*C.uint8)(&data[beg]), C.uint32(end-beg), &sent, C.uint32(opt)))
		total += uint32(sent)
		if !stat.OK() {
			return total, stat
		}
	}
	return total, nil
}
