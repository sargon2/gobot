package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	// teamID := "T0EUDFVLM"    // solsar
	channelID := "G6Z2PFJM8" // #vatcave
	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText("Hello", false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Println(err)
	}
}
