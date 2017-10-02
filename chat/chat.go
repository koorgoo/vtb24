package chat

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/koorgoo/telegram"
	"github.com/koorgoo/vtb24/api"
	"github.com/koorgoo/vtb24/bank"
)

func MakeMessage(n float64, ex []bank.Ex) (text string, mode telegram.ParseMode) {
	var buf bytes.Buffer
	for _, e := range ex {
		s1, ok1 := formatOp(n, e)
		s2, ok2 := formatOp(n, bank.Invert(e))

		if !ok1 && !ok2 {
			continue
		}

		fmt.Fprintln(&buf, FormatEx(e))
		if ok1 {
			fmt.Fprintln(&buf, s1)
		}
		if ok2 {
			fmt.Fprintln(&buf, s2)
		}

		// Line break.
		fmt.Fprintln(&buf)
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

func FormatEx(ex bank.Ex) string {
	s := api.GroupText(ex.Group())
	return fmt.Sprintf("_%s_", s)
}
