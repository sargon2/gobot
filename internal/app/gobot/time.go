package gobot

import (
	"fmt"

	swatch "github.com/djdv/go-swatch"
)

type Time struct {
	hub Hub
}

func NewTime(hub Hub, hooks *Hooks) *Time {
	ret := &Time{
		hub: hub,
	}
	hooks.RegisterBangHandler("time", ret.handleMessage)
	return ret
}

func (t *Time) handleMessage(source *MessageSource, message string) {
	t.hub.Message(source, fmt.Sprintf("Swatch Internet Time: %s", swatch.Now(swatch.Centi)))
}
