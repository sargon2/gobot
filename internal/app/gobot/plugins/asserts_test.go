package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wire "github.com/sargon2/gobot/internal/app/gobot/wire"
)

func AssertGobotResponseContains(t *testing.T, input string, expectedOutput string) {
	testEventProcessor := wire.GetTestEventProcessor()

	output := testEventProcessor.GetResponseFor(input)

	assert.Contains(t, output, expectedOutput)
}

func AssertGobotResponseIs(t *testing.T, input string, expectedOutput string) {
	testEventProcessor := wire.GetTestEventProcessor()

	output := testEventProcessor.GetResponseFor(input)

	assert.Equal(t, expectedOutput, output)
}
