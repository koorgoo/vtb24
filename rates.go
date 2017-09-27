package vtb24

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	RatesRequest = `{"action":"{\"action\":\"currency\"}","scopeData":"{\"currencyRate\":\"ExchangePersonal\"}"}`
	DefaultURL   = "https://www.vtb24.ru/services/ExecuteAction"
)

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

type RateValue float64

func (v *RateValue) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	v2 := RateValue(f)
	*v = v2
	return nil
}

type RateTime time.Time

const (
	rateTimePrefix = "/Date("
	rateTimeSuffix = ")/"
)

func (d *RateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.TrimPrefix(s, rateTimePrefix)
	s = strings.TrimSuffix(s, rateTimeSuffix)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	var secs int64 = n / 1000
	var nsecs int64 = n - secs*1000
	d2 := time.Unix(secs, nsecs).UTC()
	*d = RateTime(d2)
	return nil
}

type RateItem struct {
	CurrencyGroupAbbr string    `json:"currencyGroupAbbr"`
	CurrencyAbbr      string    `json:"currencyAbbr"`
	Title             string    `json:"title"`
	Quantity          float64   `json:"quantity"`
	Buy               RateValue `json:"buy"`
	BuyArrow          string    `json:"buyArrow"`
	Sell              RateValue `json:"sell"`
	SellArrow         string    `json:"sellArrow"`
	Gradation         float64   `json:"gradation"`
	DateActiveFrom    RateTime  `json:"dateActiveFrom"`
	IsMetal           bool      `json:"isMetal"`
}
