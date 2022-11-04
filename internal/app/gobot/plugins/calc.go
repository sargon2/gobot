package gobot

import (
	"fmt"

	wolfram "github.com/Krognol/go-wolfram"
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Calc struct {
	hub    *gobot.Hub
	client *wolfram.Client
}

func NewCalc(hub *gobot.Hub) (*Calc, error) {
	ret := &Calc{
		hub:    hub,
		client: &wolfram.Client{AppID: *gobot.GetSecret("wolfram_alpha_key")},
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
	result = "Wolfram Alpha says: " + result
	fmt.Println("calc returning result: " + result)
	c.hub.Message(source, result)
}
