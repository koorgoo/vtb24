package main

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
	// Meta is an optional extra information about exchanger.
	Meta string
	exchange.Interface
}

func (e *Exchanger) String() string {
	s := fmt.Sprintf("%v › %v", e.Src, e.Dst)
	if len(e.Meta) > 0 {
		s += fmt.Sprintf(" (%v)", e.Meta)
	}
	return s
}

func ParseExchangers(resp *api.Response) []*Exchanger {
	// Group rates by src, dst, and meta.
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
		meta := item.CurrencyGroupAbbr
		m[src][dst][meta] = append(m[src][dst][meta], exchange.Rate{
			Buy:       float64(item.Buy),
			Sell:      float64(item.Sell),
			Threshold: item.Gradation,
		})
	}

	var v []*Exchanger
	for src := range m {
		for dst := range m[src] {
			for meta, rates := range m[src][dst] {
				e := exchange.New(rates...)
				ex := &Exchanger{Src: src, Dst: dst, Meta: meta, Interface: e}
				v = append(v, ex)
			}
		}
	}
	return v
}
