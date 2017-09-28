package exchange

import (
	"fmt"
	"testing"
)

func TestNew_panicWithoutRates(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("New must panic")
		}
	}()
	New([]Rate{}...)
}

type Table map[float64]struct {
	Amount float64
	Err    error
}

var InterfaceTests = []struct {
	Rates []Rate
	Buy   Table
	Sell  Table
}{
	// Single rate tests
	{
		Rates: []Rate{
			Rate{Buy: 2, Sell: 3},
		},
		Buy: Table{
			-1:  {0, ErrNegativeAmount},
			0:   {0, nil},
			1:   {2, nil},
			2.5: {5, nil},
		},
		Sell: Table{
			-1:  {0, ErrNegativeAmount},
			0:   {0, nil},
			1:   {3, nil},
			2.5: {7.5, nil},
		},
	},
	// Multiple rates tests
	{
		Rates: []Rate{
			{Buy: 2, Sell: 3, Threshold: 10},
			{Buy: 3, Sell: 6, Threshold: 20},
		},
		Buy: Table{
			-1:  {0, ErrNegativeAmount},
			0:   {0, ErrNoRate},
			1:   {0, ErrNoRate},
			10:  {20, nil},
			15:  {30, nil},
			20:  {60, nil},
			100: {300, nil},
		},
		Sell: Table{
			-1:  {0, ErrNegativeAmount},
			0:   {0, ErrNoRate},
			1:   {0, ErrNoRate},
			10:  {30, nil},
			15:  {45, nil},
			20:  {120, nil},
			100: {600, nil},
		},
	},
}

func TestNew(t *testing.T) {
	for i := range InterfaceTests {
		tt := InterfaceTests[i]
		e := New(tt.Rates...)
		t.Run(fmt.Sprintf("%+v", tt.Rates), func(t *testing.T) {
			test(t, e, tt.Buy, tt.Sell)
		})
	}
}

func test(t *testing.T, e Interface, Buy, Sell Table) {
	t.Run("buy", func(t *testing.T) { testFunc(t, e.Buy, Buy) })
	t.Run("sell", func(t *testing.T) { testFunc(t, e.Sell, Sell) })
}

func testFunc(t *testing.T, exchange Func, table Table) {
	for amount, r := range table {
		n, err := exchange(amount)
		if n != r.Amount {
			t.Errorf("%v: want %v, got %v", amount, r.Amount, n)
		}
		if err != r.Err {
			t.Errorf("%v: error: want %q, got %q", amount, r.Err, err)
		}
	}
}
