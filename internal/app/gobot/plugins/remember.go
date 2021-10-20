package gobot

import (
	"errors"
	"strings"

	"github.com/sargon2/gobot/internal/app/gobot"
)

type Remember struct {
	hub *gobot.Hub
	db  *gobot.Database
}

type RememberRow struct {
	Key      string
	Username string
	Value    string
}

func NewRemember(hub *gobot.Hub, db *gobot.Database) *Remember {
	ret := &Remember{
		hub: hub,
		db:  db,
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
	if ok := p.db.Put("remember", &RememberRow{Key: key, Username: source.Username, Value: value}); ok {
		p.hub.Message(source, "Okay, "+key+" == "+value)
		return
	}
	p.hub.Message(source, "Oops, failed to remember")
}

func (p *Remember) handleWhatis(source *gobot.MessageSource, message string) {
	message = RemoveHook(message)
	if message == "" {
		p.hub.Message(source, "Usage: !whatis <key>")
		return
	}
	item := &RememberRow{}
	if ok := p.db.Get("remember", item, message); ok { // TODO searching
		p.hub.Message(source, item.Username+" taught me that "+item.Key+" == "+item.Value)
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
	item := &RememberRow{}
	if ok := p.db.Get("remember", item, message); ok {
		deleted := p.db.Delete("remember", item.Key)
		if deleted {
			p.hub.Message(source, "Okay, forgot that "+item.Key+" == "+item.Value)
			return
		}
		p.hub.Message(source, "Oops, failed to delete")
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
