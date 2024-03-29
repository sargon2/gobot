package gobot

import (
	"fmt"

	"github.com/djdv/go-swatch"
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Time struct {
	hub *gobot.Hub
}

func NewTime(hub *gobot.Hub) *Time {
	ret := &Time{
		hub: hub,
	}
	hub.RegisterBangHandler("time", ret.handleMessage)
	return ret
}

func (t *Time) handleMessage(source *gobot.MessageSource, message string) {
	t.hub.Message(source, fmt.Sprintf("Swatch Internet Time: @%f", swatch.New().PreciseBeats()))
}
