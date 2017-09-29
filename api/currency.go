package api

import (
	"fmt"
	"strings"
)

const RUB = "RUB" // Russian ruble

const (
	AUD = "AUD" // Australian dollar
	CAD = "CAD" // Canadian dollar
	CHF = "CHF" // Swiss franc
	CNY = "CNY" // Chinese yuan (Renminbi)
	DKK = "DKK" // Danish krone
	EUR = "EUR" // Euro
	GBP = "GBP" // Pound sterling
	JPY = "JPY" // Japanese yen
	NOK = "NOK" // Norwegian krone
	NZD = "NZD" // New Zealand dollar
	PLN = "PLN" // Polish złoty
	SEK = "SEK" // Swedish krona
	USD = "USD" // United States dollar
)

// SplitCurrency returns a list of currencies from abbr.
// Empty dest means RUB.
func SplitCurrency(abbr string) (src, dest string) {
	a := strings.Split(abbr, "/")
	switch len(a) {
	case 1:
		return a[0], ""
	case 2:
		return a[0], a[1]
	default:
		panic(fmt.Sprintf("too many currencies in %q", abbr))
	}
}
