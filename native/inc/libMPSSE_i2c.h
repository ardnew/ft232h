/*!
 * \file libMPSSE_i2c.h
 *
 * \author FTDI
 * \date 20110505
 *
 * Copyright © 2000-2014 Future Technology Devices International Limited
 *
 *
 * THIS SOFTWARE IS PROVIDED BY FUTURE TECHNOLOGY DEVICES INTERNATIONAL LIMITED ``AS IS'' AND ANY EXPRESS
 * OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL FUTURE TECHNOLOGY DEVICES INTERNATIONAL LIMITED
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
 * BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 *
 * Project: libMPSSE
 * Module: I2C
 *
 * Rivision History:
 * 0.1 - initial version
 * 0.2 - 20110708 - added FT_ReadGPIO, FT_WriteGPIO & 3-phase clocking
 * 0.3 - 20111025 - modified for supporting 64bit linux
 *				    added I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE
 * 0.5 - 20140912 - modified for compilation issues with either C application/C++ application
 * 0.6 - 20202402 - add support for static linking
 */

#ifndef LIBMPSSE_I2C_H
#define LIBMPSSE_I2C_H

#include "ftd2xx.h"


/******************************************************************************/
/*								Macro defines								  */
/******************************************************************************/

#ifdef __cplusplus
#define FTDI_API extern "C"
#else
#define FTDI_API
#endif

/* Options to I2C_DeviceWrite & I2C_DeviceRead */
/*Generate start condition before transmitting */
#define	I2C_TRANSFER_OPTIONS_START_BIT		0x00000001

/*Generate stop condition before transmitting */
#define I2C_TRANSFER_OPTIONS_STOP_BIT		0x00000002

/*Continue transmitting data in bulk without caring about Ack or nAck from device if this bit is
not set. If this bit is set then stop transitting the data in the buffer when the device nAcks*/
#define I2C_TRANSFER_OPTIONS_BREAK_ON_NACK	0x00000004

/* libMPSSE-I2C generates an ACKs for every byte read. Some I2C slaves require the I2C
master to generate a nACK for the last data byte read. Setting this bit enables working with such
I2C slaves */
#define I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE	0x00000008

/* no address phase, no USB interframe delays */
#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES	0x00000010
#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS	0x00000020
#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER		0x00000030

/* if I2C_TRANSFER_OPTION_FAST_TRANSFER is set then setting this bit would mean that the
address field should be ignored. The address is either a part of the data or this is a special I2C
frame that doesn't require an address*/
#define I2C_TRANSFER_OPTIONS_NO_ADDRESS		0x00000040

#define I2C_CMD_GETDEVICEID_RD	0xF9
#define I2C_CMD_GETDEVICEID_WR	0xF8

#define I2C_GIVE_ACK	1
#define I2C_GIVE_NACK	0

/* 3-phase clocking is enabled by default. Setting this bit in ConfigOptions will disable it */
#define I2C_DISABLE_3PHASE_CLOCKING	0x0001

/* The I2C master should actually drive the SDA line only when the output is LOW. It should be
tristate the SDA line when the output should be high. This tristating the SDA line during output
HIGH is supported only in FT232H chip. This feature is called DriveOnlyZero feature and is
enabled when the following bit is set in the options parameter in function I2C_Init */
#define I2C_ENABLE_DRIVE_ONLY_ZERO	0x0002



/******************************************************************************/
/*								Type defines								  */
/******************************************************************************/

typedef unsigned char   uint8;
typedef unsigned short  uint16;
typedef unsigned long long uint64;

typedef signed char   int8;
typedef signed short  int16;
typedef signed long long int64;

#ifndef __cplusplus
typedef unsigned char	bool;
#endif

#ifdef __x86_64__
	typedef unsigned int   uint32;
	typedef signed int   int32;
#else
	typedef unsigned long   uint32;
	typedef signed long   int32;
#endif

typedef enum I2C_ClockRate_t{
	I2C_CLOCK_STANDARD_MODE = 100000,						/* 100kb/sec */
	I2C_CLOCK_FAST_MODE = 400000,							/* 400kb/sec */
	I2C_CLOCK_FAST_MODE_PLUS = 1000000, 					/* 1000kb/sec */
	I2C_CLOCK_HIGH_SPEED_MODE = 3400000 					/* 3.4Mb/sec */
}I2C_CLOCKRATE;


/* Channel configuration information */
typedef struct I2C_ChannelConfig_t
{
	I2C_CLOCKRATE	ClockRate;
	uint8			LatencyTimer;
	uint32			Options;
}I2C_ChannelConfig;


/******************************************************************************/
/*								External variables							  */
/******************************************************************************/





/******************************************************************************/
/*								Function declarations						  */
/******************************************************************************/
FTDI_API FT_STATUS I2C_GetNumChannels(uint32 *numChannels);
FTDI_API FT_STATUS I2C_GetChannelInfo(uint32 index,
	FT_DEVICE_LIST_INFO_NODE *chanInfo);
FTDI_API FT_STATUS I2C_OpenChannel(uint32 index, FT_HANDLE *handle);
FTDI_API FT_STATUS I2C_InitChannel(FT_HANDLE handle, I2C_ChannelConfig *config);
FTDI_API FT_STATUS I2C_CloseChannel(FT_HANDLE handle);
FTDI_API FT_STATUS I2C_DeviceRead(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransfered, uint32 options);
FTDI_API FT_STATUS I2C_DeviceWrite(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransfered, uint32 options);
FTDI_API void Init_libMPSSE(void);
FTDI_API void Cleanup_libMPSSE(void);
FTDI_API FT_STATUS FT_WriteGPIO(FT_HANDLE handle, uint8 dir, uint8 value);
FTDI_API FT_STATUS FT_ReadGPIO(FT_HANDLE handle,uint8 *value);




/******************************************************************************/


#endif	/*LIBMPSSE_I2C_H*/

