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

	metar := "EHRD 312025Z AUTO 21017KT 170V240 9999 SCT013 OVC015 12/10 Q1005 BECMG BKN014"

	m := newMETAR(metar)

	m.markWind()
	m.markVisibility()
	m.markCloudbase()
	m.markCriticalWeather()
	m.markConvectiveClouds()

	out := m.colorAreas()

	color.Println(out)
}
