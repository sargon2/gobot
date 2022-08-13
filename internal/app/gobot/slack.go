package gobot

import (
	"errors"
	"strings"

	"github.com/slack-go/slack"
)

func NewSlackClient() (*slack.Client, error) {
	tokenString := *GetSecret("slack_bot_token")
	if !strings.HasPrefix(tokenString, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	return slack.New(tokenString), nil
}
