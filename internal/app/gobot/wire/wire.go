//go:build wireinject

package gobot

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/sargon2/gobot/internal/app/gobot"
	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

var testMode bool = false

// This is what tells wire which hooks to use
type Hooks struct {
	EventProcessor gobot.EventProcessor

	Ping      *plugins.Ping
	Roll      *plugins.Roll
	Sun       *plugins.Sun
	Time      *plugins.Time
	Hooks     *plugins.Hooks
	Stock     *plugins.Stock
	Remember  *plugins.Remember
	Calc      *plugins.Calc
	Predictit *plugins.Predictit
}

// This tells wire what type providers we have.  Ideally it would auto-detect them somehow but it doesn't support that today.
func WireHooks() (*Hooks, error) {
	wire.Build(
		// Event processors
		gobot.NewTestEventProcessor,
		gobot.NewLambdaEventProcessor,
		gobot.NewCliEventProcessor,
		NewEventProcessor,

		// Hooks
		wire.Struct(new(Hooks), "*"),

		// Supporting types
		gobot.NewSlackClient,
		gobot.NewLocationFinder,
		gobot.NewBangManager,
		gobot.NewHub,
		gobot.NewDatabase,

		// Plugins
		plugins.NewHooks,
		plugins.NewPing,
		plugins.NewRoll,
		plugins.NewSun,
		plugins.NewTime,
		plugins.NewStock,
		plugins.NewRemember,
		plugins.NewCalc,
		plugins.NewPredictit,
	)
	return &Hooks{}, nil // Will be magically replaced by wire.
}

func NewEventProcessor(lambda *gobot.LambdaEventProcessor, cli *gobot.CliEventProcessor, test *gobot.TestEventProcessor) gobot.EventProcessor {
	if testMode {
		return test
	}

	if len(os.Args) > 1 {
		return cli
	}

	return lambda
}

func GetTestEventProcessor() *gobot.TestEventProcessor {
	testMode = true
	hooks := Begin()
	return hooks.EventProcessor.(*gobot.TestEventProcessor)
}

func Begin() *Hooks {
	hooks, err := WireHooks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = hooks.EventProcessor.StartProcessingEvents()
	if err != nil {
		fmt.Println(err)
	}
	return hooks
}
