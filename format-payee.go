package main

import "strings"

// formatPayee formats the payee into YNAB format
func formatPayee(pyee string) string {
	// Trim into uniform format
	pyee = strings.ToUpper(pyee)
	pyee = strings.Trim(pyee, " \t\n")
	pyee = strings.Replace(pyee, ".", "", -1)
	pyee = strings.Replace(pyee, ",?", " EN ", -1)

	// Remove the following words that obfuscate the output
	for _, s := range []string{"HR ", "MEJ ", "DR ", "MR "} {
		pyee = strings.Replace(pyee, s, "", -1)
	}

	return pyee
}
