package gobot

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

type HistoryGrabber struct {
	api *slack.Client
}

func NewHistoryGrabber(api *slack.Client) *HistoryGrabber {
	grabber := &HistoryGrabber{api: api}
	grabber.grabHistory()
	return grabber
}

// TODO where should this method live?
func (h *HistoryGrabber) getUsernames() map[string]string {
	ret := make(map[string]string)
	users, err := h.api.GetUsers()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, user := range users {
		ret[user.ID] = user.Name
	}
	return ret
}

func (h *HistoryGrabber) grabHistory() {
	usernames := h.getUsernames()
	params := &slack.GetConversationHistoryParameters{
		ChannelID: "G6Z2PFJM8",
		// Cursor    string
		// Inclusive bool
		// Latest    string
		// Limit     int
		// Oldest    string
	}
	response, err := h.api.GetConversationHistory(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := len(response.Messages) - 1; i >= 0; i-- {
		message := response.Messages[i]
		fmt.Printf("<%v> %v\n", usernames[message.User], message.Text)
	}
	// fmt.Printf("%+v\n", response)
	// prettyPrint(response)
}

func prettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
