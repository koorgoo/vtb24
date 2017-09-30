package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	WebAddr           string `json:"web_addr"`
	TelegramToken     string `json:"telegram_token"`
	AmazonWishListURL string `json:"amazon_wish_list_url"`
	DonateCardNumber  string `json:"donate_card_number"`
}

var (
	errWebAddr       = errors.New("invalid web addr")
	errTelegramToken = errors.New("invalid telegram token")
)

func (c *Config) validate() error {
	if c.WebAddr == "" {
		return errWebAddr
	}
	if c.TelegramToken == "" {
		return errTelegramToken
	}
	return nil
}

func ParseJSON(filename string) (c Config, err error) {
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		err = fmt.Errorf("config: %s", err)
		return
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&c); err != nil {
		err = fmt.Errorf("config: %s: %s", filename, err)
		return
	}
	if err = c.validate(); err != nil {
		err = fmt.Errorf("config: validation: %s", err)
		return
	}
	return
}
