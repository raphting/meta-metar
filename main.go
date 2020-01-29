package main

import (
	"fmt"
	"regexp"
	"strconv"
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
	/*
		EHRD 262025Z AUTO 19014KT 160V220 9999 BKN035 OVC039 07/05 Q1010 NOSIG
		EHAM 262155Z 19015KT 9000 -RA FEW017 SCT025 BKN038 07/06 Q1009 BECMG 7000
		EBBR 262150Z 19012KT 9999 SCT048 08/05 Q1012 TEMPO 4500 RA BKN012
		EDDL 262150Z 16012KT CAVOK 06/04 Q1013 NOSIG
		EHBK 262225Z AUTO 20015KT 170V230 9999 NSC 07/05 Q1012
		EDDK 262220Z 13008KT 9999 SCT023 05/03 Q1015 NOSIG
		EDDF 262220Z VRB02KT 0200 R25R/0700N R25C/0550N R25L/0550N R18/0450N FG VV/// 02/02 Q1017 NOSIG
	*/

	metar := "EHRD 262025Z AUTO VRB19014KT 160V220 9999 BKN014 OVC039 07/05 Q1010 NOSIG"

	m := newMETAR(metar)

	m.markWind()
	m.markVisibility()
	m.markCloudbase()
	m.markCriticalWeather()
	m.markConvectiveClouds()
	fmt.Println(m)
}

func newMETAR(metar string) METAR {
	return METAR{metar, make([]alert, 0)}
}

func (m *METAR) markVisibility() {
	cloudbase := regexp.MustCompile(" (\\d{4}) ")
	viz := cloudbase.FindStringSubmatch(m.metarText)
	if len(viz) < 2 {
		return
	}

	// Omit error checks because regexp finds ints
	vizMeters, _ := strconv.Atoi(viz[1])
	if vizMeters > 6000 {
		return
	}

	indices := cloudbase.FindStringSubmatchIndex(m.metarText)
	if len(indices) < 4 {
		return
	}
	m.alerts = append(m.alerts, alert{
		startIndex: indices[2],
		endIndex:   indices[3],
		token:      LOW_VIZ,
	})
}

func (m *METAR) markConvectiveClouds() {
	weather := regexp.MustCompile("\\S*(?:CB|TCU)\\S*")
	indices := weather.FindAllStringIndex(m.metarText, -1)
	for _, i := range indices {
		m.alerts = append(m.alerts, alert{
			startIndex: i[0],
			endIndex:   i[1],
			token:      CONVECTIVE_CLOUDS,
		})
	}
}

func (m *METAR) markCriticalWeather() {
	//Description:
	//SH Showers
	//TS Thunderstorm

	//Precipitation:
	//GR Hail
	//GS Small Hail

	//Obscuration:
	//BR Mist
	//DU Dust
	//FG Fog
	//FU Smoke
	//HZ Haze
	//PY Spray
	//SA Sand
	//VA Volcanic Ash

	weather := regexp.MustCompile("\\S*(?:SH|TS|GR|GS|BR|DU|FG|FU|HZ|PY|SA|VA)\\S*")
	indices := weather.FindAllStringIndex(m.metarText, -1)
	for _, i := range indices {
		m.alerts = append(m.alerts, alert{
			startIndex: i[0],
			endIndex:   i[1],
			token:      WEATHER,
		})
	}
}

func (m *METAR) markCloudbase() {
	cloudbase := regexp.MustCompile("(?:FEW|SCT|BKN|OVC)(\\d{3})")
	clouds := cloudbase.FindAllStringSubmatch(m.metarText, -1)
	cloudsIndices := cloudbase.FindAllStringSubmatchIndex(m.metarText, -1)
	for i, c := range clouds {
		// Omit error checks because regexp finds ints
		cloudFeet, _ := strconv.Atoi(c[1])
		if cloudFeet < 15 {
			m.alerts = append(m.alerts, alert{
				startIndex: cloudsIndices[i][0],
				endIndex:   cloudsIndices[i][1],
				token:      LOW_CLOUDBASE,
			})
		}
	}
}

func (m *METAR) markWind() {
	windspeed := regexp.MustCompile("(?:VRB|\\d{3})(\\d{2})KT")
	wind := windspeed.FindStringSubmatch(m.metarText)
	if len(wind) == 0 {
		return
	}

	// Omit error checks because regexp finds ints
	windKT, _ := strconv.Atoi(wind[1])
	if windKT > 12 {
		indices := windspeed.FindStringSubmatchIndex(m.metarText)
		m.alerts = append(m.alerts, alert{
			startIndex: indices[0],
			endIndex:   indices[1],
			token:      STRONG_WIND,
		})
	}
}
