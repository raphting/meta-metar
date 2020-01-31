package main

import (
	"fmt"
	"log"
	"net/http"
)

type tokenType uint8

const (
	STRONG_WIND tokenType = iota
	LOW_VIZ
	LOW_CLOUDBASE
	CONVECTIVE_CLOUDS
	WEATHER
)

type alert struct {
	startIndex int
	endIndex   int
	token      tokenType
}

type METAR struct {
	metarText string
	alerts    []alert
}

func main() {
	http.HandleFunc("/", metarRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func metarRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	icaoCode := query.Get("icao")
	metar := retrieveMETAR(icaoCode)
	description := retrieveAirportName(icaoCode)
	m := newMETAR(metar)

	m.markWind()
	m.markVisibility()
	m.markCloudbase()
	m.markCriticalWeather()
	m.markConvectiveClouds()

	_, err := fmt.Fprint(w, parseMetarTemplate(description, m.colorAreas("<b>", "</b>")))
	if err != nil {
		fmt.Println(err.Error())
	}
}
