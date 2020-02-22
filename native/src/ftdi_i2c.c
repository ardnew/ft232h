/*!
 * \file ftdi_i2c.c
 *
 * \author FTDI
 * \date 20110321
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
 * 0.2 - 20110708 - Added 3-phase-clocking functionality in I2C_InitChannel
 * 0.3 - 20111103 - Added I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE
 *					Added I2C_TRANSFER_OPTIONS_FAST_TRANSFER(bit/byte/noAddress modes)
 *					Returns FT_DEVICE_NOT_FOUND if addressed slave doesn't respond
 *					Adjustment to clock rate if 3-phase-clocking is enabled
 * 0.5 - 20140918 - Modified R/W functions to fix glitch issue
 *					I2C_Read8bitsAndGiveAck, I2C_Write8bitsAndGetAck,
 *					I2C_FastRead and I2C_FastWrite
 */


/******************************************************************************/
/*								Include files					  			  */
/******************************************************************************/
#include "ftdi_infra.h"		/*Common portable infrastructure(datatypes, libraries, etc)*/
#include "ftdi_common.h"	/*Common across I2C, SPI, JTAG modules*/
#include "ftdi_i2c.h"		/*I2C specific*/
#include "ftdi_mid.h"		/*Middle layer*/


/******************************************************************************/
/*								Macro and type defines					  		  */
/******************************************************************************/
/* Enabling the macro will lead to checking of all parameters that are passed to the library from
    the user application*/
#define ENABLE_PARAMETER_CHECKING	1

/* This macro will enable all code that is needed to support I2C_GETDEVICEID feature */
/* #define I2C_CMD_GETDEVICEID_SUPPORTED	1 */

/* This macro will enable the code to read acknowledgements from slaves in I2C_RawWrite */
#define FASTWRITE_READ_ACK

#define START_DURATION_1	10
#define START_DURATION_2	20

#define STOP_DURATION_1	10
#define STOP_DURATION_2	10
#define STOP_DURATION_3	10

#define SEND_ACK			0x00
#define SEND_NACK			0x80

#define I2C_ADDRESS_READ_MASK	0x01	/*LSB 1 = Read*/
#define I2C_ADDRESS_WRITE_MASK	0xFE	/*LSB 0 = Write*/

#ifdef I2C_CMD_GETDEVICEID_SUPPORTED
/* This enum lists the supported I2C modes*/
typedef enum I2C_Modes_t{
I2C_STANDARD_MODE = 0,
I2C_FAST_MODE,
I2C_FAST_MODE_PLUS,
I2C_HIGH_SPEED_MODE,
I2C_MAXIMUM_SUPPORTED_MODES
}I2C_Modes;

/* This enum lists the various I2C bus condition*/
typedef enum I2C_Bus_Condition_t{
I2C_CONDITION_PRESTART,
I2C_CONDITION_START,
I2C_CONDITION_POSTSTART,
I2C_CONDITION_PRESTOP,
I2C_CONDITION_STOP,
I2C_CONDITION_POSTSTOP,
I2C_MAXIMUM_SUPPORTED_CONDITIONS
}I2C_Bus_Condition;
#endif

/******************************************************************************/
/*								Local function declarations					  */
/******************************************************************************/
FT_STATUS I2C_Restart(FT_HANDLE handle);
FT_STATUS I2C_Write8bitsAndGetAck(FT_HANDLE handle, uint8 data, bool *ack);
FT_STATUS I2C_Read8bitsAndGiveAck(FT_HANDLE handle, uint8 *data, bool ack);
FT_STATUS I2C_WriteDeviceAddress(FT_HANDLE handle, uint32 deviceAddress,
			bool direction, bool AddLen10Bit, bool *ack);
FT_STATUS I2C_SaveChannelConfig(FT_HANDLE handle, ChannelConfig *config);
FT_STATUS I2C_GetChannelConfig(FT_HANDLE handle, ChannelConfig *config);
FT_STATUS I2C_Start(FT_HANDLE handle);
FT_STATUS I2C_Stop(FT_HANDLE handle);
FT_STATUS I2C_FastWrite(FT_HANDLE handle, uint32 deviceAddress,
uint32 bitsToTransfer, uint8 *buffer, uint8 *ack, uint32 *bytesTransferred,
uint32 options);
FT_STATUS I2C_FastRead(FT_HANDLE handle,uint32 deviceAddress,
uint32 bitsToTransfer, uint8 *buffer, uint8 *ack, uint32 *bytesTransferred,
uint32 options);


/******************************************************************************/
/*								Global variables							  */
/******************************************************************************/

#ifdef I2C_CMD_GETDEVICEID_SUPPORTED
/*!
 * \brief I2C bus condition timings table
 *
 * This table contains the minimim time that the I2C bus needs to be held in, in order to register
 * a condition such as start, stop restart, etc. The table has one row each for the different bus
 * speeds.
 *
 * \sa
 * \note
 * \warning
 */
const uint64 I2C_Timings[I2C_MAXIMUM_SUPPORTED_MODES]
[I2C_MAXIMUM_SUPPORTED_CONDITIONS] = {
/*Durations for conditions are as follows:
pre-start, start, post-start, pre-stop, stop, post-stop */
{0,0,0,0,0,0,}, /* I2C_CLOCK_STANDARD_MODE */
{0,0,0,0,0,0,}, /*I2C_CLOCK_FAST_MODE */
{0,0,0,0,0,0,}, /* I2C_CLOCK_FAST_MODE_PLUS */
{0,0,0,0,0,0}	/* I2C_CLOCK_HIGH_SPEED_MODE */
};
#endif


/******************************************************************************/
/*						Public function definitions						  */
/******************************************************************************/

/*!
 * \brief Gets the number of I2C channels connected to the host
 *
 * This function gets the number of I2C channels that are connected to the host system
 * The number of ports available in each of these chips are different.
 *
 * \param[out] *numChannels Pointer to variable in which the no of channels will be returned
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note This function doesn't return the number of FTDI chips connected to the host system
 * \note FT2232D has 1 MPSSE port
 * \note FT2232H has 2 MPSSE ports
 * \note FT4232H has 4 ports but only 2 of them have MPSSEs
 * so call to this function will return 2 if a FT4232 is connected to it.
 * \warning
 */
FTDI_API FT_STATUS I2C_GetNumChannels(uint32 *numChannels)
{
	FT_STATUS status;

	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(numChannels);
#endif
	status = FT_GetNumChannels(I2C, numChannels);
	CHECK_STATUS(status);
	FN_EXIT;
	return status;
}

/*!
 * \brief Provides information about channel
 *
 * This function takes a channel index (valid values are from 1 to the value returned by
 * I2C_GetNumChannels) and provides information about the channel in the form of a populated
 * ChannelInfo structure.
 *
 * \param[in] index Index of the channel
 * \param[out] chanInfo Pointer to FT_DEVICE_LIST_INFO_NODE structure(see D2XX \
 *			  Programmer's Guide)
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \note  The channel ID can be determined by the user from the last digit of the location ID
 * \sa
 * \warning
 */
FTDI_API FT_STATUS I2C_GetChannelInfo(uint32 index,
					FT_DEVICE_LIST_INFO_NODE *chanInfo)
{
	FT_STATUS status;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(chanInfo);
#endif
	status = FT_GetChannelInfo(I2C,index+1,chanInfo);
	CHECK_STATUS(status);
	FN_EXIT;
	return status;
}

/*!
 * \brief Opens a channel and returns a handle to it
 *
 * This function opens the indexed channel and returns a handle to it
 *
 * \param[in] index Index of the channel
 * \param[out] handle Pointer to the handle of the opened channel
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note Trying to open an already open channel will return an error code
 * \warning
 */
FTDI_API FT_STATUS I2C_OpenChannel(uint32 index, FT_HANDLE *handle)
{
	FT_STATUS status;
	/* Opens a channel and returns the pointer to its handle */
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(handle);
#endif
	status = FT_OpenChannel(I2C,index+1,handle);
	DBG(MSG_DEBUG,"index=%u handle=%u\n",(unsigned)index,(unsigned)*handle);
	CHECK_STATUS(status);
	FN_EXIT;
	return status;
}

/*!
 * \brief Initializes a channel
 *
 * This function initializes the channel and the communication parameters associated with it
 *
 * \param[in] handle Handle of the channel
 * \param[out] config Pointer to ChannelConfig structure(memory to be allocated by caller)
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FTDI_API FT_STATUS I2C_InitChannel(FT_HANDLE handle, ChannelConfig *config)
{
	FT_STATUS status;
	uint8 buffer[3];//3
	uint32 noOfBytesToTransfer;
	uint32 noOfBytesTransferred;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(config);
	CHECK_NULL_RET(handle);
#endif
	if(!(config->Options & I2C_DISABLE_3PHASE_CLOCKING))
	{/* Adjust clock rate if 3phase clocking should be enabled */
		config->ClockRate = (config->ClockRate * 3)/2;
	}
	DBG(MSG_DEBUG,"handle=0x%x ClockRate=%u LatencyTimer=%u Options=0x%x\n",\
		(unsigned)handle,(unsigned)config->ClockRate,	\
		(unsigned)config->LatencyTimer,(unsigned)config->Options);
	status = FT_InitChannel(I2C,handle,(uint32)config->ClockRate,	\
		(uint32)config->LatencyTimer,(uint32)config->Options);
	CHECK_STATUS(status);

	if(!(config->Options & I2C_DISABLE_3PHASE_CLOCKING))
	{
		DBG(MSG_DEBUG,"Enabling 3 phase clocking\n");
		noOfBytesToTransfer = 1;
		noOfBytesTransferred = 0;
		buffer[0] = MPSSE_CMD_ENABLE_3PHASE_CLOCKING;/* MPSSE command */
		status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
			buffer,&noOfBytesTransferred);
		CHECK_STATUS(status);
	}

	/*Save the channel's config data for later use*/
	status = I2C_SaveChannelConfig(handle,config);
	CHECK_STATUS(status);
	FN_EXIT;
	return status;
}

/*!
 * \brief Closes a channel
 *
 * Closes a channel and frees all resources that were used by it
 *
 * \param[in] handle Handle of the channel
 * \param[out] none
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FTDI_API FT_STATUS I2C_CloseChannel(FT_HANDLE handle)
{
	FT_STATUS status;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
		CHECK_NULL_RET(handle);
#endif
	status = FT_CloseChannel(I2C,handle);
	CHECK_STATUS(status);
	FN_EXIT;
	return status;
}

/*!
 * \brief Reads data from I2C slave
 *
 * This function reads the specified number of bytes from an addressed I2C slave
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C slave
 * \param[in] sizeToTransfer Number of bytes to be read
 * \param[out] buffer Pointer to the buffer where data is to be read
 * \param[out] sizeTransferred Pointer to variable containing the number of bytes read
 * \param[in] options This parameter specifies data transfer options. Namely, if a start/stop bits
 *			are required, if the transfer should continue or stop if device nAcks, etc
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa	Definitions of macros I2C_TRANSFER_OPTIONS_START_BIT,
 *		I2C_TRANSFER_OPTIONS_STOP_BIT, I2C_TRANSFER_OPTIONS_BREAK_ON_NACK,
 *		I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES,
 *		I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS &
 *		I2C_TRANSFER_OPTIONS_NO_ADDRESS
 * \note
 * \warning
 */
FTDI_API FT_STATUS I2C_DeviceRead(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransferred, uint32 options)
{
	FT_STATUS status=FT_OK;
	bool ack=FALSE;
	uint32 i;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(handle);
	CHECK_NULL_RET(buffer);
	CHECK_NULL_RET(sizeTransferred);
	if(deviceAddress>127)
	{
		DBG(MSG_WARN,"deviceAddress(0x%x) is greater than 127\n", \
			(unsigned)deviceAddress);
		return FT_INVALID_PARAMETER;
	}
#endif
	DBG(MSG_DEBUG,"handle=0x%x deviceAddress=0x%x sizeToTransfer=%u options= \
	0x%x\n",(unsigned)handle, (unsigned)deviceAddress, (unsigned)sizeToTransfer,
	(unsigned)options);

	LOCK_CHANNEL(handle);
	Mid_PurgeDevice(handle);

	if(options & I2C_TRANSFER_OPTIONS_FAST_TRANSFER)
	{
		status = I2C_FastRead(handle, deviceAddress, sizeToTransfer, buffer, NULL, sizeTransferred, options);
	}
	else
	{
		/* Write START bit */
		if(options & I2C_TRANSFER_OPTIONS_START_BIT)
		{
			status = I2C_Start(handle);
			CHECK_STATUS(status);
		}

		/* Write device address (with LSB=1 => READ)  & Get ACK */
		status = I2C_WriteDeviceAddress(handle,deviceAddress,TRUE,FALSE,&ack);
		CHECK_STATUS(status);

		if(!ack) /*ack bit set actually means device nAcked*/
		{
			/* LOOP until sizeToTransfer */
			for(i=0; ((i<sizeToTransfer) && (status == FT_OK)); i++)
			{
				/* Read byte to buffer & Give ACK 
				(or nACK if it is last byte and I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE is set)*/
				status = I2C_Read8bitsAndGiveAck(handle,&(buffer[i]),			\
				(i<(sizeToTransfer-1))?TRUE:									\
				((options & I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE)?FALSE:TRUE));
			}

			*sizeTransferred = i;
			if(*sizeTransferred != sizeToTransfer)
			{
				DBG(MSG_ERR," sizeToTransfer=%u sizeTransferred=%u\n",\
					(unsigned)sizeToTransfer, (unsigned)*sizeTransferred);
				status = FT_IO_ERROR;
			}
			else
			{
				/* Write STOP bit */
				if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
				{
					status = I2C_Stop(handle);
					CHECK_STATUS(status);
				}
			}
		}
		else
		{
			DBG(MSG_ERR,"I2C device with address 0x%x didn't ack when addressed\n",(unsigned)deviceAddress);
			/* Write STOP bit */
			if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
			{
				status = I2C_Stop(handle);
				CHECK_STATUS(status);
			}
			/*20111102 : FT_IO_ERROR was returned when a device doesn't respond to the
	 		master when it is addressed, as well as when a data transfer fails. To distinguish
	 		between these to errors, FT_DEVICE_NOT_FOUND is now returned after a device
	 		doesn't respond when its addressed*/
			/* old code: status = FT_IO_ERROR; */
			status = FT_DEVICE_NOT_FOUND;
		}
	}
	UNLOCK_CHANNEL(handle);
	FN_EXIT;
	return status;
}

/*!
 * \brief Writes data from I2C slave
 *
 * This function writes the specified number of bytes from an addressed I2C slave
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C slave
 * \param[in] sizeToTransfer Number of bytes to be written
 * \param[out] buffer Pointer to the buffer from where data is to be written
 * \param[out] sizeTransferred Pointer to variable containing the number of bytes written
 * \param[in] options This parameter specifies data transfer options. Namely if a start/stop bits
 *			are required, if the transfer should continue or stop if device nAcks, etc
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa	Definitions of macros I2C_TRANSFER_OPTIONS_START_BIT,
 *		I2C_TRANSFER_OPTIONS_STOP_BIT, I2C_TRANSFER_OPTIONS_BREAK_ON_NACK,
 *		I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES,
 *		I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS &
 *		I2C_TRANSFER_OPTIONS_NO_ADDRESS
 * \note
 * \warning
 */
FTDI_API FT_STATUS I2C_DeviceWrite(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint32 *sizeTransferred, uint32 options)
{
	FT_STATUS status=FT_OK;
	bool ack=FALSE;
	uint32 i;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
	CHECK_NULL_RET(handle);
	CHECK_NULL_RET(buffer);
	CHECK_NULL_RET(sizeTransferred);
	if(deviceAddress>127)
	{
		DBG(MSG_WARN,"deviceAddress(0x%x) is greater than 127\n", \
			(unsigned)deviceAddress);
		return FT_INVALID_PARAMETER;
	}
#endif
	DBG(MSG_DEBUG,"handle=0x%x deviceAddress=0x%x sizeToTransfer=%u options= \
	0x%x\n",(unsigned)handle, (unsigned)deviceAddress, (unsigned)sizeToTransfer,
	(unsigned)options);

	LOCK_CHANNEL(handle);
	Mid_PurgeDevice(handle);

	if(options & I2C_TRANSFER_OPTIONS_FAST_TRANSFER)
	{
		status = I2C_FastWrite(handle, deviceAddress, sizeToTransfer, buffer, NULL, sizeTransferred, options);
	}
	else
	{
		/* Write START bit */
		if(options & I2C_TRANSFER_OPTIONS_START_BIT)
		{
			status = I2C_Start(handle);
			CHECK_STATUS(status);
		}

		/* Write device address (with LSB=0 => WRITE) & Get ACK*/
		status = I2C_WriteDeviceAddress(handle,deviceAddress,FALSE,FALSE,&ack);
		CHECK_STATUS(status);

		if(!ack) /*ack bit set actually means device nAcked*/
		{
			/* LOOP until sizeToTransfer */
			for(i=0; ((i<sizeToTransfer) && (status == FT_OK)); i++)
			{
				/* Write byte to buffer & Get ACK */
				ack=0;
				status = I2C_Write8bitsAndGetAck(handle,buffer[i],&ack);
				DBG(MSG_DEBUG,"handle=0x%x buffer[%u]=0x%x ack=0x%x \n", \
					(unsigned)handle, (unsigned)i, (unsigned)buffer[i],
					(unsigned)(1&ack));
				if(ack)
				{
					DBG(MSG_WARN,"I2C device(address 0x%x) nAcked while writing\
	byte no %d(i.e. 0x%x\n",(unsigned)deviceAddress,(int)i,(unsigned)buffer[i]);
					/* add bit in options to return with error if device nAcked
					sizeTransferred = number of correctly transfered bytes */
					if(options & I2C_TRANSFER_OPTIONS_BREAK_ON_NACK)
					{
						/*status = FT_FAILED_TO_WRITE_DEVICE;
						break;*/
						DBG(MSG_WARN,"returning FT_FAILED_TO_WRITE_DEVICE \
						options=0x%x ack=0x%x\n",options,ack);
						
						/* Write STOP bit */
						if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
						{
							status = I2C_Stop(handle);
							CHECK_STATUS(status);
						}
						return FT_FAILED_TO_WRITE_DEVICE;
					}
				}
			}
			*sizeTransferred = i;
			if(*sizeTransferred != sizeToTransfer)
			{
				DBG(MSG_ERR," sizeToTransfer=%u sizeTransferred=%u\n",\
					(unsigned)sizeToTransfer, (unsigned)*sizeTransferred);
				status = FT_IO_ERROR;
			}
			else
			{
				/* Write STOP bit */
				if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
				{
					status = I2C_Stop(handle);
					CHECK_STATUS(status);
				}
			}
		}
		else
		{
			DBG(MSG_ERR,"I2C device with address 0x%x didn't ack when addressed\n",(unsigned)deviceAddress);
			/* Write STOP bit */
			if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
			{
				status = I2C_Stop(handle);
				CHECK_STATUS(status);
			}
			/*20111102 : FT_IO_ERROR was returned when a device doesn't respond to the
	 		master when it is addressed, as well as when a data transfer fails. To distinguish
	 		between these to errors, FT_DEVICE_NOT_FOUND is now returned after a device
	 		doesn't respond when its addressed*/
			/* old code: status = FT_IO_ERROR; */
			status = FT_DEVICE_NOT_FOUND;
		}
	}
	UNLOCK_CHANNEL(handle);
	FN_EXIT;
	return status;
}


#ifdef I2C_CMD_GETDEVICEID_SUPPORTED
/*!
 * \brief Get the I2C device ID
 *
 * This function retrieves the I2C device ID
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C slave
 * \param[out] deviceID Address of memory where the 3byte I2C device ID will be stored
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FTDI_API FT_STATUS I2C_GetDeviceID(FT_HANDLE handle, uint8 deviceAddress,
	uint8* deviceID)
{
	FT_STATUS status=FT_OTHER_ERROR;
	bool ack;
	FN_ENTER;
#ifdef ENABLE_PARAMETER_CHECKING
		CHECK_NULL_RET(handle);
		CHECK_NULL_RET(deviceID);
		if(deviceAddress>127)
		{
			DBG(MSG_WARN,"deviceAddress(0x%x) is greater than 127\n", \
				(unsigned)deviceAddress);
			status = FT_INVALID_PARAMETER;
			return status;
		}
#endif

	status = I2C_Start(handle);
	CHECK_STATUS(status);
	status = I2C_Write8bitsAndGetAck(handle,(uint8)I2C_CMD_GETDEVICEID_RD,&ack);
	CHECK_STATUS(status);
	status = I2C_Write8bitsAndGetAck(handle,deviceAddress, &ack);
	CHECK_STATUS(status);
	status = I2C_Restart(handle);
	CHECK_STATUS(status);
	status = I2C_Write8bitsAndGetAck(handle,(uint8)I2C_CMD_GETDEVICEID_WR,&ack);
	CHECK_STATUS(status);
	status = I2C_Read8bitsAndGiveAck(handle,&(deviceID[0]),I2C_GIVE_ACK);
	CHECK_STATUS(status);
	status = I2C_Read8bitsAndGiveAck(handle,&(deviceID[1]),I2C_GIVE_ACK);
	CHECK_STATUS(status);
	/*NACK 3rd byte*/
	status = I2C_Read8bitsAndGiveAck(handle,&(deviceID[2]),I2C_GIVE_NACK);
	CHECK_STATUS(status);

	FN_EXIT;
	return status;
}
#endif


/******************************************************************************/
/*						Local function definations						  */
/******************************************************************************/
#ifdef I2C_CMD_GETDEVICEID_SUPPORTED
/*!
 * \brief Generate I2C bus restart condition
 *
 * This function generates the restart condition in the I2C bus
 *
 * \param[in] handle Handle of the channel
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_Restart(FT_HANDLE handle)
{
	FT_STATUS status;
	uint8 buffer[3];
	uint32 noOfBytesToTransfer;
	uint32 noOfBytesTransferred;
	I2C_Modes mode;

	FN_ENTER;
#if 0
	status = I2C_GetChannelConfig(handle,config);
	if(FT_OK != status)
		Infra_DbgPrintStatus(status);
	CHECK_NULL_RET(config);
	switch(config->ClockRate)
	{
		case I2C_CLOCK_STANDARD_MODE:
			mode = I2C_STANDARD_MODE;
		break;

		case I2C_CLOCK_FAST_MODE:
			mode = I2C_FAST_MODE;
		break;

		case I2C_CLOCK_FAST_MODE_PLUS:
			mode = I2C_FAST_MODE_PLUS;
		break;

		case I2C_CLOCK_HIGH_SPEED_MODE:
			mode = I2C_HIGH_SPEED_MODE;
		break;

		default:
			mode = I2C_MAXIMUM_SUPPORTED_MODES;
	}
#else
	mode = I2C_STANDARD_MODE;
#endif

	/*Restart condition SDA, SCL: 0,1->1,1->1,0->1,1->0,1*/

	/*I2C_CONDITION_RESTART_1*/
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLHIGH_SDALOW; /* Value */
	buffer[2] = DIRECTION_SCLOUT_SDAOUT; /* Direction */
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_PRESTART]);

	/*I2C_CONDITION_RESTART_2*/
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLHIGH_SDALOW;//  _SDAHIGH; /* Value */										// should NOT be driving SDA high
	buffer[2] = DIRECTION_SCLOUT_SDAIN;//  _SDAOUT; /* Direction */									// Make this input instead to let line be pulled up
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_PRESTART]);

	/*I2C_CONDITION_RESTART_3*/
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLLOW_SDALOW;// _SDAHIGH; /* Value */										// should NOT be driving SDA high
	buffer[2] = DIRECTION_SCLOUT_SDAIN;//  _SDAOUT; /* Direction */									// Make this input instead to let line be pulled up
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_PRESTART]);

	/*I2C_CONDITION_RESTART_4*/
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLHIGH_SDALOW;// _SDAHIGH; /* Value */									// should NOT be driving SDA high
	buffer[2] = DIRECTION_SCLOUT_SDAIN;// _SDAOUT; /* Direction */								// Make this input instead to let line be pulled up
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_PRESTART]);

	/*I2C_CONDITION_RESTART_5*/
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLLOW_SDALOW;// _SDAHIGH; /* Value */									// should NOT be driving SDA high
	buffer[2] = DIRECTION_SCLOUT_SDAIN;//  _SDAOUT; /* Direction */								// Make this input instead to let line be pulled up
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_PRESTART]);

	/*I2C_CONDITION_RESTART -tristate SCL & SDA */
	noOfBytesToTransfer = 3;
	noOfBytesTransferred = 0;
	buffer[0] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;/* MPSSE command */
	buffer[1] = VALUE_SCLLOW_SDALOW; /* Value(0x00=SCL low, SDA low) */
	buffer[2] = DIRECTION_SCLIN_SDAIN; /* Direction */
	status = FT_Channel_Write(I2C,handle,noOfBytesToTransfer,
		buffer,&noOfBytesTransferred);
	if( (FT_OK != status) && (noOfBytesToTransfer != noOfBytesTransferred) )
	{
		Infra_DbgPrintStatus(status);
		DBG(MSG_ERR,"noOfBytesToTransfer=%d noOfBytesTransferred=%d\n",
			(int)noOfBytesToTransfer,(int)noOfBytesTransferred);
	}
	Infra_Delay(I2C_Timings[mode][I2C_CONDITION_POSTSTOP]);

	FN_EXIT;
	return status;
}
#endif



/*!
 * \brief Writes 8 bits and gets the ack bit
 *
 * This function writes 8 bits of data to the I2C bus and gets the ack bit from the device
 *
 * \param[in] handle Handle of the channel
 * \param[in] data The 8bits of data that are to be written to the I2C bus
 * \param[out] ack The acknowledgment bit returned by the I2C device
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_Write8bitsAndGetAck(FT_HANDLE handle, uint8 data, bool *ack)
{
	FT_STATUS status = FT_OTHER_ERROR;
	uint8 buffer[20] = {0};
	uint8 inBuffer[3] = {0};
	uint32 noOfBytes = 0;
	uint32 noOfBytesTransferred = 0;


	FN_ENTER;
	DBG(MSG_DEBUG,"----------Writing byte 0x%x \n",data);

	/* Set direction */
	buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
	buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAOUT;

	/* Command to write 8 bits */
	buffer[noOfBytes++]= MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
	buffer[noOfBytes++]= DATA_SIZE_8BITS;
	buffer[noOfBytes++] = data;

	/* Set SDA to input mode before reading ACK bit */
	buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
	buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAIN;

	/* Command to get ACK bit */
	buffer[noOfBytes++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
	buffer[noOfBytes++] = DATA_SIZE_1BIT;

	/* Command MPSSE to send data to PC immediately */
	buffer[noOfBytes++] = MPSSE_CMD_SEND_IMMEDIATE;

	status = FT_Channel_Write(I2C, handle, noOfBytes, buffer, &noOfBytesTransferred);
	if(FT_OK != status)
	{
		DBG(MSG_DEBUG, "FT_OK != status \n");
		Infra_DbgPrintStatus(status);
	}
	else if(noOfBytes != noOfBytesTransferred)
	{
		DBG(MSG_ERR, "Requested to send %u bytes, no. of bytes sent is %u \
				bytes",(unsigned)noOfBytes,(unsigned)noOfBytesTransferred);
		status = FT_IO_ERROR;
	}
	else
	{
		noOfBytes=1;
		noOfBytesTransferred=0;
#if 0
		{
			uint32 noOfBytesInQueue = 0;
			do
			{
				status = Mid_GetQueueStatus(handle, &noOfBytesInQueue);
			} while (noOfBytesInQueue < noOfBytes && status == FT_OK);
		}
#else
		INFRA_SLEEP(1);
#endif
		status = FT_Channel_Read(I2C, handle, noOfBytes, inBuffer, &noOfBytesTransferred);
		if(FT_OK != status)
		{
			Infra_DbgPrintStatus(status);
		}
		else if(noOfBytes != noOfBytesTransferred)
		{
			DBG(MSG_ERR, "Requested to send %u bytes, no. of bytes sent is %u \
				bytes",(unsigned)noOfBytes,(unsigned)noOfBytesTransferred);
			status = FT_IO_ERROR;
		}
		else
		{
			*ack = (bool)(inBuffer[0] & 0x01);
			DBG(MSG_DEBUG,"	*ack = 0x%x\n", (unsigned)*ack);
		}
	}

	FN_EXIT;
	return status;
}

/*!
 * \brief Reads 8 bits of data and sends ack bit
 *
 * This function reads 8 bits of data from the I2C bus and then writes an ack bit to the bus
 *
 * \param[in] handle Handle of the channel
 * \param[in] *data Pointer to the buffer where the 8bits would be read to
 * \param[in] ack Gives ack to device if set, otherwise gives nAck
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_Read8bitsAndGiveAck(FT_HANDLE handle, uint8 *data, bool ack)
{
	FT_STATUS status = FT_OTHER_ERROR;
	uint8 buffer[20] = {0};
	uint8 inBuffer[3] = {0};
	uint32 noOfBytes = 0;
	uint32 noOfBytesTransferred = 0;


	FN_ENTER;


	/* Set pin directions - SCL is output driven low, SDA is input (set high but does not matter) */
	buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
	buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAIN;

	/* Command to read 8 bits */
	buffer[noOfBytes++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
	buffer[noOfBytes++] = DATA_SIZE_8BITS;

	/* Set directions to make SDA drive out. Pre-set state of SDA first though to avoid glitch */
	if (ack)
	{
		/* We will drive the ACK bit to a '0' so pre-set pin to a '0' */
		buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
		buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAOUT;

		/* Clock out the ack bit as a '0' on negative edge */
		buffer[noOfBytes++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
		buffer[noOfBytes++] = DATA_SIZE_1BIT;
		buffer[noOfBytes++] = SEND_ACK;
	}
	else
	{
		/* We will release the ACK bit to a '1' so pre-set pin to a '1' by making it an input */
		buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
		buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAIN;

		/* Clock out the ack bit as a '1' on negative edge - never actually seen on line since SDA is input but burns off one bit time */
		buffer[noOfBytes++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
		buffer[noOfBytes++] = DATA_SIZE_1BIT;
		buffer[noOfBytes++] = SEND_NACK;
	}

	/* Back to Idle */
	buffer[noOfBytes++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[noOfBytes++] = VALUE_SCLLOW_SDALOW;
	buffer[noOfBytes++] = DIRECTION_SCLOUT_SDAIN;

	/* Command MPSSE to send data to PC immediately */
	buffer[noOfBytes++] = MPSSE_CMD_SEND_IMMEDIATE;

	status = FT_Channel_Write(I2C, handle, noOfBytes, buffer, &noOfBytesTransferred);
	if(FT_OK != status)
	{
		DBG(MSG_DEBUG, "FT_OK != status \n");
		Infra_DbgPrintStatus(status);
	}
	else if(noOfBytes != noOfBytesTransferred)
	{
		DBG(MSG_ERR, "Requested to send %u bytes, no. of bytes sent is %u \
				bytes",(unsigned)noOfBytes,(unsigned)noOfBytesTransferred);
		status = FT_IO_ERROR;
	}
	else
	{
		noOfBytes=1;
		noOfBytesTransferred=0;
#if 0
		{
			uint32 noOfBytesInQueue = 0;
			do
			{
				status = Mid_GetQueueStatus(handle, &noOfBytesInQueue);
			} while (noOfBytesInQueue < noOfBytes && status == FT_OK);
		}
#else
		INFRA_SLEEP(1);
#endif
		status = FT_Channel_Read(I2C, handle, noOfBytes, inBuffer, &noOfBytesTransferred);
		if(FT_OK != status)
		{
			Infra_DbgPrintStatus(status);
		}
		else if(noOfBytes != noOfBytesTransferred)
		{
			DBG(MSG_ERR, "Requested to read %u bytes, no. of bytes read is %u \
					bytes",(unsigned)noOfBytes,(unsigned)noOfBytesTransferred);
			status = FT_IO_ERROR;
		}
		else
		{
			*data = inBuffer[0];
			DBG(MSG_DEBUG,"	*data = 0x%x\n", (unsigned)*data);
		}
	}

	FN_EXIT;
	return status;
}



/*!
 * \brief This function generates the START, ADDRESS, DATA(write) & STOP phases in the I2C
 *		bus without having delays between these phases
 *
 * This function allocates memory locally, makes MPSSE command frames to write each data
 * byte/bit, makes MPSSE command frames to read the acknowledgement bits, and then writes
 * all these to the MPSSE in one shot. This function is useful where delays between START, DATA
 * and STOP phases are not prefered.
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C Slave. This parameter is ignored if flag
 *			I2C_TRANSFER_OPTIONS_NO_ADDRESS is set in the options parameter
 * \param[in] sizeToTransfer Number of bytes or bits to be written, depending on if
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS or
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES is set in options parameter
 * \param[in] *buffer Pointer to the buffer from where data is to be written. The user application
 * 			is expected to send a byte array of length sizeToTransfer. However if
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS is set then the length of the byte
 *			array should be sizeToTransfer/8.
 * \param[out] *ack Pointer to the buffer to where the ack bits are to be stored. The ack bits are
 *			 ignored if NULL is passed
 * \param[out] sizeTransferred Pointer to variable containing the number of bytes/bits written
 * \param[in] options This parameter specifies data transfer options. Namely if a start/stop
 *			conditions are required, if size is in bytes or bits, if address is provided.
 * \return 	Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa		Definitions of I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES,
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS,
 *			I2C_TRANSFER_OPTIONS_NO_ADDRESS,
 *			I2C_TRANSFER_OPTIONS_START_BIT,
 *			I2C_TRANSFER_OPTIONS_STOP_BIT
 * \note The I2C_TRANSFER_OPTIONS_BREAK_ON_NACK bit in the options parameter is not
 *          applicable for this function.
 * \warning
 */
FT_STATUS I2C_FastWrite(FT_HANDLE handle, uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint8 *ack, uint32 *sizeTransferred, uint32 options)
{
	FT_STATUS status=FT_OK;
	uint32 i=0; /* index of cmdBuffer that is filled */
	uint32 j=0; /* scratch register */
	uint32 sizeTotal;
	uint8* outBuffer;
	uint32 bytesRead;
	uint32 bytesToTransfer;
	uint8 tempAddress;
	uint8* inBuffer;


	FN_ENTER;


	if(!(options & I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES))
	{
		// Fast read supports I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES only, not I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS
		return FT_INVALID_PARAMETER;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Allocate COMMAND buffer #############
	
	bytesToTransfer = sizeToTransfer;
	
	/* Calculate size of required buffer */
	sizeTotal = 0 /* for send immediate */
	+ (bytesToTransfer*(6+5)) /* the size of data itself */
	+ ((!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))?(11):(0)) /* for address byte */
	+ ((options & I2C_TRANSFER_OPTIONS_START_BIT)?((START_DURATION_1+	\
			START_DURATION_2+1)*3):0) /* size required for START */
	+ ((options & I2C_TRANSFER_OPTIONS_STOP_BIT)?((STOP_DURATION_1+		\
			STOP_DURATION_2+STOP_DURATION_3+1)*3):0); /* size for STOP */

	/* Allocate buffers */
	outBuffer = (uint8*) INFRA_MALLOC(sizeTotal);
	if(NULL == outBuffer)
	{
		return FT_INSUFFICIENT_RESOURCES;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write START command #############
	if(options & I2C_TRANSFER_OPTIONS_START_BIT)
	{
		DBG(MSG_DEBUG,"adding START condition\n");
		/* SCL high, SDA high */
		for(j=0;j<START_DURATION_1;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;
		}
		/* SCL high, SDA low */
		for(j=0;j<START_DURATION_2;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/*SCL low, SDA low */
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write ADDRESS #############
	if(!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))
	{
		tempAddress = (uint8)deviceAddress;
		tempAddress = (tempAddress << 1);
		tempAddress = (tempAddress & I2C_ADDRESS_WRITE_MASK);
		DBG(MSG_DEBUG,"7bit I2C address plus direction bit = 0x%x\n", tempAddress);

		/*set direction*/
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;

		/* write address + direction bit */
		outBuffer[i++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
		outBuffer[i++] = DATA_SIZE_8BITS;
		outBuffer[i++] = tempAddress;

		/* Set SDA to input mode before reading ACK bit */
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;

		/* Command to get ACK bit */
		outBuffer[i++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
		outBuffer[i++] = DATA_SIZE_1BIT;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write ACTUAL DATA #############
	for (j=0; j<bytesToTransfer; j++)
	{
		/*set direction*/
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;

		/* Command to write 8 bits */
		outBuffer[i++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
		outBuffer[i++] = DATA_SIZE_8BITS;
		outBuffer[i++] = buffer[j];

		/* Set SDA to input mode before reading ACK bit */
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;

		/* Command to get ACK bit */
		outBuffer[i++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
		outBuffer[i++] = DATA_SIZE_1BIT;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write STOP command #############
	if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
	{
		/* SCL low, SDA low */
		for(j=0;j<STOP_DURATION_1;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLLOW_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/* SCL high, SDA low */
		for(j=0;j<STOP_DURATION_2;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/* SCL high, SDA high */
		for(j=0;j<STOP_DURATION_3;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;
		}
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
		outBuffer[i++] = DIRECTION_SCLIN_SDAIN; /* Tristate the SCL & SDA pins */
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write commands #############
	DBG(MSG_DEBUG,"i=%u bytesToTransfer=%u\n",(unsigned)i,(unsigned)bytesToTransfer);
	status = FT_Channel_Write(I2C,handle,i,outBuffer,&bytesRead);
	*sizeTransferred = sizeToTransfer;
	INFRA_FREE(outBuffer);
	CHECK_STATUS(status);


	//--------------------------------------------------------------------------------------------
	// ############## Read ACKS #############
	/* Read ack of address */
	if(!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))
	{
		uint8 addAck;
		status = FT_Channel_Read(I2C, handle, 1, &addAck, &bytesRead);
		CHECK_STATUS(status);
	}

	/* Read 1bit ack after each 8bits written */
	inBuffer = INFRA_MALLOC(bytesToTransfer);
	if(NULL != inBuffer)
	{
		status = FT_Channel_Read(I2C, handle, bytesToTransfer, inBuffer, &bytesRead);
		if (status == FT_OK && ack)
		{/* Copy the ack bits into the ack buffer if provided */
			INFRA_MEMCPY(ack,inBuffer,bytesRead);
		}
		INFRA_FREE(inBuffer);
		CHECK_STATUS(status);
	}


	FN_EXIT;
	return status;
}


/*!
 * \brief This function generates the START, ADDRESS, DATA(read) & STOP phases in the I2C
 *		bus without having delays between these phases
 *
 * This function allocates memory locally, makes MPSSE command frames to read each data
 * byte/bit, makes MPSSE command frames to write the acknowledgement bits, and then writes
 * all these to the MPSSE in one shot. This function is useful where delays between START, DATA
 * and STOP phases are not prefered.
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C Slave. This parameter is ignored if flag
 *			I2C_TRANSFER_OPTIONS_NO_ADDRESS is set in the options parameter
 * \param[in] sizeToTransfer Number of bytes or bits to be written, depending on if
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS or
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES is set in options parameter
 * \param[in] *buffer Pointer to the buffer to where data is to be read. The user application
 * 			is expected to send a byte array of length sizeToTransfer. However if
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS is set then the length of the byte
 *			array should be sizeToTransfer/8
 * \param[out] *ack Reserved (place holder for buffer in which user can provide ack/nAck bits)
 * \param[out] sizeTransferred Pointer to variable containing the number of bytes/bits written
 * \param[in] options This parameter specifies data transfer options. Namely if a start/stop
 *			conditions are required, if size is in bytes or bits, if address is provided.
 * \return 	Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa		Definitions of I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES,
 *			I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS,
 *			I2C_TRANSFER_OPTIONS_NO_ADDRESS,
 *			I2C_TRANSFER_OPTIONS_START_BIT,
 *			I2C_TRANSFER_OPTIONS_STOP_BIT
 * \note The I2C_TRANSFER_OPTIONS_NACK_LAST_BYTE bit in the options parameter is not
 *          applicable for this function.
 * \warning
 */
FT_STATUS I2C_FastRead(FT_HANDLE handle,uint32 deviceAddress,
uint32 sizeToTransfer, uint8 *buffer, uint8 *ack, uint32 *sizeTransferred, uint32 options)
{
	FT_STATUS status=FT_OK;
	uint32 i=0; /* index of cmdBuffer that is filled */
	uint32 j=0; /* scratch register */
	uint32 sizeTotal;
	uint8* outBuffer;
	uint32 bytesRead;
	uint32 bytesToTransfer;
	uint8 tempAddress;


	FN_ENTER;


	if(!(options & I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES))
	{
		// Fast read supports I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BYTES only, not I2C_TRANSFER_OPTIONS_FAST_TRANSFER_BITS
		return FT_INVALID_PARAMETER;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Allocate COMMAND buffer #############
	
	bytesToTransfer = sizeToTransfer;
	
	/* Calculate size of required buffer */
	sizeTotal = 0 /* for send immediate */
	+ (bytesToTransfer*(8+6)) /* the size of data itself */
	+ ((!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))?(11):(0)) /* for address byte */
	+ ((options & I2C_TRANSFER_OPTIONS_START_BIT)?((START_DURATION_1+	\
			START_DURATION_2+1)*3):0) /* size required for START */
	+ ((options & I2C_TRANSFER_OPTIONS_STOP_BIT)?((STOP_DURATION_1+		\
			STOP_DURATION_2+STOP_DURATION_3+1)*3):0); /* size for STOP */

	/* Allocate buffers */
	outBuffer = (uint8*) INFRA_MALLOC(sizeTotal);
	if(NULL == outBuffer)
	{
		return FT_INSUFFICIENT_RESOURCES;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write START command #############
	if(options & I2C_TRANSFER_OPTIONS_START_BIT)
	{
		DBG(MSG_DEBUG,"adding START condition\n");
		/* SCL high, SDA high */
		for(j=0;j<START_DURATION_1;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;
		}
		/* SCL high, SDA low */
		for(j=0;j<START_DURATION_2;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/*SCL low, SDA low */
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write ADDRESS #############
	if(!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))
	{
		tempAddress = (uint8)deviceAddress;
		tempAddress = (tempAddress << 1);
		tempAddress = (tempAddress | I2C_ADDRESS_READ_MASK);
		DBG(MSG_DEBUG,"7bit I2C address plus direction bit = 0x%x\n", tempAddress);

		/*set direction*/
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;

		/* write address + direction bit */
		outBuffer[i++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
		outBuffer[i++] = DATA_SIZE_8BITS;
		outBuffer[i++] = tempAddress;

		/* Set SDA to input mode before reading ACK bit */
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;

		/* Command to get ACK bit */
		outBuffer[i++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
		outBuffer[i++] = DATA_SIZE_1BIT;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Read ACTUAL DATA #############	
	for (j=0; j<bytesToTransfer; j++)
	{
		/*set direction*/
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;

		/*Command to read 8 bits*/
		outBuffer[i++] = MPSSE_CMD_DATA_IN_BITS_POS_EDGE;
		outBuffer[i++] = DATA_SIZE_8BITS;

		/*Command MPSSE to send data to PC immediately */
		if (j < bytesToTransfer - 1)// then ack
		{
			// We will drive the ACK bit to a '0' so pre-set pin to a '0'
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLLOW_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;

			// Clock out the ack bit as a '0' on negative edge
			outBuffer[i++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
			outBuffer[i++] = DATA_SIZE_1BIT;
			outBuffer[i++] = SEND_ACK;
		}
		else// otherwise nack last byte
		{
			// We will release the ACK bit to a '1' so pre-set pin to a '1' by making it an input
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLLOW_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;

			// Clock out the ack bit as a '1' on negative edge - never actually seen on line since SDA is input but burns off one bit time
			outBuffer[i++] = MPSSE_CMD_DATA_OUT_BITS_NEG_EDGE;
			outBuffer[i++] = DATA_SIZE_1BIT;
			outBuffer[i++] = SEND_NACK;
		}

		/*Back to Idle*/
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLLOW_SDALOW;
		outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write STOP command #############
	if(options & I2C_TRANSFER_OPTIONS_STOP_BIT)
	{
		/* SCL low, SDA low */
		for(j=0;j<STOP_DURATION_1;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLLOW_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/* SCL high, SDA low */
		for(j=0;j<STOP_DURATION_2;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDALOW;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAOUT;
		}
		/* SCL high, SDA high */
		for(j=0;j<STOP_DURATION_3;j++)
		{
			outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
			outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
			outBuffer[i++] = DIRECTION_SCLOUT_SDAIN;
		}
		outBuffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		outBuffer[i++] = VALUE_SCLHIGH_SDAHIGH;
		outBuffer[i++] = DIRECTION_SCLIN_SDAIN; /* Tristate the SCL & SDA pins */
	}


	//--------------------------------------------------------------------------------------------
	// ############## Write commands #############
	DBG(MSG_DEBUG,"i=%u bytesToTransfer=%u\n",(unsigned)i,(unsigned)bytesToTransfer);
	status = FT_Channel_Write(I2C,handle,i,outBuffer,&bytesRead);
	*sizeTransferred = sizeToTransfer;
	INFRA_FREE(outBuffer);
	CHECK_STATUS(status);


	//--------------------------------------------------------------------------------------------
	// ############## Read ACKS #############
	/* Read ack of address */
	if(!(options & I2C_TRANSFER_OPTIONS_NO_ADDRESS))
	{
		uint8 addAck;
		status = FT_Channel_Read(I2C, handle, 1, &addAck, &bytesRead);
		CHECK_STATUS(status);
	}


	//--------------------------------------------------------------------------------------------
	// ############## Read ACTUAL DATA #############
	/* Read the actual data from the MPSSE-chip into the host system */
	status = FT_Channel_Read(I2C, handle, bytesToTransfer, buffer, &bytesRead);
	CHECK_STATUS(status);


	FN_EXIT;
	return status;
}



/*!
 * \brief Write I2C device address
 *
 * This function writes the direction and address bits to the I2C bus, and then gets the ACK
 *
 * \param[in] handle Handle of the channel
 * \param[in] deviceAddress Address of the I2C device
 * \param[in] direction 0=Write; 1=Read
 * \param[in] AddLen10Bit Setting this bit specifies 10bit addressing, otherwise 7bit
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_WriteDeviceAddress(FT_HANDLE handle, uint32 deviceAddress,
			bool direction, bool AddLen10Bit, bool *ack)
{
	FT_STATUS status=FT_OTHER_ERROR;
	uint8 tempAddress;
	FN_ENTER;
	if(!AddLen10Bit)
	{/* 7bit addressing */
		tempAddress = (uint8)deviceAddress;
		DBG(MSG_DEBUG,"7bit I2C address = 0x%x\n",tempAddress);
		tempAddress = (tempAddress << 1);
		if(direction)
			tempAddress = (tempAddress | I2C_ADDRESS_READ_MASK);
		else
			tempAddress = (tempAddress & I2C_ADDRESS_WRITE_MASK);
		DBG(MSG_DEBUG,"7bit I2C address plus direction bit = 0x%x\n",\
			tempAddress);
		status = I2C_Write8bitsAndGetAck(handle,tempAddress,ack);
		if(FT_OK != status)
			Infra_DbgPrintStatus(status);
		if((*ack))
		{
			DBG(MSG_ERR,"Didn't receieve an ACK from the addressed device\n");
		}
	}
	else
	{/* 10bit addressing */
		/* TODO: Add support for 10bit addressing */
		DBG(MSG_ERR, "10 bit addressing yet to be supported\n");
		status = FT_NOT_SUPPORTED;
	}
	FN_EXIT;
	return status;
}

/*!
 * \brief Saves the channel's configuration data
 *
 * This function saves the channel's configuration data
 *
 * \param[in] handle Handle of the channel
 * \param[in] config Pointer to ChannelConfig structure(memory to be allocated by caller)
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_SaveChannelConfig(FT_HANDLE handle, ChannelConfig *config)
{
	FT_STATUS status=FT_OTHER_ERROR;
	FN_ENTER;

	status = FT_OK;
	FN_EXIT;
	return status;
}

/*!
 * \brief Retrieves channel's configuration data
 *
 * This function retrieves the channel's configuration data that was previously saved
 *
 * \param[in] handle Handle of the channel
 * \param[in] config Pointer to ChannelConfig structure(memory to be allocated by caller)
 * \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
 * \sa
 * \note
 * \warning
 */
FT_STATUS I2C_GetChannelConfig(FT_HANDLE handle, ChannelConfig *config)
{
	FT_STATUS status=FT_OTHER_ERROR;
	FN_ENTER;

	status = FT_OK;
	FN_EXIT;
	return status;
}



/*!
* \brief Generates the I2C Start condition
*
* This function generates the I2C start condition on the bus.
*
* \param[in] handle Handle of the channel
* \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
* \sa
* \note
* \warning
*/
FT_STATUS I2C_Start(FT_HANDLE handle)
{
	FT_STATUS status;
	uint8 buffer[(START_DURATION_1 + START_DURATION_2 + 1) * 3];
	uint32 i = 0, j = 0;
	uint32 noOfBytesTransferred;
	FN_ENTER;

	/* SCL high, SDA high */
	for (j = 0; j<START_DURATION_1; j++)
	{
		buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[i++] = VALUE_SCLHIGH_SDAHIGH;
		buffer[i++] = DIRECTION_SCLOUT_SDAIN; // Make this input instead to let line be pulled up
	}
	/* SCL high, SDA low */
	for (j = 0; j<START_DURATION_2; j++)
	{
		buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[i++] = VALUE_SCLHIGH_SDALOW;
		buffer[i++] = DIRECTION_SCLOUT_SDAOUT;
	}
	/*SCL low, SDA low */
	buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[i++] = VALUE_SCLLOW_SDALOW;
	buffer[i++] = DIRECTION_SCLOUT_SDAOUT;

	status = FT_Channel_Write(I2C, handle, i, buffer, &noOfBytesTransferred);

	FN_EXIT;
	return status;
}

/*!
* \brief Generates the I2C Stop condition
*
* This function generates the I2C stop condition on the bus.
*
* \param[in] handle Handle of the channel
* \return Returns status code of type FT_STATUS(see D2XX Programmer's Guide)
* \sa
* \note
* \warning
*/
FT_STATUS I2C_Stop(FT_HANDLE handle)
{
	FT_STATUS status;
	uint8 buffer[(STOP_DURATION_1 + STOP_DURATION_2 + STOP_DURATION_3 + 1) * 3];
	uint32 i = 0, j = 0;
	uint32 noOfBytesTransferred;

	FN_ENTER;
	/* SCL low, SDA low */
	for (j = 0; j<STOP_DURATION_1; j++)
	{
		buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[i++] = VALUE_SCLLOW_SDALOW;
		buffer[i++] = DIRECTION_SCLOUT_SDAOUT;
	}
	/* SCL high, SDA low */
	for (j = 0; j<STOP_DURATION_2; j++)
	{
		buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[i++] = VALUE_SCLHIGH_SDALOW;
		buffer[i++] = DIRECTION_SCLOUT_SDAOUT;
	}
	/* SCL high, SDA high */
	for (j = 0; j<STOP_DURATION_3; j++)
	{
		buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
		buffer[i++] = VALUE_SCLHIGH_SDAHIGH;
		buffer[i++] = DIRECTION_SCLOUT_SDAIN; // Make this input instead to let line be pulled up
	}
	buffer[i++] = MPSSE_CMD_SET_DATA_BITS_LOWBYTE;
	buffer[i++] = VALUE_SCLHIGH_SDAHIGH;
	buffer[i++] = DIRECTION_SCLIN_SDAIN; /* Tristate the SCL & SDA pins */

	status = FT_Channel_Write(I2C, handle, i, buffer, &noOfBytesTransferred);

	FN_EXIT;
	return status;
}

