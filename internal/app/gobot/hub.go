package gobot

import "github.com/slack-go/slack/slackevents"

// Hub is an interface layer between plugins and functionality providers.
// This helps to avoid dependency loops between the types that provide functionality
// to plugins.
type Hub struct {
	eventProcessor EventProcessor
	bangManager    *BangManager
}

func NewHub(eventProcessor EventProcessor, bangManager *BangManager) *Hub {
	return &Hub{
		eventProcessor: eventProcessor,
		bangManager:    bangManager,
	}
}

func (h *Hub) RegisterBangHandler(cmd string, handler func(*MessageSource, string)) {
	h.bangManager.RegisterBangHandler(cmd, handler)
}

func (h *Hub) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	h.eventProcessor.RegisterMessageCallback(cb)
}

func (h *Hub) Message(source *MessageSource, m string) {
	h.eventProcessor.Message(source, m)
}
