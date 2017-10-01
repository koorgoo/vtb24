package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
	"github.com/koorgoo/vtb24/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	termc := make(chan os.Signal)
	signal.Notify(termc, os.Interrupt, syscall.SIGTERM)

	servec := make(chan error)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.WebAddr, nil); err != nil {
			servec <- err
		}
	}()

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

	select {
	case <-termc:
	case err := <-servec:
		log.Fatal(err)
	}
}
