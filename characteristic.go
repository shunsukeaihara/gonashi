package gonashi

import "github.com/flemay/gatt"

const (
	KonashiServiceUUID = "229bff0003fb40da98a7b0def65c2d4b"

	// PIO
	KonashiPioSettingUUID              = "229b300003fb40da98a7b0def65c2d4b"
	KonashiPioPullUpUUID               = "229b300103fb40da98a7b0def65c2d4b"
	KonashiPioOutputUUID               = "229b300203fb40da98a7b0def65c2d4b"
	KonashiPioInputNotificationUUID    = "229b300303fb40da98a7b0def65c2d4b"
	KonashiPioInputNotificationReadLen = 1

	// PWM
	KonashiPWmConfigUUID = 0x3004
	KonashiPWmParamUUID  = 0x3005
	KonashiPWmDutyUUID   = 0x3006

	// Analog
	KonashiAnalogDriveUUID = 0x3007
	KonashiAnalogRead0UUID = 0x3008
	KonashiAnalogRead1UUID = 0x3009
	KonashiAnalogRead2UUID = 0x300A
	KonashiAnalogReadLen   = 2

	// I2C
	KonashiI2cConfigUUID    = 0x300B
	KonashiI2cStartStopUUID = 0x300C
	KonashiI2cWriteUUID     = 0x300D
	KonashiI2cReadParamUIUD = 0x300E
	KonashiI2cReadUUID      = 0x300F

	// Uart
	KonashiUartConfigUUID            = 0x3010
	KonashiUartBaudrateUUID          = 0x3011
	KonashiUartTxUUID                = 0x3012
	KonashiUartRXNotificationUUID    = 0x3013
	KonashiUartRXNotificationReadLen = 1

	// Hardware
	KonashiHardwareResetUUID                 = 0x3014
	KonashiHardwareLowBATNotificationUUID    = 0x3015
	KonashiHardwareLowBATNotificationReadLen = 1
)

var (
	KonashiService *gatt.Service = gatt.NewService(gatt.MustParseUUID(KonashiServiceUUID))

	// PIO
	KonashiPioSetting *gatt.Characteristic = gatt.NewCharacteristic(gatt.MustParseUUID(KonashiPioSettingUUID),
		KonashiService, 0, 0, 0)
)
