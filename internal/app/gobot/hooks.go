package gobot

import (
	"sort"
	"strings"
)

type Hooks struct {
	hub Hub
}

func NewHooks(hub Hub) *Hooks {
	ret := &Hooks{
		hub: hub,
	}
	hub.RegisterBangHandler("hooks", ret.handleMessage)
	return ret
}

func (h *Hooks) handleMessage(source *MessageSource, msg string) {
	ret := make([]string, 0)
	for cmd, _ := range h.hub.GetBangHandlers() {
		if cmd != "hooks" {
			ret = append(ret, cmd)
		}
	}
	sort.Strings(ret)
	h.hub.Message(source, strings.Join(ret, ", "))
}
