package gonashi

import (
	"strings"
	"sync"
	"time"

	"github.com/eapache/channels"
	"github.com/flemay/gatt"
)

var ClientOptions = []gatt.Option{
	gatt.MacDeviceRole(gatt.CentralManager),
}

type Gonashi struct {
	device     gatt.Device
	discovered *discovered
	connected  *konashiMap
}

func NewGonashi() (Gonashi, error) {
	device, err := gatt.NewDevice(ClientOptions...)
	if err != nil {
		return Gonashi{}, err
	}
	dp := discovered{
		konashiMap{
			sync.RWMutex{},
			map[string]*Konashi{},
			channels.NewOverflowingChannel(1)},
		make(chan struct{}),
	}

	cn := konashiMap{
		sync.RWMutex{},
		map[string]*Konashi{},
		channels.NewOverflowingChannel(1),
	}

	g := Gonashi{device, &dp, &cn}
	g.device.Handle(
		gatt.PeripheralDiscovered(g.onPeriphDiscovered),
		gatt.PeripheralConnected(g.onPeriphConnected),
		gatt.PeripheralDisconnected(g.onPeriphDisconnected),
	)
	device.Init(func(d gatt.Device, s gatt.State) {})
	return g, nil
}

type konashiMap struct {
	mu       sync.RWMutex
	konashis map[string]*Konashi
	Update   *channels.OverflowingChannel
}

type discovered struct {
	konashiMap
	Stop chan struct{}
}

func (km *konashiMap) GetKonashiMap() map[string]*Konashi {
	km.mu.RLock()
	defer km.mu.RUnlock()
	ret := map[string]*Konashi{}
	for id, konashi := range km.konashis {
		ret[id] = konashi
	}
	return ret
}

func (km *konashiMap) AddKonashi(k *Konashi) {
	km.mu.Lock()
	defer km.mu.Unlock()
	idStr := strings.ToUpper(k.Peripheral.ID())
	km.konashis[idStr] = k

	km.Update.In() <- km.konashis
}

func (km *konashiMap) DelKonashi(idStr string) {
	km.mu.Lock()
	defer km.mu.Unlock()
	delete(km.konashis, idStr)
	km.Update.In() <- km.konashis
}

func (km *konashiMap) Clear() {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.konashis = map[string]*Konashi{}
}

func (dp *discovered) AddDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	idStr := strings.ToUpper(p.ID())
	if k, ok := dp.konashis[idStr]; ok {
		k.Update(a, rssi)
	} else {
		dp.konashis[idStr] = NewKonashi(p, a, rssi)
	}
	dp.Update.In() <- dp.konashis
}

func (dp *discovered) Discard() {
	func() {
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				dp.mu.Lock()
				count := 0
				for idStr, d := range dp.konashis {
					now := time.Now()
					if now.Sub(d.T) > time.Second*20 {
						delete(dp.konashis, idStr)
						count++
					}
				}
				if count > 0 {
					dp.Update.In() <- dp.konashis
				}
				dp.mu.Unlock()
			case <-dp.Stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *Gonashi) Discovered() <-chan map[string]*Konashi {
	konaMap := <-g.discovered.Update.Out()
	ch := make(chan map[string]*Konashi, 1)
	ret, _ := konaMap.(map[string]*Konashi)
	ch <- ret
	return ch
}

func (g *Gonashi) Scan() {
	g.device.Scan([]gatt.UUID{}, false)
	go g.discovered.Discard()
}

func (g *Gonashi) StopScanning() {
	g.device.Scan([]gatt.UUID{}, false)
	g.device.StopScanning()
	g.discovered.Stop <- struct{}{}
}

func (g *Gonashi) onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		g.device.Scan([]gatt.UUID{}, false)
		return
	default:
		g.device.StopScanning()
		g.discovered.Stop <- struct{}{}
	}
}

func (g *Gonashi) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if strings.HasPrefix(a.LocalName, "konashi") {
		g.discovered.AddDiscovered(p, a, rssi)
	}
}

func (g *Gonashi) onPeriphConnected(p gatt.Peripheral, err error) {
	discovered := g.discovered.GetKonashiMap()
	idStr := strings.ToUpper(p.ID())
	if k, ok := discovered[idStr]; ok {
		k.SetPeripheral(p)
		g.connected.AddKonashi(k)
		k.Connected <- struct{}{}
	}
}

func (g *Gonashi) onPeriphDisconnected(p gatt.Peripheral, err error) {
	connected := g.connected.GetKonashiMap()
	idStr := strings.ToUpper(p.ID())
	if k, ok := connected[idStr]; ok {
		g.connected.DelKonashi(idStr)
		k.Disconnected <- struct{}{}
	}
}

func (g *Gonashi) GetDiscovered() map[string]*Konashi {
	return g.discovered.GetKonashiMap()
}
