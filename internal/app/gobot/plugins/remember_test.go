package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

// TODO:
// - How should we test with dynamodb?
// - Should we allow multiple people to remember the same key?
// - Make remember replace so you don't need to forget to change it
// - Fuzzy lookups; !whatis thing should match "the thing"

func TestRemember(t *testing.T) {
	RunGobotCommand("!forget asdf")
	AssertGobotResponseIs(t, "!remember asdf == jkl", "Okay, asdf == jkl")
	AssertGobotResponseIs(t, "!whatis asdf", "tests taught me that asdf == jkl")
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

func TestSearching(t *testing.T) {
	RunGobotCommand("!forget asdf")
	RunGobotCommand("!forget asdf2")
	RunGobotCommand("!forget asdf3")
	AssertGobotResponseIs(t, "!remember asdf == jkl", "Okay, asdf == jkl")
	AssertGobotResponseIs(t, "!remember asdf2 == jkl2", "Okay, asdf2 == jkl2")
	AssertGobotResponseIs(t, "!whatis asdf", "tests taught me that asdf == jkl\n(also asdf2)")
	AssertGobotResponseIs(t, "!whatis as", "tests taught me that asdf == jkl\n(also asdf2)")
	AssertGobotResponseIs(t, "!remember asdf3 == jkl3", "Okay, asdf3 == jkl3")
	AssertGobotResponseIs(t, "!whatis as", "tests taught me that asdf == jkl\n(also asdf2, asdf3)")
	AssertGobotResponseIs(t, "!whatis asdf2", "tests taught me that asdf2 == jkl2")
	AssertGobotResponseIs(t, "!whatis asdf3", "tests taught me that asdf3 == jkl3")
	RunGobotCommand("!forget asdf")
	RunGobotCommand("!forget asdf2")
	RunGobotCommand("!forget asdf3")
}

func TestReplace(t *testing.T) {
	RunGobotCommand("!forget asdf")
	AssertGobotResponseIs(t, "!remember asdf == jkl", "Okay, asdf == jkl")
	AssertGobotResponseIs(t, "!whatis asdf", "tests taught me that asdf == jkl")
	AssertGobotResponseIs(t, "!remember asdf == jkl2", "Okay, asdf == jkl2\n(was: jkl by tests)")
	AssertGobotResponseIs(t, "!whatis asdf", "tests taught me that asdf == jkl2")
	RunGobotCommand("!forget asdf")
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

func TestShortestRowFinder(t *testing.T) {
	s := plugins.NewShortestRowFinder()
	item := &plugins.RememberRow{
		Key: "key",
	}
	s.AddItem(item)
	result := s.Result()
	assert.Equal(t, 1, len(result))
}
