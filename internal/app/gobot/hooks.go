package gobot

import (
	"github.com/slack-go/slack/slackevents"
	"sort"
	"strings"
)

type Hooks struct {
	hub          Hub
	bangHandlers map[string]func(*MessageSource, string)
}

func NewHooks(hub Hub) *Hooks {
	ret := &Hooks{
		hub:          hub,
		bangHandlers: make(map[string]func(*MessageSource, string)),
	}
	hub.RegisterMessageCallback(ret.handleBangs)
	ret.RegisterBangHandler("hooks", ret.handleMessage)
	return ret
}

func (h *Hooks) RegisterBangHandler(cmd string, handler func(*MessageSource, string)) {
	h.bangHandlers[cmd] = handler
}

func (h *Hooks) handleBangs(event *slackevents.MessageEvent) {
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

func (h *Hooks) handleMessage(source *MessageSource, msg string) {
	ret := make([]string, 0)
	for cmd, _ := range h.bangHandlers {
		if cmd != "hooks" {
			ret = append(ret, cmd)
		}
	}
	sort.Strings(ret)
	h.hub.Message(source, strings.Join(ret, ", "))
}
