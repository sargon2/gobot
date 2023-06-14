package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

func assertParse(t *testing.T, input string, expectedRolls []plugins.OneRoll) {
	result, err := plugins.ParseRoll(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedRolls, result)
}

func assertParseOneResult(t *testing.T, input string, numDice, diceSize int) {
	assertParse(t, input, []plugins.OneRoll{{
		NumDice:  numDice,
		DiceSize: diceSize,
	}})
}

func assertParseError(t *testing.T, input string) {
	_, err := plugins.ParseRoll(input)
	assert.NotNil(t, err)
}

func TestParseRoll(t *testing.T) {
	assertParseOneResult(t, "0", 0, 1)
	assertParseOneResult(t, "3", 3, 1)
	assertParse(t, "d6 d6 d6", []plugins.OneRoll{{1, 6}, {1, 6}, {1, 6}})
	assertParse(t, "d6+d6+d6", []plugins.OneRoll{{1, 6}, {1, 6}, {1, 6}})
	assertParse(t, "3d6", []plugins.OneRoll{{1, 6}, {1, 6}, {1, 6}})
	assertParseOneResult(t, "13d6", 13, 6)
	assertParseOneResult(t, "13D6", 13, 6)
	assertParseOneResult(t, "13d6 asdf", 13, 6)
	assertParseOneResult(t, "1d2 1=asdf 2=jkl", 1, 2)
	assertParseOneResult(t, "1d2 1=d3 2=d4", 1, 2)
	assertParseOneResult(t, "1d2 1 = yes, 2 = no", 1, 2)
	assertParseOneResult(t, "1d2  1  =  yes,  2  =  no", 1, 2)
	assertParseOneResult(t, "1d100; spezbrain", 1, 100)
	assertParseError(t, "1d2=1")
	assertParseError(t, "1d2=a")
	assertParseOneResult(t, "12d8", 12, 8)
	assertParse(t, "2d8+4", []plugins.OneRoll{{1, 8}, {1, 8}, {4, 1}})
	assertParse(t, "12d8+4", []plugins.OneRoll{{12, 8}, {4, 1}})
	assertParse(t, "  2  d  8  +  4  ", []plugins.OneRoll{{2, 1}, {1, 8}, {4, 1}})
	assertParse(t, "12d8+4 5", []plugins.OneRoll{{12, 8}, {4, 1}, {5, 1}})
	assertParse(t, "12d8+13d6", []plugins.OneRoll{{12, 8}, {13, 6}})
	assertParse(t, "12d8 13d6", []plugins.OneRoll{{12, 8}, {13, 6}})
	assertParse(t, "3+4", []plugins.OneRoll{{3, 1}, {4, 1}})
	assertParse(t, "3 4", []plugins.OneRoll{{3, 1}, {4, 1}})
	assertParse(t, "3 asdf 4", []plugins.OneRoll{{3, 1}})
	assertParse(t, "d6", []plugins.OneRoll{{1, 6}})
	assertParse(t, "d6 d8", []plugins.OneRoll{{1, 6}, {1, 8}})
	assertParse(t, "1d9223372036854775807", []plugins.OneRoll{{1, 9223372036854775807}})
	assertParseOneResult(t, "10000d6", 10000, 6)
	assertParseError(t, "asdf")
	assertParseError(t, "3d6d7")
	assertParseError(t, "3.5d6")
	assertParseError(t, "3d6.5")
	assertParseError(t, "1dSTONKS")
	assertParseError(t, "10001d6")
	assertParseError(t, "-1d6")
	assertParseError(t, "1d-1")
	assertParseError(t, "-1d-1")
	assertParseError(t, "1d9223372036854775808")
}
