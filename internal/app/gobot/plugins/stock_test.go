package gobot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	plugins "github.com/sargon2/gobot/internal/app/gobot/plugins"
)

func assertStockSplit(t *testing.T, input string, expected []string) {
	result := plugins.StockSplit(input)
	assert.Equal(t, expected, result)
}

func TestStockSplit(t *testing.T) {
	assertStockSplit(t, "", []string{})
	assertStockSplit(t, "aapl", []string{"aapl"})
	assertStockSplit(t, "aapl btc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "    aapl   btc-usd    ", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "aapl,btc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "aapl, btc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "aapl,,btc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "aapl  ,  ,  btc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "aapl\tbtc-usd", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "1,2,3,4,5", []string{"1", "2", "3", "4", "5"})
	assertStockSplit(t, "1;,/\\ 2", []string{"1", "2"})
	assertStockSplit(t, "\"aapl\"", []string{"aapl"})
	assertStockSplit(t, "\"aapl\",\"btc-usd\"", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "\"aapl,btc-usd\"", []string{"aapl", "btc-usd"})
	assertStockSplit(t, "'aapl'", []string{"aapl"})
	assertStockSplit(t, "`aapl`", []string{"aapl"})
}

func assertFloatFormat(t *testing.T, f float64, expected string) {
	result := plugins.FloatFormat(f)
	assert.Equal(t, expected, result)
}

func TestFormatFloat(t *testing.T) {
	assertFloatFormat(t, 1.1, "1.1  ")
	assertFloatFormat(t, 0.1, "0.1  ")
	assertFloatFormat(t, 400, "400    ")
	assertFloatFormat(t, 3124, "3,124    ")
	assertFloatFormat(t, 3124.567, "3,124.567")
	assertFloatFormat(t, 1111111.111, "1,111,111.111")
	assertFloatFormat(t, 11111111.111, "11,111,111.111")
	assertFloatFormat(t, 111111111.1111, "111,111,111.111")
	assertFloatFormat(t, 1, "1    ")

	assertFloatFormat(t, 0, "0    ")

	assertFloatFormat(t, -1.1, "-1.1  ")
	assertFloatFormat(t, -0.1, "-0.1  ")
	assertFloatFormat(t, -400, "-400    ")
	assertFloatFormat(t, -3124, "-3,124    ")
	assertFloatFormat(t, -3124.567, "-3,124.567")
	assertFloatFormat(t, -1111111.111, "-1,111,111.111")
	assertFloatFormat(t, -11111111.111, "-11,111,111.111")
	assertFloatFormat(t, -111111111.1111, "-111,111,111.111")
	assertFloatFormat(t, -1, "-1    ")
}
