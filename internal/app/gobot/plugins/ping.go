package gobot

import (
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Ping struct {
	hub *gobot.Hub
}

func NewPing(hub *gobot.Hub) *Ping {
	ret := &Ping{
		hub: hub,
	}
	hub.RegisterBangHandler("ping", ret.handleMessage)
	return ret
}

func (p *Ping) handleMessage(source *gobot.MessageSource, message string) {
	p.hub.Message(source, "pong")
}
