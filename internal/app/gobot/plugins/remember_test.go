package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

// TODO:
// - Remember who set each key
// - How should we test with dynamodb?
// - Should we allow multiple people to remember the same key?
// - Make remember replace so you don't need to forget to change it

func TestRemember(t *testing.T) {
	AssertGobotResponseIs(t, "!remember asdf == jkl", "Okay, asdf == jkl")
	AssertGobotResponseIs(t, "!whatis asdf", "jkl")
	AssertGobotResponseIs(t, "!forget asdf", "Okay, forgot that asdf == jkl")
	AssertGobotResponseIs(t, "!whatis asdf", "asdf not found")
}

func TestUsages(t *testing.T) {
	AssertGobotResponseIs(t, "!remember", "Usage: !remember <key> == <value>")
	AssertGobotResponseIs(t, "!remember asdf", "Usage: !remember <key> == <value>")
	AssertGobotResponseIs(t, "!remember asdf ==", "Usage: !remember <key> == <value>")
	AssertGobotResponseIs(t, "!remember == jkl", "Usage: !remember <key> == <value>")
	AssertGobotResponseIs(t, "!whatis", "Usage: !whatis <key>")
	AssertGobotResponseIs(t, "!forget", "Usage: !forget <key>")
}

func TestNotFound(t *testing.T) {
	AssertGobotResponseIs(t, "!whatis notfound", "notfound not found")
	AssertGobotResponseIs(t, "!forget notfound", "notfound not found")
}

func assertParseRememberMessage(t *testing.T, input string, expectedKey string, expectedValue string) {
	key, value, err := plugins.ParseRememberMessage(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedKey, key)
	assert.Equal(t, expectedValue, value)
}

func assertParseRememberMessageError(t *testing.T, input string) {
	key, value, err := plugins.ParseRememberMessage(input)
	assert.Equal(t, "", key)
	assert.Equal(t, "", value)
	assert.NotNil(t, err)
}
func TestRememberParse(t *testing.T) {
	assertParseRememberMessage(t, "asdf == jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf==jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf== jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf ==jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf    ==      jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf==      jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf     ==jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "    asdf == jkl     ", "asdf", "jkl")
	assertParseRememberMessage(t, "!remember asdf == jkl", "asdf", "jkl")
	assertParseRememberMessage(t, "asdf == jkl == foo", "asdf", "jkl == foo")
	assertParseRememberMessageError(t, "asdf jkl")
	assertParseRememberMessageError(t, "asdf")
}
