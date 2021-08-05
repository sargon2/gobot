package gobot

type Hub interface {
	StartEventLoop()
	RegisterBangHandler(string, func(*MessageSource, string))
	Message(*MessageSource, string)
}
