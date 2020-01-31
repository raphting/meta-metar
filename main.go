package main

import (
	"github.com/gookit/color"
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

	metar := retrieveMETAR("EDDH")

	m := newMETAR(metar)

	m.markWind()
	m.markVisibility()
	m.markCloudbase()
	m.markCriticalWeather()
	m.markConvectiveClouds()

	out := m.colorAreas()
	color.Println(out)
}
