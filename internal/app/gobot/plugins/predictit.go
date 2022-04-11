package gobot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sargon2/gobot/internal/app/gobot"
)

type Predictit struct {
	hub *gobot.Hub
}

type PredictitResponse struct {
	Contracts []PredictitResponseContract `json:contracts`
}

type PredictitResponseContract struct {
	Name           string  `json:name`
	LastClosePrice float64 `json:lastClosePrice`
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
	for _, contract := range response.Contracts {
		ret += fmt.Sprintf("%s %v\n", contract.Name, contract.LastClosePrice)
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
