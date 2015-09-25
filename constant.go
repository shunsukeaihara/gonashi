package gonashi

type DioPin int

const (
	KonashiDigitalIO0 IOPin = iota
	KonashiDigitalIO1
	KonashiDigitalIO2
	KonashiDigitalIO3
	KonashiDigitalIO4
	KonashiDigitalIO5
	KonashiDigitalIO6
	KonashiDigitalIO7
)

type TactSwitch int

const (
	KonashiS1 TactSwitch = iota
)

type LED int

const (
	KonashiLED2 LED = 1 + iota
	KonashiLED3
	KonashiLED4
	KonashiLED5
)

type AioPin int

const (
	KonashiAnalogIO0 AioPin = iota
	KonashiAnalogIO1
	KonashiAnalogIO2
)

type I2CPin int

const (
	KonashiI2C_SDA I2CPin = 6 + iota
	KonashiI2C_SCL
)

type PinLevel int

const (
	KonashiLevelHigh PinLevel = iota
	KonashiLevelLow
)

type PinIOMode int

const (
	KonashiPinModeOutput PinIOMode = iota
	KonashiPinModeInput
)

type PinPullupMode int

const (
	KonashiPinModePullup PinPullupMode = iota
	KonashiPinModeNoPulls
)

const (
	analogReference = 1300
)

type PwmMode int

const (
	KonashiPWMModeDisable PwmMode = iota
	KonashiPWMModeEnable
	KonashiPWMModeEnableLED
	KonashiLEDPeriod = 10000
)

type UartMode int

const (
	KonashiUartModeDisable UartMode = iota
	KonashiUartModeEnable
	KonashiUartBaudrateRate2K4   = 0x000a // 2400bps
	KonashiUartBaudrateRate9K6   = 0x0028 // 9600bps
	KonashiUartBaudrateRate19K2  = 0x0050 // 19200bps
	KonashiUartBaudrateRate38K4  = 0x00a0 // 38400pbs
	KonashiUartBaudrateRate57K6  = 0x00f0 // 57600pbs
	KonashiUartBaudrateRate76K8  = 0x0140 // 76800pbs
	KonashiUartBaudrateRate115K2 = 0x01e0 // 115200pbs
)

type I2CMode int

const (
	KonashiI2CModeDisable I2CMode = iota
	KonashiI2CModeEnable
)

type I2CCondition int

const (
	KonashiI2CConditionStop I2CCondition = iota
	KonashiI2CConditionStart
	KonashiI2CConditionRestart
)

const (
	KonashiResultSuccess = iota
	KonashiResultFailure
)
