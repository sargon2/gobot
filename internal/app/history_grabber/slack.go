package gobot

import (
	"errors"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func NewSlackClient() (*slack.Client, error) {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}
	return slack.New(botToken), nil
}
