package gonashi

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/eapache/channels"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type Gonashi struct {
	device     gatt.Device
	discovered *discovered
	connected  *konashiMap
}

func NewGonashi() (Gonashi, error) {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
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
		gatt.PeripheralDisconnected(g.onPeriphDisconnected),
		gatt.PeripheralConnected(g.onPeriphConnected),
	)
	device.Init(func(d gatt.Device, s gatt.State) {})

	return g, nil
}

type Konashi struct {
	Peripheral    gatt.Peripheral
	Advertisement *gatt.Advertisement
	Rssi          int
	T             time.Time
	connected     chan bool
}

func (k *Konashi) Connect() {
	// idStr := k.Peripheral.ID()
	k.Peripheral.Device().Connect(k.Peripheral)
	// ToDo: gonashiのdiscoveredに追加して、そのmapは、onPeriphDiscoveredで参照してchanでやり取りする
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

func (dp *konashiMap) GetDiscovered() map[string]*Konashi {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	ret := map[string]*Konashi{}
	for id, konashi := range dp.konashis {
		ret[id] = konashi
	}
	return ret
}

func (km *konashiMap) AddKonashi(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	km.mu.Lock()
	defer km.mu.Unlock()
	idStr := strings.ToUpper(p.ID())
	if k, ok := km.konashis[idStr]; ok {
		k.Peripheral = p
		k.Advertisement = a
		k.Rssi = rssi
		k.T = time.Now()
	} else {
		km.konashis[idStr] = &Konashi{p, a, rssi, time.Now(), make(chan bool)}
	}

	km.Update.In() <- km.konashis
}

func (km *konashiMap) Clear() {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.konashis = map[string]*Konashi{}
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
				dp.Update.In() <- dp.konashis
				dp.mu.Unlock()
			case <-dp.Stop:
				log.Println("stop")
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *Gonashi) Discoverd() <-chan interface{} {
	//ToDo: インターフェイスじゃなくてmap[string]*Konashiにしたいけど……
	return g.discovered.Update.Out()
}

func (g *Gonashi) Scan() {
	g.device.Scan([]gatt.UUID{}, false)
	go g.discovered.Discard()
}

func (g *Gonashi) StopScanning() {
	g.device.Scan([]gatt.UUID{}, false)
	g.device.StopScanning()
	g.discovered.Stop <- struct{}{}
	g.discovered.Clear()
}

func (g *Gonashi) onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		g.device.Scan([]gatt.UUID{}, false)
		return
	default:
		g.device.StopScanning()
		g.discovered.Stop <- struct{}{}
		g.discovered.Clear()
	}
}

func (g *Gonashi) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if strings.HasPrefix(a.LocalName, "konashi") {
		g.discovered.AddKonashi(p, a, rssi)
	}
}

func (g *Gonashi) onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

func (g *Gonashi) onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
}

func (g *Gonashi) GetDiscovered() map[string]*Konashi {
	return g.discovered.GetDiscovered()
}
