package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func retrieveMETAR(icao string) string {
	// Input checks
	if len(icao) != 4 {
		return ""
	}
	icao = strings.ToUpper(icao)

	// Prepare request
	h := http.Header{}
	h.Add("X-API-Key", os.Getenv("M_API_KEY"))

	uParsed, err := url.Parse(fmt.Sprintf("https://api.checkwx.com/metar/%s", icao))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	c := http.Client{}
	resp, err := c.Do(&http.Request{Header: h, URL: uParsed})
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
