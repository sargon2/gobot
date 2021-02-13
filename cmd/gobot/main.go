package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/sargon2/gobot/internal/app/gobot"
)

// This is what tells wire which hooks to use
type Hooks struct {
	Hub      gobot.Hub
	Ping     *gobot.Ping
	Roll     *gobot.Roll
	Twilight *gobot.Sun
}

// This tells wire what type providers we have.  Ideally it would auto-detect them somehow but it doesn't support that today.
func WireHooks() (*Hooks, error) {
	wire.Build(
		gobot.NewSlackHub,
		wire.Bind(new(gobot.Hub), new(*gobot.SlackHub)),
		wire.Struct(new(Hooks), "*"),
		gobot.NewPing,
		gobot.NewRoll,
		gobot.NewSun,
	)
	return &Hooks{}, nil // Will be magically replaced by wire.
}

func main() {
	hooks, err := WireHooks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hooks.Hub.StartEventLoop()
}
