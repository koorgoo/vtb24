package chat

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/koorgoo/telegram"
	"github.com/koorgoo/vtb24/bank"
)

func MakeMessage(n float64, ex []bank.Ex) (text string, mode telegram.ParseMode, err error) {
	var buf bytes.Buffer
	var buy, sell float64

	for _, e := range ex {
		buy, err = e.Buy(n)
		if err != nil {
			return
		}
		sell, err = e.Sell(n)
		if err != nil {
			return
		}

		fmt.Fprintf(&buf, "__%s__\n", e)
		fmt.Fprintf(&buf, "**%v** %v - **%v** (покупка) **%v** (продажа) %v\n\n",
			n, e.Src(), FormatValue(buy), FormatValue(sell), e.Dst())
	}

	return buf.String(), telegram.ModeMarkdown, nil
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
