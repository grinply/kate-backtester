package kate

import (
	"math"
	"testing"
)

func TestLiquidationPriceCOINMargined(t *testing.T) {
	var tests = []struct {
		inputPosition *Position
		expectedPrice float64
	}{
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 50, Direction: LONG}, 7882},
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 50, Direction: SHORT}, 8121.50},
	}

	marketHandler := newMarketHandler(CoinMarginedFutures, 0.0002, 0.0004)

	for _, test := range tests {
		if liquidationPrice := marketHandler.liquidationPrice(test.inputPosition); math.Abs(liquidationPrice-test.expectedPrice) > maxErrorPrice {
			t.Errorf("The expected liquidation price was %f but the result is %f", test.expectedPrice, liquidationPrice)
		}
	}
}
