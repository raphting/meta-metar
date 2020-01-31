package main

import "strings"

func parseMetarTemplate(header, metarInfo string) string {
	metarTemplate := `
<html>
	<header>
		<style type="text/css">
			body {
				font-family: "Source Code Pro";
				background-color: #001f3f;
				color: #7FDBFF;
			}
			b {
				color: #FFDC00;
			}
		</style>
	</header>
	<body>
		<h1>#HEADER#</h1>
		<p>#METAR#</p>
	</body>
</html>
`
	m := strings.Replace(metarTemplate, "#METAR#", metarInfo, -1)
	return strings.Replace(m, "#HEADER#", header, -1)
}
