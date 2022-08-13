//go:build wireinject

package gobot

import (
	"fmt"

	"github.com/google/wire"
	"github.com/sargon2/gobot/internal/app/gobot"
	history_grabber "github.com/sargon2/gobot/internal/app/history_grabber"
)

func Begin() {
	_, err := WireHistoryGrabber()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func WireHistoryGrabber() (*history_grabber.HistoryGrabber, error) {
	wire.Build(
		gobot.NewSlackClient,
		history_grabber.NewHistoryGrabber,
	)
	return &history_grabber.HistoryGrabber{}, nil // Will be magically replaced by wire.
}
