package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	RequestBody = `{"action":"{\"action\":\"currency\"}","scopeData":"{\"currencyRate\":\"ExchangePersonal\"}"}`
	RequestURL  = "https://www.vtb24.ru/services/ExecuteAction"
)

type Client struct {
	Client *http.Client
}

func (c *Client) Request() (*Response, error) {
	req := newRequest()
	return doRequest(c.Client, req)
}

func doRequest(client *http.Client, req *http.Request) (*Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var rr *Response
	err = json.NewDecoder(resp.Body).Decode(&rr)
	return rr, err
}

func newRequest() *http.Request {
	b := bytes.NewReader([]byte(RequestBody))
	r, err := http.NewRequest("POST", RequestURL, b)
	if err != nil {
		panic(err)
	}
	r.Header.Set("Content-Type", "application/json")
	return r
}

type Response struct {
	Items []*Item `json:"items"`
}

type Item struct {
	CurrencyGroupAbbr string    `json:"currencyGroupAbbr"`
	CurrencyAbbr      string    `json:"currencyAbbr"`
	Title             string    `json:"title"`
	Quantity          float64   `json:"quantity"`
	Buy               ItemValue `json:"buy"`
	BuyArrow          string    `json:"buyArrow"`
	Sell              ItemValue `json:"sell"`
	SellArrow         string    `json:"sellArrow"`
	Gradation         float64   `json:"gradation"`
	DateActiveFrom    ItemTime  `json:"dateActiveFrom"`
	IsMetal           bool      `json:"isMetal"`
}

type ItemValue float64

func (v *ItemValue) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	v2 := ItemValue(f)
	*v = v2
	return nil
}

type ItemTime time.Time

const (
	timePrefix = "/Date("
	timeSuffix = ")/"
)

func (d *ItemTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.TrimPrefix(s, timePrefix)
	s = strings.TrimSuffix(s, timeSuffix)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	var secs int64 = n / 1000
	var nsecs int64 = n - secs*1000
	d2 := time.Unix(secs, nsecs).UTC()
	*d = ItemTime(d2)
	return nil
}
