package main

import (
	"log"

	"github.com/koorgoo/vtb24/api"
)

func main() {
	c := new(api.Client)
	resp, err := c.Request()
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
