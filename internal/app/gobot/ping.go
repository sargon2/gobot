package gobot

type Ping struct {
	hub Hub
}

func NewPing(hub Hub, hooks *Hooks) *Ping {
	ret := &Ping{
		hub: hub,
	}
	hooks.RegisterBangHandler("ping", ret.handleMessage)
	return ret
}

func (p *Ping) handleMessage(source *MessageSource, message string) {
	p.hub.Message(source, "pong")
}
