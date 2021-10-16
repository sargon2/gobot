package gobot

import (
	"sort"
	"strings"

	"github.com/sargon2/gobot/internal/app/gobot"
)

type Hooks struct {
	hub         *gobot.Hub
	bangManager *gobot.BangManager
}

func NewHooks(hub *gobot.Hub, bangManager *gobot.BangManager) *Hooks {
	ret := &Hooks{
		hub:         hub,
		bangManager: bangManager,
	}
	hub.RegisterBangHandler("hooks", ret.handleMessage)
	return ret
}

func (h *Hooks) handleMessage(source *gobot.MessageSource, msg string) {
	bangHandlers := h.bangManager.GetBangHandlers()

	ret := make([]string, 0)
	for _, cmd := range bangHandlers {
		if cmd != "hooks" {
			ret = append(ret, cmd)
		}
	}
	sort.Strings(ret)
	h.hub.Message(source, strings.Join(ret, ", "))
}

// Remove the hook from the given message, if it has one.
func RemoveHook(message string) string {
	if message == "" {
		return ""
	}
	if message[0] == '!' {
		parts := strings.SplitN(message, " ", 2)
		if len(parts) == 1 {
			return ""
		} else if len(parts) == 2 {
			return parts[1]
		}
	}
	return message
}
