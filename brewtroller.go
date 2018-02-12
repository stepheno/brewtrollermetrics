package main

import (
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"log"
	"strconv"
)

type Vessel int

const (
	HLT Vessel = 0
	MASH = 1
	BOIL = 2
)

var vesselName = map[Vessel]string{
	HLT:"HLT",
	MASH:"MASH",
	BOIL:"BOIL",
}

type Brewtroller struct {
	Hostname url.URL
}

func (bt *Brewtroller) brewtrollerStatus() {
	resp, err := http.Get(bt.Hostname.String()+ "?a")
	if (err != nil) {
		fmt.Errorf("fuck")
	}
	words := parseBody(resp, err)
	fmt.Println(words)

}

func (bt *Brewtroller) getTargetVolume(vessel Vessel) *client.Point {
	resp, err := http.Get(fmt.Sprintf("%s?|%d",bt.Hostname.String(), vessel))
	if (err != nil) {
		fmt.Errorf("fuck")
	}
	words := parseBody(resp, err)

	retVal, _ := strconv.Atoi(strings.Trim(words[1], "\""))

	tags := map[string]string{"vessel": vesselName[vessel]}
	fields := map[string]interface{}{
		"value": retVal / 1000.0,
	}
	pt, err := client.NewPoint("target_volume", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pt)

	return pt

}
func parseBody(resp *http.Response, err error) []string {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	words := strings.Split(string(body)[1:], ",")
	return words
}
