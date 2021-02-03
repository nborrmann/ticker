package util

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ConvertFloatToString(f float64, decimals bool) string {
	if !decimals {
		return ConvertFloatToStringNoDecimals(f)
	}

	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", f)
}

func ConvertFloatToStringNoDecimals(f float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.f", f)
}
