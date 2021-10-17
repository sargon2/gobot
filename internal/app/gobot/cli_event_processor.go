package gobot

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack/slackevents"
)

type CliEventProcessor struct {
	// TODO CLI has nothing to do with slack, so we shouldn't have a reference to slackevents.
	// Maybe the fix is to make our own event type that just has Text and Channel.
	messageCallback func(*slackevents.MessageEvent) // TODO this should be a list of callbacks
}

func NewCliEventProcessor() *CliEventProcessor {
	ret := &CliEventProcessor{}

	return ret
}

func (s *CliEventProcessor) GetUsername(userID string) string {
	return userID
}

func (s *CliEventProcessor) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	s.messageCallback = cb
}

func (s *CliEventProcessor) Message(source *MessageSource, m string) {
	fmt.Println(m)
}

func (s *CliEventProcessor) StartProcessingEvents() error {
	// Since we're CLI, we only have one event to process, and it's in the cli args.
	command := strings.Join(os.Args[1:], " ")
	ev := &slackevents.MessageEvent{
		Text: command,
		User: "cli",
	}
	s.messageCallback(ev)
	return nil
}
