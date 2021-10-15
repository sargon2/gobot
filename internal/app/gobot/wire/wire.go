//+build wireinject

package gobot

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
	Stock *plugins.Stock
}

// This tells wire what type providers we have.  Ideally it would auto-detect them somehow but it doesn't support that today.
func WireHooks() (*Hooks, error) {
	wire.Build(
		gobot.NewLambdaEventProcessor,
		gobot.NewCliEventProcessor,
		NewEventProcessor,
		wire.Struct(new(Hooks), "*"),
		gobot.NewLocationFinder,
		gobot.NewBangManager,
		gobot.NewHub,

		plugins.NewHooks,
		plugins.NewPing,
		plugins.NewRoll,
		plugins.NewSun,
		plugins.NewTime,
		plugins.NewStock,
	)
	return &Hooks{}, nil // Will be magically replaced by wire.
}

// TODO dup'd
func WireHooksForTest() (*Hooks, error) {
	wire.Build(
		gobot.NewTestEventProcessor,
		NewTestEventProcessor,
		wire.Struct(new(Hooks), "*"),
		gobot.NewLocationFinder,
		gobot.NewBangManager,
		gobot.NewHub,

		plugins.NewHooks,
		plugins.NewPing,
		plugins.NewRoll,
		plugins.NewSun,
		plugins.NewTime,
		plugins.NewStock,
	)
	return &Hooks{}, nil // Will be magically replaced by wire.
}

func NewEventProcessor(lambda *gobot.LambdaEventProcessor, cli *gobot.CliEventProcessor) gobot.EventProcessor {
	if len(os.Args) > 1 {
		return cli
	} else {
		return lambda
	}
}

func NewTestEventProcessor(test *gobot.TestEventProcessor) gobot.EventProcessor { // TODO this function shouldn't be needed
	return test
}

func Begin() {
	hooks, err := WireHooks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = hooks.EventProcessor.StartProcessingEvents()
	if err != nil {
		fmt.Println(err)
	}
}

func GetTestEventProcessor() *gobot.TestEventProcessor {
	// TODO c&p'd with Begin
	hooks, err := WireHooksForTest() // TODO this should only be done once per test run...
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = hooks.EventProcessor.StartProcessingEvents()
	if err != nil {
		fmt.Println(err)
	}
	return hooks.EventProcessor.(*gobot.TestEventProcessor)
}
