package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wire "github.com/sargon2/gobot/internal/app/gobot/wire"
)

func assertGobotResponseContains(t *testing.T, input string, expectedOutput string) {
	testEventProcessor := wire.GetTestEventProcessor()

	output := testEventProcessor.GetResponseFor(input)

	assert.Contains(t, output, expectedOutput)
}

func TestPing(t *testing.T) {
	assertGobotResponseContains(t, "!ping", "pong")
}
