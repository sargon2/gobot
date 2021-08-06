package gobot

import (
	"github.com/slack-go/slack/slackevents"
)

type Hub interface {
	StartEventLoop()
	RegisterMessageCallback(func(*slackevents.MessageEvent))
	Message(*MessageSource, string)
}
