package gonashi

import (
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

func (k *Konashi) Connect() {
	k.Peripheral.Device().Connect(k.Peripheral)
}

func (k *Konashi) DisConnect() {
	k.Peripheral.Device().CancelConnection(k.Peripheral)
}

func (k *Konashi) DiscoverCharacteristics() []*gatt.Service {
	s, _ := k.Peripheral.DiscoverServices([]gatt.UUID{gatt.UUID16(0xFF00)})
	//s, _ := k.Peripheral.DiscoverServices(nil)
	return s
}

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
