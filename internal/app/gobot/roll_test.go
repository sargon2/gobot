package gobot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertParse(t *testing.T, input string, expectedRolls []OneRoll) {
	result, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedRolls, result)
}

func assertParseOneResult(t *testing.T, input string, numDice, diceSize int) {
	assertParse(t, input, []OneRoll{{
		NumDice:  numDice,
		DiceSize: diceSize,
	}})
}

func assertParseError(t *testing.T, input string) {
	_, err := Parse(input)
	assert.NotNil(t, err)
}

func TestWhatever(t *testing.T) {
	assertParseOneResult(t, "0", 0, 1)
	assertParseOneResult(t, "3", 3, 1)
	assertParseOneResult(t, "3d6", 3, 6)
	assertParseOneResult(t, "3D6", 3, 6)
	assertParseOneResult(t, "3d6 asdf", 3, 6)
	assertParseOneResult(t, "2d8", 2, 8)
	assertParse(t, "2d8+4", []OneRoll{{2, 8}, {4, 1}})
	assertParse(t, "  2  d  8  +  4  ", []OneRoll{{2, 1}, {1, 8}, {4, 1}})
	assertParse(t, "2d8+4 5", []OneRoll{{2, 8}, {4, 1}, {5, 1}})
	assertParse(t, "2d8+3d6", []OneRoll{{2, 8}, {3, 6}})
	assertParse(t, "2d8 3d6", []OneRoll{{2, 8}, {3, 6}})
	assertParse(t, "3+4", []OneRoll{{3, 1}, {4, 1}})
	assertParse(t, "3 4", []OneRoll{{3, 1}, {4, 1}})
	assertParse(t, "3 asdf 4", []OneRoll{{3, 1}})
	assertParse(t, "d6", []OneRoll{{1, 6}})
	assertParse(t, "d6 d8", []OneRoll{{1, 6}, {1, 8}})
	assertParseError(t, "asdf")
	assertParseError(t, "3d6d7")
	assertParseError(t, "3.5d6")
	assertParseError(t, "3d6.5")
	assertParseError(t, "1dSTONKS")
}
