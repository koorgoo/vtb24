package exchange

import (
	"errors"
	"sort"
)

var (
	ErrNegativeAmount = errors.New("exchange: negative amount")
	ErrNoRate         = errors.New("exchange: no rate")
)

// errThreshold is returned by Func when provided amount does not match a
// threshold.
var errThreshold = errors.New("")

const nilThreshold = 0.0

type Func func(float64) (float64, error)

type Interface interface {
	Buy(float64) (float64, error)
	Sell(float64) (float64, error)
}

type Rate struct {
	Buy  float64
	Sell float64

	// Threshold sets the minimum amount for exchange. Threshold is used to
	// find a rate matching a provided amount only when a few rates are passed
	// into New().
	Threshold float64
}

func (r *Rate) doBuy(x float64) (float64, error)  { return r.exchange(x, r.Buy) }
func (r *Rate) doSell(x float64) (float64, error) { return r.exchange(x, r.Sell) }

func (r *Rate) exchange(x, rate float64) (y float64, err error) {
	switch {
	case x < 0:
		err = ErrNegativeAmount
	case x < r.Threshold:
		err = errThreshold
	default:
		y = x * rate
	}
	return
}

// New returns an Interface.
//
// The function panics if provided slice of rates has 0 length.
func New(rates ...Rate) Interface {
	switch len(rates) {
	case 0:
		panic("exchange: New needs 1 or more rates")
	case 1:
		return newRateEx(&rates[0])
	default:
		return newRatesEx(rates)
	}
}

func newRateEx(rate *Rate) Interface {
	rate.Threshold = nilThreshold
	return &rateEx{Rate: rate}
}

type rateEx struct{ Rate *Rate }

func (e *rateEx) Buy(x float64) (float64, error)  { return e.Rate.doBuy(x) }
func (e *rateEx) Sell(x float64) (float64, error) { return e.Rate.doSell(x) }

func newRatesEx(rates []Rate) Interface {
	e := new(ratesEx)
	// Use index to iterate over the slice not to copy structs.
	for i := range rates {
		// TODO: panic when multiple rates have same thresholds?
		e.Rates = append(e.Rates, &rates[i])
	}
	e.sortRates()
	return e
}

type ratesEx struct{ Rates []*Rate }

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
		y, err = chooseFunc(rate)(x)
		if err == errThreshold {
			continue
		}
		return
	}
	y, err = 0, ErrNoRate
	return
}
