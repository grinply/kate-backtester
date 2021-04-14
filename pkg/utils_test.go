package pkg

import (
	"math"
	"testing"
)

const maxErrorPrice = 0.5

func TestLiquidationPriceUSDMargined(t *testing.T) {
	var tests = []struct {
		inputPosition *Position
		expectedPrice float64
	}{
		{&Position{EntryPrice: 10000, Margin: 1, Leverage: 50, Direction: LONG}, 9850},
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 40, Direction: SHORT}, 8160},
	}

	for _, test := range tests {
		if liquidationPrice := USDMarginedLiquidationPrice(test.inputPosition); math.Abs(liquidationPrice-test.expectedPrice) > maxErrorPrice {
			t.Errorf("The expected liquidation price was %f but the result is %f", test.expectedPrice, liquidationPrice)
		}
	}
}

func TestLiquidationPriceCOINMargined(t *testing.T) {
	var tests = []struct {
		inputPosition *Position
		expectedPrice float64
	}{
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 50, Direction: LONG}, 7882},
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 50, Direction: SHORT}, 8121.50},
	}

	for _, test := range tests {
		if liquidationPrice := CoinMarginedLiquidationPrice(test.inputPosition); math.Abs(liquidationPrice-test.expectedPrice) > maxErrorPrice {
			t.Errorf("The expected liquidation price was %f but the result is %f", test.expectedPrice, liquidationPrice)
		}
	}
}
