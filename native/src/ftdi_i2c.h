/*!
 * \file FTDI_I2C.h
 *
 * \author FTDI
 * \date 20110127
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
 * 0.2 - 20110708 - added macro I2C_DISABLE_3PHASE_CLOCKING
 */

#ifndef FTDI_I2C_H
#define FTDI_I2C_H

#include "ftdi_infra.h"


/******************************************************************************/
/*								Macro defines								  */
/******************************************************************************/
/* Options to I2C_DeviceWrite & I2C_DeviceRead */
	/*Generate start condition before transmitting */
	#define	I2C_TRANSFER_OPTIONS_START_BIT		0x00000001

	/*Generate stop condition before transmitting */
	#define I2C_TRANSFER_OPTIONS_STOP_BIT		0x00000002

	/*Continue transmitting data in bulk without caring about Ack or nAck from device if this bit
	is not set. If this bit is set then stop transferring the data in the buffer when the device
	nACKs */
	#define I2C_TRANSFER_OPTIONS_BREAK_ON_NACK	0x00000004

	/* libMPSSE-I2C generates an ACKs for every byte read. Some I2C slaves require the I2C
	master to generate a nACK for the last data byte read. Setting this bit enables working with
	such I2C slaves */
	#define I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE	0x00000008

	/*Fast transfers prepare a buffer containing commands to generate START/STOP/ADDRESS
	   conditions and commands to read/write data. This buffer is sent to the MPSSE in one shot,
	   hence delays between different phases of the I2C transfer are eliminated. Fast transfers
	   can have data length in terms of bits or bytes. The user application should call
	   I2C_DeviceWrite or I2C_DeviceRead with either
	   I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES or
	   I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS bit set to perform a fast transfer.
	   I2C_TRANSFER_OPTIONS_START_BIT and I2C_TRANSFER_OPTIONS_STOP_BIT have
	   their usual meanings when used in fast transfers, however
	   I2C_TRANSFER_OPTIONS_BREAK_ON_NACK and
	   I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE are not applicable in fast transfers */
	#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER		0x00000030/*not visible to user*/

	/* When the user calls I2C_DeviceWrite or I2C_DeviceRead with this bit set then libMPSSE
	     packs commands to transfer sizeToTransfer number of bytes, and to read/write
	     sizeToTransfer number of ack bits. If data is written then the read ack bits are ignored, if
	     data is being read then an acknowledgement bit(SDA=LOW) is given to the I2C slave
	     after each byte read */
	#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES	0x00000010

	/* When the user calls I2C_DeviceWrite or I2C_DeviceRead with this bit set then libMPSSE
	     packs commands to transfer sizeToTransfer number of bits. There is no ACK phase when
	     this bit is set */
	#define I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS	0x00000020

	/* The address parameter is ignored in fast transfers if this bit is set. This would mean that
	     the address is either a part of the data or this is a special I2C frame that doesn't require
	     an address. However if this bit is not set then 7bit address and 1bit direction will be
	     written to the I2C bus each time I2C_DeviceWrite or I2C_DeviceRead is called and a
	     1bit acknowledgement will be read after that, which will however be just ignored */
	#define I2C_TRANSFER_OPTIONS_NO_ADDRESS		0x00000040


#define I2C_CMD_GETDEVICEID_RD	0xF9
#define I2C_CMD_GETDEVICEID_WR	0xF8

#define I2C_GIVE_ACK	1
#define I2C_GIVE_NACK	0

/* 3-phase clocking is enabled by default. Setting this bit in ConfigOptions will disable it */
#define I2C_DISABLE_3PHASE_CLOCKING	0x0001



/******************************************************************************/
/*								Type defines								  */
/******************************************************************************/
/*
Valid range for clock divisor is 0 to 65535
Highest clock freq is 6MHz represented by 0
The next highest is   3MHz represented by 1
Lowest is                 91Hz  represented by 65535

User can pass either pass I2C_DataRate_100K, I2C_DataRate_400K or
I2C_DataRate_3_4M for the standard clock rates
or a clock divisor value may be passed
*/
typedef enum I2C_ClockRate_t{
	I2C_CLOCK_STANDARD_MODE = 100000,		/* 100kb/sec */
	I2C_CLOCK_FAST_MODE = 400000,			/* 400kb/sec */
	I2C_CLOCK_FAST_MODE_PLUS = 1000000,		/* 1000kb/sec */
	I2C_CLOCK_HIGH_SPEED_MODE = 3400000		/* 3.4Mb/sec */
}I2C_CLOCKRATE;


/**/
typedef struct ChannelConfig_t
{
	I2C_CLOCKRATE	ClockRate; /*There were 2 functions I2C_TurnOn/OffDivideByFive
	ClockinghiSpeedDevice (FTC_HANDLE fthandle) in the old DLL. This function turns on the
	divide by five for the MPSSE clock to allow the hi-speed devices FT2232H and FT4232H to
	clock at the same rate as the FT2232D device. This allows for backward compatibility
	NOTE: This feature is probably a per chip feature and not per device*/

	uint8			LatencyTimer; /*Required value, in milliseconds, of latency timer.
	Valid range is 2 – 255
	In the FT8U232AM and FT8U245AM devices, the receive buffer timeout that is used to flush
	remaining data from the receive buffer was fixed at 16 ms. In all other FTDI devices, this
	timeout is programmable and can be set at 1 ms intervals between 2ms and 255 ms.  This
	allows the device to be better optimized for protocols requiring faster response times from
	short data packets
	NOTE: This feature is probably a per chip feature and not per device*/

	/*Review comment: Maybe we could call use a bitmask called Options as a replacement
	parameter for Enable3PhaseClocking if that is preferable to Flags?  It needs to be quite
	vague and generic since at this point we don't know what it will be setting in the future */

	uint32			Options;	/*This member provides a way to enable/disable features
	specific to the protocol that are implemented in the chip
	BIT0		: 3PhaseDataClocking - Setting this bit will turn on 3 phase data clocking for a
			FT2232H dual hi-speed device or FT4232H quad hi-speed device. Three phase
			data clocking, ensures the data is valid on both edges of a clock
	BIT1		: Loopback
	BIT2		: Clock stretching
	BIT3 -BIT31		: Reserved
	*/
}ChannelConfig;


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
FTDI_API FT_STATUS I2C_InitChannel(FT_HANDLE handle, ChannelConfig *config);
FTDI_API FT_STATUS I2C_CloseChannel(FT_HANDLE handle);
FTDI_API FT_STATUS I2C_DeviceRead(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransfered, uint32 options);
FTDI_API FT_STATUS I2C_DeviceWrite(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransfered, uint32 options);
FTDI_API FT_STATUS I2C_GetDeviceID(FT_HANDLE handle, uint8 deviceAddress,
uint8* deviceID);



/******************************************************************************/


#endif	/*FTDI_I2C_H*/

