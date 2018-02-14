package main

import (
	"time"
	"net/url"
	"log"
	"github.com/influxdata/influxdb/client/v2"
)


func main() {
	tickChan := time.NewTicker(time.Millisecond * 100).C

	host, _ := url.Parse("http://192.168.1.179/btnic.cgi")
	controller := &Brewtroller{Hostname:*host}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://192.168.1.5:8086",
		Username: "root",
		Password: "root",
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <- tickChan:
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database: "brewtroller",
				Precision: "s",
			})
			if err != nil {
				log.Fatal(err)
			}

			controller.brewtrollerStatus()

			bp.AddPoint(controller.getTargetVolume(MASH))

			err = c.Write(bp)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}

