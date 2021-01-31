package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_BOT_TOKEN")
	// teamID := os.Getenv("SLACK_BOT_TEAM_ID")
	channelID := os.Getenv("SLACK_BOT_CHANNEL_ID")
	api := slack.New(token)

	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText("Hello", false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Println(err)
	}
}
