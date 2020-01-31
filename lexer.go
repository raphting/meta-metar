package main

import (
	"regexp"
	"strconv"
)

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
	windspeed := regexp.MustCompile("(?:VRB)?(?:\\d{3})?(\\d{2})KT")
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
