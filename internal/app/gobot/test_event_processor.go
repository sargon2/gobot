package gobot

import (
	"github.com/slack-go/slack/slackevents"
)

type TestEventProcessor struct {
	// TODO this class has nothing to do with slack, so we shouldn't have a reference to slackevents.
	// Maybe the fix is to make our own event type that just has Text and Channel.
	messageCallback func(*slackevents.MessageEvent)
	response        string
}

func NewTestEventProcessor() *TestEventProcessor {
	ret := &TestEventProcessor{}

	return ret
}

func (s *TestEventProcessor) GetResponseFor(input string) string {
	ev := &slackevents.MessageEvent{
		Text: input,
	}
	s.messageCallback(ev)

	return s.response
}

func (s *TestEventProcessor) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	s.messageCallback = cb
}

func (s *TestEventProcessor) Message(source *MessageSource, m string) {
	s.response = m
}

func (s *TestEventProcessor) StartProcessingEvents() error {
	return nil
}
