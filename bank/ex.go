package bank

import (
	"fmt"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/exchange"
)

type Ex interface {
	// Src returns a currency to exchange (buy or sell).
	Src() string
	// Dst returns a currency to exchange to (sold or bought).
	Dst() string
	// Group returns a currency group affecting rates.
	Group() string

	exchange.Interface
}

type ex struct {
	src, dst, group string
	exchange.Interface
}

func (e *ex) Src() string   { return e.src }
func (e *ex) Dst() string   { return e.dst }
func (e *ex) Group() string { return e.group }

func (e *ex) String() string {
	group := api.GroupText(e.group)
	return fmt.Sprintf("%s › %s (%s)", e.src, e.dst, group)
}

func ParseEx(resp *api.Response) []Ex {
	// Group rates by src, dst, and group.
	m := map[string]map[string]map[string][]exchange.Rate{}
	for _, item := range resp.Items {
		src, dst := api.SplitCurrency(item.CurrencyAbbr)
		if dst == "" {
			dst = api.RUB
		}
		if _, ok := m[src]; !ok {
			m[src] = map[string]map[string][]exchange.Rate{}
		}
		if _, ok := m[src][dst]; !ok {
			m[src][dst] = map[string][]exchange.Rate{}
		}
		group := item.CurrencyGroupAbbr
		m[src][dst][group] = append(m[src][dst][group], exchange.Rate{
			Buy:       float64(item.Buy),
			Sell:      float64(item.Sell),
			Threshold: exchange.NewThreshold(item.Gradation, item.Gradation),
		})
	}

	var v []Ex
	for src := range m {
		for dst := range m[src] {
			for group, rates := range m[src][dst] {
				e := exchange.New(rates...)
				v = append(v, &ex{src: src, dst: dst, group: group, Interface: e})
			}
		}
	}
	return v
}

func FilterEx(v []Ex, filters ...ExFilter) []Ex {
	var a []Ex
	for _, ex := range v {
		if filterEx(ex, filters) {
			a = append(a, ex)
		}
	}
	return a
}

func filterEx(e Ex, filters []ExFilter) bool {
	for _, filter := range filters {
		if !filter(e) {
			return false
		}
	}
	return true
}

type ExFilter func(Ex) bool

func WithSrcDst(srcdst ...string) ExFilter {
	if len(srcdst)%2 == 1 {
		panic("odd number of src and dst")
	}
	return func(e Ex) bool {
		for i := 0; i < len(srcdst); i = i + 2 {
			if e.Src() == srcdst[i] && e.Dst() == srcdst[i+1] {
				return true
			}
		}
		return false
	}
}

func WithGroup(groups ...string) ExFilter {
	return func(e Ex) bool {
		for i := range groups {
			if e.Group() == groups[i] {
				return true
			}
		}
		return false
	}
}

func Invert(e Ex) Ex {
	i := exchange.Invert(e)
	return &ex{src: e.Dst(), dst: e.Src(), group: e.Group(), Interface: i}
}
