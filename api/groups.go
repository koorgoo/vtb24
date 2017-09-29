package api

import "fmt"

// Typical currency group abbreviations.
const (
	GroupCash        = "cash"
	GroupCashDesk    = "cash-desk"
	GroupCentralBank = "central-bank"
	// GroupCentralBankJur = "central-bank-jur"
	// GroupOfficeCash     = "pp_curcur_office_cash"
	// GroupOfficeCashless = "pp_curcur_office_cashless"
	// GroupSpecKassaCash  = "pp_curcur_speckassa_cash"
	GroupTele = "tele"
	// GroupW4   = "w4"
)

// GroupText returns a text in Russian for provided currency group.
func GroupText(group string) string {
	if text, ok := groupText[group]; ok {
		return text
	}
	return fmt.Sprintf("unknown group %q", group)
}

var groupText = map[string]string{
	"cash":         "в офисе, наличные",
	"cash-desk":    "в спецкассе",
	"central-bank": "в офисе, безналичные",
	// "central-bank-jur":          "",
	// "pp_curcur_office_cash":     "",
	// "pp_curcur_office_cashless": "",
	// "pp_curcur_speckassa_cash":  "",
	"tele": "в ВТБ24 - онлайн",
	// "w4":   "",
}
