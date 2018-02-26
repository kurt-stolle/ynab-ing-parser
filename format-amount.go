package main

import "strings"

// formatAmount formats the amount into YNAB format
func formatAmount(amount string) string {
	return strings.Replace(amount, ",", ".", -1)
}
