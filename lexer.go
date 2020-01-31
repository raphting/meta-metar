package main

import (
	"regexp"
	"strconv"
)

func newMETAR(metar string) METAR {
	return METAR{metar, make([]alert, 0)}
}

func (m *METAR) markVisibility() {
	cloudbase := regexp.MustCompile(` (\d{4}) `)
	viz := cloudbase.FindAllStringSubmatch(m.metarText, -1)
	indices := cloudbase.FindAllStringSubmatchIndex(m.metarText, -1)

	for i, v := range viz {
		// Prevent out-of-bounds
		if len(v) < 2 {
			continue
		}
		// Check club limits
		vizMeters, _ := strconv.Atoi(v[1])
		if vizMeters > 6000 {
			continue
		}

		// Prevent out-of-bounds
		if len(indices[i]) < 4 {
			continue
		}
		m.alerts = append(m.alerts, alert{
			startIndex: indices[i][2],
			endIndex:   indices[i][3],
			token:      LOW_VIZ,
		})
	}
}

func (m *METAR) markConvectiveClouds() {
	weather := regexp.MustCompile(`\S*(?:CB|TCU)\S*`)
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

	weather := regexp.MustCompile(`\S*(?:SH|TS|GR|GS|BR|DU|FG|FU|HZ|PY|SA|VA)\S*`)
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
	cloudbase := regexp.MustCompile(`(?:FEW|SCT|BKN|OVC)(\d{3})`)
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
	windspeed := regexp.MustCompile(`(?:VRB)?(?:\d{3})?(?:(\d{2})G)?(\d{2})KT`)
	wind := windspeed.FindAllStringSubmatch(m.metarText, -1)
	// Check if wind is given
	if len(wind) == 0 {
		return
	}

	indices := windspeed.FindAllStringSubmatchIndex(m.metarText, -1)

	for i, w := range wind {
		// prevent out-of-bounds
		if len(w) < 2 {
			continue
		}

		highWind := false
		w1, _ := strconv.Atoi(w[1])
		if w1 > 12 {
			highWind = true
		}

		if len(w) == 3 {
			w2, _ := strconv.Atoi(w[2])
			if w2 > 12 {
				highWind = true
			}
		}

		if highWind {
			m.alerts = append(m.alerts, alert{
				startIndex: indices[i][0],
				endIndex:   indices[i][1],
				token:      STRONG_WIND,
			})
		}
	}

}
