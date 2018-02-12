package main

import "net/http"
import (
	"fmt"
	"io/ioutil"
	"time"
	"strings"
	"net/url"
)

type Vessel int

const (
	HLT = 0
	MASH = 1
	BOIL = 2
)

type Brewtroller struct {
	Hostname url.URL
}

func main() {
	tickChan := time.NewTicker(time.Millisecond * 100).C

	host, _ := url.Parse("http://192.168.1.179/btnic.cgi")
	controller := &Brewtroller{Hostname:*host}
	for {
		select {
		case <- tickChan:
			controller.brewtrollerStatus()
			controller.getTargetVolume(HLT)
		}
	}
}

func (bt *Brewtroller) brewtrollerStatus() {
	resp, err := http.Get(bt.Hostname.String()+ "?a")
	if (err != nil) {
		fmt.Errorf("fuck")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	words := strings.Fields(string(body))
	fmt.Println(words)

}

func (bt *Brewtroller) getTargetVolume(vessel Vessel) {
	resp, err := http.Get(fmt.Sprintf("%s?p%d",bt.Hostname.String(), vessel))
	if (err != nil) {
		fmt.Errorf("fuck")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	words := strings.Fields(string(body))
	fmt.Println(words)

}
