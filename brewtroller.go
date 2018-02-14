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
	"errors"
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
	responseFields, err := makeRequest(bt, "a")
	if err != nil {
		print(err)
	}
	fmt.Println(responseFields)

}

func (bt *Brewtroller) getTargetVolume(vessel Vessel) *client.Point {
	responseFields, err := makeRequest(bt, fmt.Sprintf("|%d", vessel))
	if (err != nil) {
		log.Println("Error making target volume request")
		return nil
	}

	retVal, _ := strconv.Atoi(responseFields[1])

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

func makeRequest(bt *Brewtroller, requestParams string) ([]string, error) {
	response, err := http.Get(fmt.Sprintf("%s?%s",bt.Hostname.String(), requestParams))
	if (err != nil) {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Non 200 response received for request %s", response.Request))
	}

	responseBody, err := parseBody(response)
	if (err != nil) {
		return nil, err
	}

	return responseBody, nil
}

func parseBody(response *http.Response) ([]string, error) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if (err != nil) {
		return nil, err
	}

	words := strings.Split(string(body)[1:], ",")
	for i := range words {
		words[i] = strings.Trim(words[i], "\"")
	}

	return words, nil
}
