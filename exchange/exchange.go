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

type Func func(float64) (float64, error)

type Interface interface {
	Buy(float64) (float64, error)
	Sell(float64) (float64, error)
	Rates() []Rate
}

type Rate struct {
	Buy  float64
	Sell float64

	// Threshold sets the minimum amount for exchange. Threshold is used to
	// find a rate matching a provided amount only when a few rates are passed
	// into New().
	Threshold Threshold

	inverted bool
}

func (r *Rate) doBuy(x float64) (float64, error)  { return doExchange(x, r.Buy, r.Threshold.Buy()) }
func (r *Rate) doSell(x float64) (float64, error) { return doExchange(x, r.Sell, r.Threshold.Sell()) }

func doExchange(x, rate, threshold float64) (y float64, err error) {
	switch {
	case x < 0:
		err = ErrNegativeAmount
	case x < threshold:
		err = errThreshold
	default:
		y = x * rate
	}
	return
}

func (r *Rate) invert() *Rate {
	s, b := 1/r.Buy, 1/r.Sell
	th := r.invertThreshold()
	return &Rate{Buy: b, Sell: s, Threshold: th, inverted: !r.inverted}
}

func (r *Rate) invertThreshold() Threshold {
	b, s := r.Threshold.Buy(), r.Threshold.Sell()
	if r.inverted {
		s, b = b/r.Buy, s/r.Sell
	} else {
		s, b = b*r.Buy, s*r.Sell
	}
	return NewThreshold(b, s)
}

type Threshold interface {
	Buy() float64
	Sell() float64
}

func NewThreshold(buy, sell float64) Threshold {
	return &threshold{buy: buy, sell: sell}
}

type threshold struct{ buy, sell float64 }

func (t *threshold) Buy() float64  { return t.buy }
func (t *threshold) Sell() float64 { return t.sell }

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
		// TODO: panic when rate has nil threshold?
		return newRatesEx(rates)
	}
}

var nilThreshold = NewThreshold(0, 0)

func newRateEx(rate *Rate) Interface {
	rate.Threshold = nilThreshold
	return &rateEx{Rate: rate}
}

type rateEx struct{ Rate *Rate }

func (e *rateEx) Buy(x float64) (float64, error)  { return e.Rate.doBuy(x) }
func (e *rateEx) Sell(x float64) (float64, error) { return e.Rate.doSell(x) }
func (e *rateEx) Rates() []Rate                   { return []Rate{*e.Rate} }

func newRatesEx(rates []Rate) Interface {
	e := new(ratesEx)
	// Use index to iterate over the slice not to copy structs.
	for i := range rates {
		// TODO: panic when multiple rates have same thresholds?
		rate := &rates[i]
		e.buy = append(e.buy, rate)
		e.sell = append(e.sell, rate)
		e.rates = append(e.rates, rate)
	}
	e.sortRates()
	return e
}

type ratesEx struct {
	buy   []*Rate
	sell  []*Rate
	rates []*Rate
}

func (e *ratesEx) sortRates() {
	sort.Slice(e.buy, func(i, j int) bool {
		return e.buy[i].Threshold.Buy() > e.buy[j].Threshold.Buy()
	})
	sort.Slice(e.sell, func(i, j int) bool {
		return e.sell[i].Threshold.Sell() > e.sell[j].Threshold.Sell()
	})
}

func (e *ratesEx) Buy(x float64) (float64, error)  { return exchange(x, e.buy, chooseBuy) }
func (e *ratesEx) Sell(x float64) (float64, error) { return exchange(x, e.sell, chooseSell) }

func (e *ratesEx) Rates() []Rate {
	v := make([]Rate, len(e.rates))
	for i, r := range e.rates {
		v[i] = *r
	}
	return v
}

type chooseFunc func(*Rate) Func

var (
	chooseBuy  chooseFunc = func(r *Rate) Func { return r.doBuy }
	chooseSell            = func(r *Rate) Func { return r.doSell }
)

func exchange(x float64, rates []*Rate, chooseFunc chooseFunc) (y float64, err error) {
	for _, rate := range rates {
		y, err = chooseFunc(rate)(x)
		if err == errThreshold {
			continue
		}
		return
	}
	y, err = 0, ErrNoRate
	return
}

func Invert(v Interface) Interface {
	rates := v.Rates()
	for i := range rates {
		rates[i] = *rates[i].invert()
	}
	return New(rates...)
}
