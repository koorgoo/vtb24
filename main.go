package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
	"github.com/koorgoo/vtb24/config"
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
		api.EUR, api.RUB,
	),
}

var cfgPath = flag.String("config.file", "config.json", "path to configuration file")

func main() {
	flag.Parse()

	cfg, err := config.ParseJSON(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	_ = cfg

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
