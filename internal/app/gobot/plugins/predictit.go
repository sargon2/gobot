package gobot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Predictit struct {
	hub *gobot.Hub
}

type PredictitResponse struct {
	Name      string                      `json:name`
	Contracts []PredictitResponseContract `json:contracts`
}

type PredictitResponseContract struct {
	Name           string  `json:name`
	LastTradePrice float64 `json:lastTradePrice`
}

func NewPredictit(hub *gobot.Hub) *Predictit {
	ret := &Predictit{
		hub: hub,
	}
	hub.RegisterBangHandler("predictit", ret.handleMessage)
	return ret
}

func (p *Predictit) handleMessage(source *gobot.MessageSource, message string) {
	contents, err := getURLContents("https://www.predictit.org/api/marketdata/markets/7456")
	if err != nil {
		p.hub.Message(source, "Error! "+err.Error())
		return
	}
	var response PredictitResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		p.hub.Message(source, "Error unmarshalling! "+err.Error())
		return
	}

	ret := "```"
	ret += response.Name + "\n"
	for _, contract := range response.Contracts {
		ret += fmt.Sprintf("%s %v\n", contract.Name, contract.LastTradePrice)
	}
	fivethirtyeight, err := getFiveThirtyEight()
	if err != nil {
		p.hub.Message(source, "Error getting 538: "+err.Error())
	} else {
		ret += "\n538 says: " + fivethirtyeight
	}
	ret += "```"

	p.hub.Message(source, ret)
}

// TODO where should this live?
func getURLContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %v", err)
	}

	return data, nil
}

func getFiveThirtyEight() (string, error) {
	// https://projects.fivethirtyeight.com/2024-election-forecast/
	// https://github.com/gocolly/colly
	// https://stackoverflow.com/questions/65971880/scrape-only-a-certain-div-using-gocolly

	cly := colly.NewCollector(
		colly.AllowedDomains("projects.fivethirtyeight.com"),
	)
	// cly.OnHTML("body", func(e *colly.HTMLElement) {
	// 	link := e.Attr("div")
	// 	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	// 	// cly.Visit(e.Request.AbsoluteURL(link))
	// })
	content := ""
	cly.OnHTML("div.odds-text-large.mb-10", func(e *colly.HTMLElement) {
		content = strings.Join(e.ChildTexts("div"), " ")
	})
	cly.Visit("https://projects.fivethirtyeight.com/2024-election-forecast/")

	if content != "" {
		return content, nil
	}
	return "", fmt.Errorf("Div not found")
}
