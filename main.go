package main

import (
	"fmt"
	"log"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
)

func main() {
	c := new(api.Client)
	resp, err := c.Request()
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range bank.ParseExchangers(resp) {
		fmt.Println(e)
	}
}
