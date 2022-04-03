//go:build wireinject

package gobot

import (
	"fmt"

	"github.com/google/wire"
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
		history_grabber.NewSlackClient,
		history_grabber.NewHistoryGrabber,
	)
	return &history_grabber.HistoryGrabber{}, nil // Will be magically replaced by wire.
}
