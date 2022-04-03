package gobot

import (
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

func (h *HistoryGrabber) grabHistory() {
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
	fmt.Println(response)
}
