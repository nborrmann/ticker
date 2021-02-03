package watchlist

import (
	"fmt"
	"strconv"
	"strings"
	"ticker/internal/position"
	"ticker/internal/quote"
	. "ticker/internal/ui/util"

	. "ticker/internal/ui/util/text"

	"github.com/novalagung/gubrak/v2"
)

var (
	styleNeutral      = NewStyle("#d4d4d4", "", false)
	styleNeutralBold  = NewStyle("#d4d4d4", "", true)
	styleNeutralFaded = NewStyle("#616161", "", false)
	styleLine         = NewStyle("#3a3a3a", "", false)
	styleTag          = NewStyle("#d4d4d4", "#3a3a3a", false)
	styleTagEnd       = NewStyle("#3a3a3a", "#3a3a3a", false)
	styleGain         = NewStyle("#aae61e", "", false)
	styleLoss         = NewStyle("#FF7940", "", false)
	styleFooter       = NewStyle("#d4d4d4", "", true)
	styleGainFooter   = NewStyle("#aae61e", "", true)
	styleLossFooter   = NewStyle("#FF7940", "", true)
	//stylePricePositive = newStyleFromGradient("#779929", "#C6FF40")
	//stylePriceNegative = newStyleFromGradient("#FF7940", "#994926")
)

const (
	maxPercentChangeColorGradient = 100
)

type Model struct {
	Width                 int
	Quotes                []quote.Quote
	Positions             map[string]position.Position
	Separate              bool
	ExtraInfoExchange     bool
	ExtraInfoFundamentals bool
}

// NewModel returns a model with default values.
func NewModel(separate bool, extraInfoExchange bool, extraInfoFundamentals bool) Model {
	return Model{
		Width:                 80,
		Separate:              separate,
		ExtraInfoExchange:     extraInfoExchange,
		ExtraInfoFundamentals: extraInfoFundamentals,
	}
}

func (m Model) View() string {

	if m.Width < 80 {
		return fmt.Sprintf("Terminal window too narrow to render content\nResize to fix (%d/80)", m.Width)
	}

	horizontalLine := strings.Repeat("─", m.Width)
	quotes := sortQuotes(m.Quotes)
	items := make([]string, 0)
	items = append(items, header(m.Width))
	items = append(items, horizontalLine)
	for _, quote := range quotes {
		items = append(
			items,
			strings.Join(
				[]string{
					item(quote, m.Positions[quote.Symbol], m.Width),
					extraInfoFundamentals(m.ExtraInfoFundamentals, quote, m.Width),
					extraInfoExchange(m.ExtraInfoExchange, quote, m.Width),
				},
				"",
			),
		)
	}

	items = append(items, horizontalLine)
	items = append(items, showTotals(quotes, m.Positions, m.Width))
	return strings.Join(items, separator(m.Separate, m.Width))
}

func separator(isSeparated bool, width int) string {
	if isSeparated {
		return "\n" + Line(
			width,
			false,
			Cell{
				Text: styleLine(strings.Repeat("⎯", width)),
			},
		) + "\n"
	}

	return "\n"
}

func header(width int) string {
	return Line(
		width,
		false,
		Cell{
			Text: styleNeutralBold("SECURITY"),
		},
		Cell{
			Width: 5,
			Text:  "",
			Align: RightAlign,
		},
		Cell{
			Width: 10,
			Text:  styleNeutralBold("Pos"),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  styleNeutralBold("Daily Chg."),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  styleNeutralBold("Total Chg."),
			Align: RightAlign,
		},
	)
}

func item(q quote.Quote, p position.Position, width int) string {

	return Line(
		width,
		false,
		Cell{
			Text: styleNeutral(q.ShortName),
		},
		Cell{
			Width: 5,
			Text:  marketStateText(q),
			Align: RightAlign,
		},
		Cell{
			Width: 10,
			Text:  valueText(p.Value),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  valueChangeText(p.DayChange, p.DayChangePercent, false),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  valueChangeText(p.AbsoluteReturn, p.RelativeReturn, false),
			Align: RightAlign,
		},
	)
}

func showTotals(quotes []quote.Quote, positions map[string]position.Position, width int) string {
	totalReturn := 0.0
	totalCurrent := 0.0
	totalCost := 0.0
	dailyReturn := 0.0
	valuePreviousClose := 0.0
	for _, quote := range quotes {
		if p, ok := positions[quote.Symbol]; ok {
			totalReturn += p.AbsoluteReturn
			totalCurrent += p.Value
			totalCost += p.Cost
			dailyReturn += p.DayChange
			valuePreviousClose += p.ValuePreviousClose
		}
	}
	if totalReturn <= 0.00001 {
		// don't bother displaying if the total basis is near zero
		return ""
	}

	return Line(
		width,
		true,
		Cell{
			Text: styleFooter("TOTAL"),
		},
		Cell{
			Width: 5,
			Text:  "",
			Align: RightAlign,
		},
		Cell{
			Width: 10,
			Text:  styleFooter(ConvertFloatToStringNoDecimals(totalCurrent)),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  valueChangeText(dailyReturn, dailyReturn/valuePreviousClose*100, true),
			Align: RightAlign,
		},
		Cell{
			Width: 22,
			Text:  valueChangeText(totalReturn, totalReturn/totalCost*100, true),
			Align: RightAlign,
		},
	)

}
func extraInfoExchange(show bool, q quote.Quote, width int) string {
	if !show {
		return ""
	}
	return "\n" + Line(
		width,
		false,
		Cell{
			Text: tagText(q.Currency) + " " + tagText(exchangeDelayText(q.ExchangeDelay)) + " " + tagText(q.ExchangeName),
		},
	)
}

func extraInfoFundamentals(show bool, q quote.Quote, width int) string {
	if !show {
		return ""
	}

	return "\n" + Line(
		width,
		false,
		Cell{
			Width: 25,
			Text:  styleNeutralFaded("Prev Close: ") + styleNeutral(ConvertFloatToString(q.RegularMarketPreviousClose, true)),
		},
		Cell{
			Width: 20,
			Text:  styleNeutralFaded("Open: ") + styleNeutral(ConvertFloatToString(q.RegularMarketOpen, true)),
		},
		Cell{
			Text: dayRangeText(q.RegularMarketDayRange),
		},
	)
}

func dayRangeText(dayRange string) string {
	if len(dayRange) <= 0 {
		return ""
	}
	return styleNeutralFaded("Day Range: ") + styleNeutral(dayRange)
}

func exchangeDelayText(delay float64) string {
	if delay <= 0 {
		return "Real-Time"
	}

	return "Delayed " + strconv.FormatFloat(delay, 'f', 0, 64) + "min"
}

func tagText(text string) string {
	return styleTagEnd(" ") + styleTag(text) + styleTagEnd(" ")
}

func marketStateText(q quote.Quote) string {
	if q.IsRegularTradingSession {
		return styleNeutralFaded(" ⦿  ")
	}

	if !q.IsRegularTradingSession && q.IsActive {
		return styleNeutralFaded(" ⦾  ")
	}

	return ""
}

func valueText(value float64) string {
	if value <= 0.0 {
		return ""
	}

	return styleNeutral(ConvertFloatToStringNoDecimals(value))
}

func valueTextBold(value float64) string {
	if value <= 0.0 {
		return ""
	}

	return styleNeutralBold(ConvertFloatToStringNoDecimals(value))
}

func valueChangeText(change float64, changePercent float64, bold bool) string {
	return quoteChangeText(change, changePercent, bold)
}

func quoteChangeText(change float64, changePercent float64, bold bool) string {
	if bold {
		if change == 0.0 {
			return styleNeutralFaded("  " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
		}

		if change > 0.0 {
			return styleGainFooter("↑ " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
		}

		return styleLossFooter("↓ " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
	} else {
		if change == 0.0 {
			return styleNeutralFaded("  " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
		}

		if change > 0.0 {
			return styleGain("↑ " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
		}

		return styleLoss("↓ " + ConvertFloatToString(change, false) + " (" + ConvertFloatToString(changePercent, true) + "%)")
	}
}

// Sort by change percent and keep all inactive quotes at the end
func sortQuotes(q []quote.Quote) []quote.Quote {
	if len(q) <= 0 {
		return q
	}

	activeQuotes, inactiveQuotes, _ := gubrak.
		From(q).
		Partition(func(v quote.Quote) bool {
			return v.IsActive
		}).
		ResultAndError()

	concatQuotes := gubrak.
		From(activeQuotes).
		OrderBy(func(v quote.Quote) float64 {
			return v.ChangePercent
		}, false).
		Concat(inactiveQuotes).
		Result()

	return (concatQuotes).([]quote.Quote)
}
