package gobot

import (
	"fmt"
	"regexp"

	equity "github.com/piquette/finance-go/equity"

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
	stocks := StockSplit(message)
	if len(stocks) == 0 {
		p.hub.Message(source, "Need one or more symbols")
		return
	}
	totalMsg := "```\n"
	for _, stock := range stocks {
		q, err := equity.Get(stock)
		if q != nil {
			// fmt.Printf("%+v\n", q)
			msg := fmt.Sprintf("%50s: %s %s %s%%\n", q.ShortName+" ("+q.Symbol+")", FloatFormat(q.RegularMarketPrice), FloatFormat(q.RegularMarketChange), FloatFormat(q.RegularMarketChangePercent))
			totalMsg += msg
			// TODO set left-width (and right?) based on max room needed for display?
			// TODO RegularMarketChangePercent, RegularMarketChange
		} else {
			if err == nil {
				totalMsg += stock + " not found\n"
			} else {
				totalMsg += "Error getting " + stock + ", err was " + string(err.Error()) + "\n"
			}
		}
	}
	totalMsg += "```"
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
	result := []rune(fmt.Sprintf("%12.3f", f))
	for i := len(result) - 1; result[i] == '0' || result[i] == '.'; i-- {
		result[i] = ' '
	}
	return string(result)
}
