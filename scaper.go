package main

import (
	"time"
	"net/url"
)





func main() {
	tickChan := time.NewTicker(time.Millisecond * 100).C

	host, _ := url.Parse("http://192.168.1.179/btnic.cgi")
	controller := &Brewtroller{Hostname:*host}
	for {
		select {
		case <- tickChan:
			controller.brewtrollerStatus()
			controller.getTargetVolume(MASH)
		}
	}
}

