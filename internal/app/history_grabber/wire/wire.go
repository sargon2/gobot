//go:build wireinject

package gobot

import (
	"github.com/google/wire"
	history_grabber "github.com/sargon2/gobot/internal/app/history_grabber"
)

func Begin() {
	WireHistoryGrabber()
}

func WireHistoryGrabber() *history_grabber.HistoryGrabber {
	wire.Build(
		history_grabber.NewHistoryGrabber,
	)
	return &history_grabber.HistoryGrabber{} // Will be magically replaced by wire.
}
