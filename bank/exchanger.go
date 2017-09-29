package bank

import (
	"fmt"

	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/exchange"
)

type Exchanger struct {
	// Src is a currency to exchange.
	Src string
	// Dst is a currency to exchange to.
	Dst string
	// Group is a currency group affecting rates.
	Group string
	exchange.Interface
}

func (e *Exchanger) String() string {
	group := api.GroupText(e.Group)
	return fmt.Sprintf("%s › %s (%s)", e.Src, e.Dst, group)
}

func ParseExchangers(resp *api.Response) []*Exchanger {
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
			Threshold: item.Gradation,
		})
	}

	var v []*Exchanger
	for src := range m {
		for dst := range m[src] {
			for group, rates := range m[src][dst] {
				e := exchange.New(rates...)
				ex := &Exchanger{Src: src, Dst: dst, Group: group, Interface: e}
				v = append(v, ex)
			}
		}
	}
	return v
}
