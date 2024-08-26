package gobot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sargon2/gobot/internal/app/gobot"
	"github.com/tidwall/gjson"
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
	// Polymarket
	ret += "Polymarket: presidential-election-winner-2024 (> .035)\n"
	ret += getPolymarket()

	ret += "\n"

	// Predictit
	ret += "Predictit: " + response.Name + " (> .035)\n"
	for _, contract := range response.Contracts {
		if contract.LastTradePrice > .035 {
			ret += fmt.Sprintf("%s %v\n", contract.Name, contract.LastTradePrice)
		}
	}

	// 538
	fivethirtyeight, err := getFiveThirtyEight()
	if err != nil {
		ret += "Error getting 538: " + err.Error()
	} else {
		ret += "\n538: " + fivethirtyeight
	}

	ret += "\n\n"

	// Nate silver
	ret += "Nate Silver polls: " + getNateSilver()
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
	suspended := ""
	cly.OnHTML("div#forecast-suspended-box", func(e *colly.HTMLElement) {
		suspended = e.Text
	})
	cly.OnHTML("div.odds-text-large.mb-10", func(e *colly.HTMLElement) {
		content = strings.Join(e.ChildTexts("div"), " ")
	})
	cly.Visit("https://projects.fivethirtyeight.com/2024-election-forecast/")

	if suspended != "" {
		return suspended, nil
	}
	if content != "" {
		return content, nil
	}
	return "", fmt.Errorf("Div not found")
}

func getPolymarket() string {
	ret := ""
	url := "https://gamma-api.polymarket.com/events?slug=presidential-election-winner-2024"
	contents, err := getURLContents(url)
	if err != nil {
		return "Error getting polymarket: " + err.Error()
	}

	result := gjson.Get(string(contents), "#.markets.#.{groupItemTitle,outcomePrices}")

	type marketData struct {
		title string
		price float64
	}

	var markets []marketData

	for _, outer_item := range result.Array() {
		for _, item := range outer_item.Array() {
			m := item.Map()
			title := m["groupItemTitle"].String()
			price := gjson.Get(m["outcomePrices"].String(), "0").Float()
			if price > .035 {
				markets = append(markets, marketData{title: title, price: price})
			}
		}
	}

	// Sort markets by price
	sort.Slice(markets, func(i, j int) bool {
		return markets[i].price > markets[j].price
	})

	for _, market := range markets {
		ret += market.title
		ret += " "
		ret += fmt.Sprintf("%.2f", market.price)
		ret += "\n"
	}

	return ret
}

func getNateSilver() string {
	// https://www.natesilver.net/p/nate-silver-2024-president-election-polls-model
	// https://datawrapper.dwcdn.net/wB0Zh/16/
	// https://static.dwcdn.net/data/wB0Zh.csv

	ret := ""

	contents, err := getURLContents("https://static.dwcdn.net/data/wB0Zh.csv")
	if err != nil {
		return "Error getting Nate Silver csv: " + err.Error()
	}

	// Read CSV data
	reader := csv.NewReader(strings.NewReader(string(contents)))
	records, err := reader.ReadAll()
	if err != nil {
		return "Error reading Nate Silver csv: " + err.Error()
	}

	// Get the header row and the last data row
	headers := records[0]
	latestRow := records[len(records)-1]

	type resultData struct {
		name string
		val  float64
	}

	var results []resultData

	// Create a map to store the latest poll results
	for i, header := range headers {
		if header == "state" {
			// state is always "National"
			continue
		}
		if header == "modeldate" {
			ret += "as of " + latestRow[i]
		} else if !strings.HasSuffix(header, "_poll") {
			if latestRow[i] == "" {
				continue
			}
			valf, err := strconv.ParseFloat(latestRow[i], 64)
			if err != nil {
				return "Error parsing Nate Silver csv: " + err.Error()
			}
			results = append(results, resultData{name: header, val: valf})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].val > results[j].val
	})

	for _, candidate := range results {
		ret += "\n"
		ret += fmt.Sprintf("%s: %.2f", candidate.name, candidate.val)
	}

	return ret

}
