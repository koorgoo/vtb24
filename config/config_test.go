package config

import (
	"reflect"
	"testing"
)

var ParseJSONTests = []struct {
	Filename string
	Config   Config
	OK       bool
}{
	{
		"testdata/valid.json",
		Config{WebAddr: ":8000", TelegramToken: "test"},
		true,
	},
	{"testdata/no-web-addr.json", Config{}, false},
	{"testdata/no-telegram-token.json", Config{}, false},
	{"testdata/not-json.json", Config{}, false},
	{"testdata/does-not-exist.json", Config{}, false},
}

func TestParseJSON(t *testing.T) {
	for _, tt := range ParseJSONTests {
		t.Run(tt.Filename, func(t *testing.T) {
			c, err := ParseJSON(tt.Filename)
			if err == nil && !reflect.DeepEqual(tt.Config, c) {
				t.Errorf("want %v, got %v", tt.Config, c)
			}
			if ok := (err == nil); ok != tt.OK {
				t.Errorf("error: want %v, got %v", tt.OK, ok)
			}
		})
	}
}
