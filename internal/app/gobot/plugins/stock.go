package gobot

import (
	"fmt"
	"regexp"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/piquette/finance-go/quote"

	"github.com/sargon2/gobot/internal/app/gobot"
)

type Stock struct {
	hub *gobot.Hub
}

func NewStock(hub *gobot.Hub) *Stock {
	ret := &Stock{
		hub: hub,
	}
	hub.RegisterBangHandler("stock", ret.handleMessage)
	return ret
}

func (p *Stock) handleMessage(source *gobot.MessageSource, message string) {
	tw := table.NewWriter()
	tw.Style().Options.DrawBorder = false
	tw.Style().Options.SeparateRows = false
	tw.Style().Options.SeparateColumns = false
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft},
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
	})

	stocks := StockSplit(message)
	if len(stocks) == 0 {
		p.hub.Message(source, "Need one or more symbols")
		return
	}
	totalMsg := ""
	for _, stock := range stocks {
		q, err := quote.Get(stock)
		if q != nil {
			tw.AppendRow(table.Row{
				q.ShortName,
				q.Symbol,
				FloatFormat(q.RegularMarketPrice),
				FloatFormat(q.RegularMarketChange),
				FloatFormat(q.RegularMarketChangePercent) + "%",
			})
		} else {
			if err == nil {
				totalMsg += stock + " not found\n"
			} else {
				totalMsg += "Error getting " + stock + ", err was " + string(err.Error()) + "\n"
			}
		}
	}
	totalMsg = "```\n" + tw.Render() + "\n" + totalMsg + "```\n"
	p.hub.Message(source, totalMsg)
}

func StockSplit(message string) []string {
	ret := []string{}
	r := regexp.MustCompile("[\\s,;/\\\\\"'`]+")
	for _, item := range r.Split(message, -1) {
		if len(item) > 0 {
			ret = append(ret, item)
		}
	}
	return ret
}

func FloatFormat(f float64) string {
	result := []rune(fmt.Sprintf("%.3f", f))
	for i := len(result) - 1; result[i] == '0' || result[i] == '.'; i-- {
		result[i] = ' '
	}
	return string(result)
}
