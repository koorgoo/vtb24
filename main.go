package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	RatesRequest = `{"action":"{\"action\":\"currency\"}","scopeData":"{\"currencyRate\":\"ExchangePersonal\"}"}`
	DefaultURL   = "https://www.vtb24.ru/services/ExecuteAction"
)

func main() {
	c := new(Client)
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

type Client struct{}

func (c *Client) Rates() (*RatesResponse, error) {
	req := NewRequest()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var rr *RatesResponse
	err = json.NewDecoder(resp.Body).Decode(&rr)
	return rr, err
}

func NewRequest() *http.Request {
	b := bytes.NewReader([]byte(RatesRequest))
	r, err := http.NewRequest("POST", DefaultURL, b)
	if err != nil {
		panic(err)
	}
	r.Header.Set("Content-Type", "application/json")
	return r
}

type RatesResponse struct {
	Items []*RateItem `json:"items"`
}

type RateItem struct {
	CurrencyGroupAbbr string  `json:"currencyGroupAbbr"`
	CurrencyAbbr      string  `json:"currencyAbbr"`
	Title             string  `json:"title"`
	Quantity          float64 `json:"quantity"`
	Buy               string  `json:"buy"`
	BuyArrow          string  `json:"buyArrow"`
	Sell              string  `json:"sell"`
	SellArrow         string  `json:"sellArrow"`
	Gradiation        int     `json:"gradiation"`
	DateActiveFrom    string  `json:"dateActiveFrom"`
	IsMetal           bool    `json:"isMetal"`
}
