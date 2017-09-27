package main

import (
	"log"

	"github.com/koorgoo/vtb24"
)

func main() {
	c := new(vtb24.Client)
	resp, err := c.Rates()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Items {
		switch item.CurrencyAbbr {
		case "USD", "EUR":
			log.Printf("%v", item)
		}
	}
}
