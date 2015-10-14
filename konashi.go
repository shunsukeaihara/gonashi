package gonashi

import (
	"log"
	"sync"
	"time"

	"github.com/flemay/gatt"
)

type Konashi struct {
	Peripheral    gatt.Peripheral
	Advertisement *gatt.Advertisement
	Rssi          int
	T             time.Time
	Connected     chan struct{}
	Disconnected  chan struct{}
	mu            sync.RWMutex
}

func NewKonashi(p gatt.Peripheral, a *gatt.Advertisement, rssi int) *Konashi {

	return &Konashi{p, a, rssi, time.Now(),
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		sync.RWMutex{},
	}
}

func (k *Konashi) Connect() {
	k.Peripheral.Device().Connect(k.Peripheral)
}

func (k *Konashi) DisConnect() {
	k.Peripheral.Device().CancelConnection(k.Peripheral)
}

func (k *Konashi) DiscoverCharacteristics() []*gatt.Characteristic {
	s, err := k.Peripheral.DiscoverServices([]gatt.UUID{gatt.MustParseUUID("229bff0003fb40da98a7b0def65c2d4b")})
	if err != nil || len(s) == 0 {
		log.Println("Service Not Found")
	}
	cs, err := k.Peripheral.DiscoverCharacteristics(nil, s[0])
	return cs
}

// func (k *Konashi) StoreCharacteristics() {
// 	cs := k.DiscoverCharacteristics()
// 	k.mu.Lock()
// 	defer k.mu.Unlock()

// 	for _, c := range cs {
// 		k.characteristics[strings.ToUpper(c.UUID().String())] = c
// 	}
// }

func (k *Konashi) SetPeripheral(p gatt.Peripheral) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Peripheral = p
}

func (k *Konashi) SetMTU(mtu uint16) error {
	return k.Peripheral.SetMTU(mtu)
}

func (k *Konashi) Update(a *gatt.Advertisement, rssi int) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Advertisement = a
	k.Rssi = rssi
	k.T = time.Now()
}

func (k *Konashi) ReadCharacteristic(c *gatt.Characteristic) ([]byte, error) {
	return k.Peripheral.ReadCharacteristic(c)
}

func (k *Konashi) WriteCharacteristic(c *gatt.Characteristic, b []byte, noRsp bool) error {
	return nil
}

func (k *Konashi) SetNotifyValue(c *gatt.Characteristic) error {
	//f func(*gatt.Characteristic, []byte, error)
	return nil
}

func (k *Konashi) PinMode(pin DioPin, mode PinIOMode) error {
	var pioSetting uint8
	pioSetting |= 0x01 << pin
	b := make([]byte, 1)
	b[0] = byte(pioSetting)
	k.WriteCharacteristic(KonashiPioSetting, b, true)
	return nil
}
