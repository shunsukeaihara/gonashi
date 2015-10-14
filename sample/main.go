package main

import (
	"log"
	"sync"
	"time"

	"github.com/flemay/gatt"
	"github.com/shunsukeaihara/gonashi"
)

func main() {
	g, err := gonashi.NewGonashi()
	if err != nil {
		return
	}
	g.Scan()
	var discovered map[string]*gonashi.Konashi

	ticker := time.NewTicker(time.Second * 20)

	select {
	case discovered = <-g.Discovered():
		break
	case <-ticker.C:
		log.Println("time out")
		g.StopScanning()
		return
	}
	g.StopScanning()
	log.Println(gatt.UUID16(0xFF00).String())
	log.Println(discovered)

	wg := new(sync.WaitGroup)
	for idStr, konashi := range discovered {
		log.Println(idStr)
		konashi.Connect()
		wg.Add(1)
		go func() {
			<-konashi.Connected
			log.Println("Connected")
			defer func() {
				konashi.DisConnect()
				<-konashi.Disconnected
				log.Println("Disconnected")
				wg.Done()
			}()
			for _, c := range konashi.DiscoverCharacteristics() {
				log.Println(c)
			}

		}()
	}
	wg.Wait()
}
