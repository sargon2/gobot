package gobot

import "github.com/slack-go/slack/slackevents"

type EventProcessor interface {
	StartProcessingEvents() error
	RegisterMessageCallback(func(*slackevents.MessageEvent))
	Message(*MessageSource, string)
}