package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/koorgoo/telegram"
	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
	"github.com/koorgoo/vtb24/chat"
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

const RatesRetryTimeout = time.Minute

var cfgPath = flag.String("config.file", "config.json", "path to configuration file")

func main() {
	flag.Parse()

	cfg, err := config.ParseJSON(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	errc := make(chan error, 1)
	termc := make(chan os.Signal)
	signal.Notify(termc, os.Interrupt, syscall.SIGTERM)

	ex, err := GetDefaultEx()
	if err != nil {
		errc <- err
	}

	var rates atomic.Value
	rates.Store(ex)

	go func() {
		for {
			t := cfg.RatesTimeout
			e, err := GetDefaultEx()
			if err == nil {
				rates.Store(e)
			} else {
				log.Printf("failed to update rates: %s", err)
				t = RatesRetryTimeout
			}
			time.Sleep(t)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.WebAddr, nil); err != nil {
			errc <- err
		}
	}()

	go func() {
		bot, err := telegram.NewBot(context.TODO(), cfg.TelegramToken)
		if err != nil {
			errc <- err
			return
		}

		go func(updatec <-chan *telegram.Update) {
			for update := range updatec {
				if update.Message == nil {
					continue
				}
				if update.Message.Text == nil {
					continue
				}

				n, err := strconv.ParseFloat(*update.Message.Text, 64)
				if err != nil {
					_, _ = bot.SendMessage(context.TODO(), &telegram.TextMessage{
						ChatID: update.Message.Chat.ID,
						Text:   "Я понимаю только числа.",
					})
					continue
				}

				ex := rates.Load().([]bank.Ex)
				text, mode := chat.MakeMessage(n, ex)
				if text == "" {
					_, _ = bot.SendMessage(context.TODO(), &telegram.TextMessage{
						ChatID: update.Message.Chat.ID,
						Text:   fmt.Sprintf("Не удалось обменять %v.", n),
					})
					continue
				}

				_, err = bot.SendMessage(context.TODO(), &telegram.TextMessage{
					ChatID:    update.Message.Chat.ID,
					Text:      text,
					ParseMode: mode,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}(bot.Updates())

		go func(errorc <-chan error) {
			for err := range errorc {
				log.Println(err)
			}
		}(bot.Errors())
	}()

	select {
	case <-termc:
	case err := <-errc:
		log.Fatal(err)
	}
}

func GetDefaultEx() ([]bank.Ex, error) {
	c := new(api.Client)
	resp, err := c.Request()
	if err != nil {
		return nil, err
	}
	ex := bank.ParseEx(resp)
	ex = bank.FilterEx(ex, DefaultFilters...)
	return ex, nil
}
