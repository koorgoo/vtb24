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

func (r *Rate) doBuy(x float64) (float64, error)  { return exchange(x, r.Buy) }
func (r *Rate) doSell(x float64) (float64, error) { return exchange(x, r.Sell) }

func exchange(x, rate float64) (y float64, err error) {
	if x < 0 {
		err = ErrNegativeAmount
	} else {
		y = x * rate
	}
	return
}

func New(rate Rate) Interface {
	return &rateEx{Rate: &rate}
}

type rateEx struct{ Rate *Rate }

func (e *rateEx) Buy(x float64) (float64, error)  { return e.Rate.doBuy(x) }
func (e *rateEx) Sell(x float64) (float64, error) { return e.Rate.doSell(x) }

func NewWithThresholds(rates ...Rate) Interface {
	return newRatesEx(rates...)
}

func newRatesEx(rates ...Rate) Interface {
	e := new(ratesEx)
	for i := range rates {
		e.Rates = append(e.Rates, &rates[i])
	}
	e.sortRates()
	return e
}

type ratesEx struct {
	Rates []*Rate
}

func (e *ratesEx) sortRates() {
	sort.Slice(e.Rates, func(i, j int) bool {
		return e.Rates[i].Threshold > e.Rates[j].Threshold
	})
}

func (e *ratesEx) Buy(x float64) (float64, error)  { return e.exchange(x, chooseBuy) }
func (e *ratesEx) Sell(x float64) (float64, error) { return e.exchange(x, chooseSell) }

type chooseFunc func(*Rate) Func

var (
	chooseBuy  chooseFunc = func(r *Rate) Func { return r.doBuy }
	chooseSell            = func(r *Rate) Func { return r.doSell }
)

func (e *ratesEx) exchange(x float64, chooseFunc chooseFunc) (y float64, err error) {
	for _, rate := range e.Rates {
		if rate.Threshold > x {
			continue
		}
		y, err = chooseFunc(rate)(x)
		return
	}
	_, err = exchange(x, 1)
	if err == nil {
		err = ErrNoRate
	}
	return
}
