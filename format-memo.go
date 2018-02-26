package main

import (
	"fmt"
	"regexp"
	"strings"
)

var regMemoGT = regexp.MustCompile("Omschrijving: (.+) IBAN: .*")
var regMemoBA = regexp.MustCompile("Pasvolgnr:(\\d+) (\\d\\d-\\d\\d-\\d\\d\\d\\d) (\\d\\d:\\d\\d) Transactie:(.+) Term:(.+)")

// formatMemo formats the ING memo into a readable YNAB format
func formatMemo(code, memo string) string {
	memo = strings.Replace(memo, "\n", "", -1)

	// Switch according to code VZ/GT/IC
	switch code {
	case "BA":
		match := regMemoBA.FindAllStringSubmatch(memo, -1)

		if len(match) <= 0 {
			break
		}

		memo = fmt.Sprintf("Payment on %s at %s using card %s: %s @ %s", match[0][2], match[0][3], match[0][1], match[0][4], match[0][5])
	case "VZ":
		fallthrough
	case "OV":
		fallthrough
	case "GT":
		match := regMemoGT.FindAllStringSubmatch(memo, -1)

		if len(match) <= 0 {
			break
		}

		memo = match[0][1]
		memo = strings.Title(memo)
	case "IC":
		match := regMemoGT.FindAllStringSubmatch(memo, -1)

		if len(match) <= 0 {
			break
		}

		memo = "Incasso " + match[0][1]
		memo = strings.Title(memo)
	}

	// Finally, trim and replace obfuscation
	memo = strings.Trim(memo, " \t")            // Tabs and spaces
	memo = strings.Replace(memo, "A: ", "", -1) // Automatic payments

	// Return parsed data
	return memo
}
