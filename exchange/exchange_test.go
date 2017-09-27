package exchange

import (
	"fmt"
	"testing"
)

type Table map[float64]float64

var BasicTests = []struct {
	Rate Rate
	Buy  Table
	Sell Table
}{
	{
		Rate: Rate{Buy: 2, Sell: 3},
		Buy:  Table{1: 2, 2.5: 5},
		Sell: Table{1: 3, 2.5: 7.5},
	},
}

func TestNew(t *testing.T) {
	for i := range BasicTests {
		tt := BasicTests[i]
		t.Run(fmt.Sprintf("%+v", tt.Rate), func(t *testing.T) {
			e := New(tt.Rate)
			for a, b := range tt.Buy {
				if c := e.Buy(a); c != b {
					t.Errorf("buy %v: want %v, got %v", a, b, c)
				}
			}
			for a, b := range tt.Sell {
				if c := e.Sell(a); c != b {
					t.Errorf("sell %v: want %v, got %v", a, b, c)
				}
			}
		})
	}
}
