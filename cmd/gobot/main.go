package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/sargon2/gobot/internal/app/gobot"
	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

// This is what tells wire which hooks to use
type Hooks struct {
	EventProcessor gobot.EventProcessor

	Ping  *plugins.Ping
	Roll  *plugins.Roll
	Sun   *plugins.Sun
	Time  *plugins.Time
	Hooks *plugins.Hooks
}

// This tells wire what type providers we have.  Ideally it would auto-detect them somehow but it doesn't support that today.
func WireHooks() (*Hooks, error) {
	wire.Build(
		gobot.NewLambdaEventProcessor,
		wire.Bind(new(gobot.EventProcessor), new(*gobot.LambdaEventProcessor)),
		wire.Struct(new(Hooks), "*"),
		gobot.NewLocationFinder,
		gobot.NewBangManager,
		gobot.NewHub,

		plugins.NewHooks,
		plugins.NewPing,
		plugins.NewRoll,
		plugins.NewSun,
		plugins.NewTime,
	)
	return &Hooks{}, nil // Will be magically replaced by wire.
}

func main() {
	hooks, err := WireHooks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hooks.EventProcessor.StartProcessingEvents()
}
