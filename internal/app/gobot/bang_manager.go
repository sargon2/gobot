package gobot

import (
	"strings"

	"github.com/slack-go/slack/slackevents"
)

type BangManager struct {
	eventProcessor EventProcessor
	bangHandlers   map[string]func(*MessageSource, string)
}

func NewBangManager(eventProcessor EventProcessor) *BangManager {
	ret := &BangManager{
		eventProcessor: eventProcessor,
		bangHandlers:   make(map[string]func(*MessageSource, string)),
	}
	eventProcessor.RegisterMessageCallback(ret.handleBangs)
	return ret
}

func (h *BangManager) RegisterBangHandler(cmd string, handler func(*MessageSource, string)) {
	h.bangHandlers[cmd] = handler
}

func (h *BangManager) handleBangs(event *slackevents.MessageEvent) {
	messageText := event.Text
	channelID := event.Channel
	for cmd, handler := range h.bangHandlers {
		bangCmd := "!" + cmd
		if messageText == bangCmd || strings.HasPrefix(messageText, bangCmd+" ") {
			source := &MessageSource{
				ChannelID: channelID,
			}
			messageText = strings.TrimSpace(messageText[len(bangCmd):])
			handler(source, messageText)
		}

	}
}

func (h *BangManager) GetBangHandlers() []string {
	ret := make([]string, len(h.bangHandlers))
	i := 0
	for cmd := range h.bangHandlers {
		ret[i] = cmd
		i++
	}
	return ret
}
