package main

func (m METAR) colorAreas(startMarker, endMarker string) string {
	alertArea := make(map[int]int)
	for _, warn := range m.alerts {
		alertArea[warn.startIndex] = warn.endIndex
	}

	out := ""
	nextEnd := -1
	for i, metarChar := range m.metarText {
		if end, found := alertArea[i]; found {
			nextEnd = end - 1 // End is exclusive
			out += startMarker
		}

		out += string(metarChar)

		if i == nextEnd {
			out += endMarker
		}
	}

	return out
}
