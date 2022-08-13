package gobot

import (
	"errors"
	"strings"

	"github.com/slack-go/slack"
)

func NewSlackClient(slackBotToken *SlackBotToken) (*slack.Client, error) {
	tokenString := string(*slackBotToken)
	if !strings.HasPrefix(tokenString, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	return slack.New(tokenString), nil
}
