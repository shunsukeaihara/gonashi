package main

import (
	"log"
	"time"

	"github.com/shunsukeaihara/gonashi"
)

func main() {
	g, err := gonashi.NewGonashi()
	if err != nil {
		return
	}
	g.Scan()
	var discovered map[string]*gonashi.Konashi
	func() {
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				log.Println("scanning...")
				discovered = g.GetDiscovered()
				log.Println(len(discovered))
				if len(discovered) > 0 {
					//一つでも見つかったら先に進む
					ticker.Stop()
					return
				}
			}
		}
	}()
	log.Println(discovered)
	for idStr, konashi := range discovered {
		log.Println(idStr)
		konashi.Connect()
	}
}
