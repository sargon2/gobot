package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

func assertRemoveHook(t *testing.T, input string, expectedOutput string) {
	result := plugins.RemoveHook(input)

	assert.Equal(t, expectedOutput, result)
}
func TestRemoveHook(t *testing.T) {
	assertRemoveHook(t, "", "")
	assertRemoveHook(t, "asdf", "asdf")
	assertRemoveHook(t, "!hook asdf", "asdf")
	assertRemoveHook(t, "!hook", "")
	assertRemoveHook(t, "!differenthook asdf", "asdf")
	assertRemoveHook(t, "!hook1 !hook2 asdf", "!hook2 asdf")
	assertRemoveHook(t, "!hook asdf jkl", "asdf jkl")
}
