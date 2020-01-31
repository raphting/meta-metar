package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func retrieveMETAR(icao string) string {
	// Input checks
	if len(icao) != 4 {
		return ""
	}
	icao = strings.ToUpper(icao)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.checkwx.com/metar/%s", icao), nil)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	req.Header.Add("X-API-Key", os.Getenv("M_API_KEY"))

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("WX response was not 200")
		return ""
	}

	// Analyze response
	type wxFormat struct {
		Data []string `json:"data"`
	}

	f := wxFormat{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if len(f.Data) != 1 {
		return ""
	}
	return f.Data[0]
}
