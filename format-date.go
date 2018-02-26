package main

import (
	"log"
	"regexp"
	"strings"
)

var regexDate = regexp.MustCompile("(\\d\\d\\d\\d)(\\d\\d)(\\d\\d)")

// formatDate formats the date given by ING into YNAB format
func formatDate(date string) string {
	// Match by regex
	match := regexDate.FindAllStringSubmatch(date, -1)

	if len(match) <= 0 {
		log.Fatalf("Failed to parse date, format unknown: %s\n", date)
	}

	// We have to use the weird american MM/DD/YYYY format instead of the superior YYYY-MM-DD or DD-MM-YYYY format
	return strings.Join([]string{match[0][2], match[0][3], match[0][1]}, "/")
}
