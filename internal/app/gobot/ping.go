package gobot

type Ping struct {
	hub Hub
}

func NewPing(hub Hub) *Ping {
	ret := &Ping{
		hub: hub,
	}
	hub.RegisterBangHandler("ping", ret.handleMessage)
	return ret
}

func (p *Ping) handleMessage(source *MessageSource, message string) {
	p.hub.Message(source, "pong")
}
