package exchange

import (
	"errors"
	"sort"
)

var (
	ErrNegativeAmount = errors.New("exchange: negative amount")
	ErrNoRate         = errors.New("exchange: no rate")
)

type Func func(float64) (float64, error)

type Interface interface {
	Buy(float64) (float64, error)
	Sell(float64) (float64, error)
}

type Rate struct {
	Buy  float64
	Sell float64

	// Threshold sets the minimum amount for exchange. Threshold is used only
	// in exchanger returned by NewWithThresholds().
	Threshold float64
}

func New(rate Rate) Interface {
	return &basicEx{Rate: &rate}
}

type basicEx struct{ Rate *Rate }

func (e *basicEx) Buy(x float64) (float64, error)  { return exchange(x, e.Rate.Buy) }
func (e *basicEx) Sell(x float64) (float64, error) { return exchange(x, e.Rate.Sell) }

func exchange(x, rate float64) (y float64, err error) {
	if x < 0 {
		err = ErrNegativeAmount
	} else {
		y = x * rate
	}
	return
}

func NewWithThresholds(rates ...Rate) Interface {
	e := new(withThresholdsEx)
	for i := range rates {
		e.entries = append(e.entries, &basicEx{Rate: &rates[i]})
	}
	e.sortEntries()
	return e
}

type withThresholdsEx struct {
	entries []*basicEx
}

func (e *withThresholdsEx) sortEntries() {
	sort.Slice(e.entries, func(i, j int) bool {
		return e.entries[i].Rate.Threshold > e.entries[j].Rate.Threshold
	})
}

func (e *withThresholdsEx) Buy(x float64) (float64, error) {
	return e.exchange(x, func(e Interface) Func { return e.Buy })
}

func (e *withThresholdsEx) Sell(x float64) (float64, error) {
	return e.exchange(x, func(e Interface) Func { return e.Sell })
}

type xChoose func(Interface) Func

func (e *withThresholdsEx) exchange(x float64, choose xChoose) (y float64, err error) {
	for _, ex := range e.entries {
		if ex.Rate.Threshold > x {
			continue
		}
		y, err = choose(ex)(x)
		return
	}
	_, err = exchange(x, 1)
	if err == nil {
		err = ErrNoRate
	}
	return
}
