package gobot

import (
	"errors"
	"strings"

	"github.com/sargon2/gobot/internal/app/gobot"
)

var storage map[string]string = make(map[string]string)
var nickStorage map[string]string = make(map[string]string)

type Remember struct {
	hub *gobot.Hub
}

func NewRemember(hub *gobot.Hub) *Remember {
	ret := &Remember{
		hub: hub,
	}
	hub.RegisterBangHandler("remember", ret.handleRemember)
	hub.RegisterBangHandler("whatis", ret.handleWhatis)
	hub.RegisterBangHandler("forget", ret.handleForget)
	return ret
}

func (p *Remember) handleRemember(source *gobot.MessageSource, message string) {
	key, value, err := ParseRememberMessage(message)
	if err != nil {
		p.hub.Message(source, err.Error())
		return
	}
	storage[key] = value
	nickStorage[key] = source.Username
	p.hub.Message(source, "Okay, "+key+" == "+value)
}

func (p *Remember) handleWhatis(source *gobot.MessageSource, message string) {
	message = RemoveHook(message)
	if message == "" {
		p.hub.Message(source, "Usage: !whatis <key>")
		return
	}
	if value, ok := storage[message]; ok {
		nick := nickStorage[message]
		p.hub.Message(source, nick+" taught me that "+message+" == "+value)
		return
	}
	p.hub.Message(source, message+" not found")
}

func (p *Remember) handleForget(source *gobot.MessageSource, message string) {
	message = RemoveHook(message)
	if message == "" {
		p.hub.Message(source, "Usage: !forget <key>")
		return
	}
	if value, ok := storage[message]; ok {
		delete(storage, message)
		p.hub.Message(source, "Okay, forgot that "+message+" == "+value)
		return
	}
	p.hub.Message(source, message+" not found")
}

func ParseRememberMessage(message string) (key, value string, err error) {
	message = RemoveHook(message)
	parts := strings.SplitN(message, "==", 2)
	if len(parts) == 1 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return "", "", errors.New("Usage: !remember <key> == <value>")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}
