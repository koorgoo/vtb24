package main

import (
	"fmt"
	"log"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
)

var DefaultFilters = []bank.ExFilter{
	bank.WithGroup(
		api.GroupCash,
		api.GroupCashDesk,
		api.GroupCentralBank,
		api.GroupTele,
	),
	bank.WithSrcDst(
		api.USD, api.RUB,
		api.RUB, api.USD,
		api.EUR, api.RUB,
		api.RUB, api.EUR,
	),
}

func main() {
	c := new(api.Client)
	resp, err := c.Request()
	if err != nil {
		log.Fatal(err)
	}

	ex := bank.ParseEx(resp)
	ex = bank.FilterEx(ex, DefaultFilters...)

	for _, e := range ex {
		fmt.Println(e)
	}
}
