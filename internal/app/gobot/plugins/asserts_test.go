package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wire "github.com/sargon2/gobot/internal/app/gobot/wire"
)

func RunGobotCommand(input string) string {
	testEventProcessor := wire.GetTestEventProcessor()

	return testEventProcessor.GetResponseFor(input)
}

func AssertGobotResponseContains(t *testing.T, input string, expectedOutput string) {
	assert.Contains(t, RunGobotCommand(input), expectedOutput)
}

func AssertGobotResponseIs(t *testing.T, input string, expectedOutput string) {
	assert.Equal(t, expectedOutput, RunGobotCommand(input))
}
