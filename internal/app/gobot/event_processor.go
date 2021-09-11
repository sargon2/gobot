package gobot

import "github.com/slack-go/slack/slackevents"

type EventProcessor interface {
	StartProcessingEvents()
	RegisterMessageCallback(func(*slackevents.MessageEvent))
	Message(*MessageSource, string)
}
