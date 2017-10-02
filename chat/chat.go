package chat

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/koorgoo/telegram"
	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
)

func MakeMessage(n float64, ex []bank.Ex, groups []string) (text string, mode telegram.ParseMode) {
	m := map[string][]bank.Ex{}
	for _, e := range ex {
		m[e.Group()] = append(m[e.Group()], e)
	}

	var buf bytes.Buffer
	var hasGroups = true

	for _, group := range groups {
		var writeGroup sync.Once

		for _, e := range m[group] {
			s, ok := formatOp(n, e)
			si, oki := formatOp(n, bank.Invert(e))

			if !ok && !oki {
				continue
			}

			writeGroup.Do(func() {
				fmt.Fprintf(&buf, formatGroup(e.Group(), !hasGroups))
				hasGroups = true
			})

			if ok {
				fmt.Fprintln(&buf, s)
			}
			if oki {
				fmt.Fprintln(&buf, si)
			}

			// Line break between ops.
			fmt.Fprintln(&buf)
		}
	}
	return buf.String(), telegram.ModeMarkdown
}

func formatOp(n float64, e bank.Ex) (s string, ok bool) {
	buy, err := e.Buy(n)
	if err != nil {
		return
	}
	sell, err := e.Sell(n)
	if err != nil {
		return
	}
	s = fmt.Sprintf("*%v* %v - *%v* (покупка) *%v* (продажа) %v",
		n, e.Src(), FormatValue(buy), FormatValue(sell), e.Dst())
	return s, true
}

func FormatValue(v float64) (s string) {
	if n := int64(v); float64(n) == v {
		s = strconv.FormatInt(n, 10)
	} else {
		s = big.NewFloat(v).Text('f', 2)
	}
	if strings.HasSuffix(s, ".00") {
		s = s[:len(s)-3]
	}
	return
}

func formatGroup(group string, isFirst bool) string {
	var suffix string
	if !isFirst {
		suffix = "\n"
	}
	text := api.GroupText(group)
	return fmt.Sprintf("%s_%s_\n\n", suffix, text)
}
