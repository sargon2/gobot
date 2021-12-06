package gobot

import (
	"errors"
	"fmt"
	"os"

	wolfram "github.com/Krognol/go-wolfram"
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Calc struct {
	hub    *gobot.Hub
	client *wolfram.Client
}

func NewCalc(hub *gobot.Hub) (*Calc, error) {
	apiKey := os.Getenv("WOLFRAM_ALPHA_KEY")
	if apiKey == "" {
		return nil, errors.New("WOLFRAM_ALPHA_KEY must be set")
	}

	ret := &Calc{
		hub:    hub,
		client: &wolfram.Client{AppID: apiKey},
	}

	hub.RegisterBangHandler("calc", ret.handleMessage)
	hub.RegisterBangHandler("ask", ret.handleMessage)
	hub.RegisterBangHandler("conv", ret.handleMessage)

	return ret, nil
}

func (c *Calc) handleMessage(source *gobot.MessageSource, message string) {
	fmt.Println("Starting calc handleMessage")
	result, err := c.client.GetShortAnswerQuery(message, wolfram.Imperial, 30)
	if err != nil {
		fmt.Println("calc returning error: " + err.Error())
		c.hub.Message(source, "Got an error: "+err.Error())
		return
	}
	fmt.Println("calc returning result: " + result)
	c.hub.Message(source, result)
}