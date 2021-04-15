package pkg

import (
	"math"
	"testing"
)

const maxErrorPrice = 0.5

func TestCalculateUnrealizedPNLUSDMargined(t *testing.T) {
	var tests = []struct {
		inputPosition *Position
		expectedPNL   float64
	}{
		{&Position{EntryPrice: 129.9, ClosePrice: 126.2195, Margin: 10, Size: 300, Leverage: 30, Direction: LONG}, -1104},
		{&Position{EntryPrice: 7000, ClosePrice: 7500, Margin: 0.1, Size: 0.2, Leverage: 2, Direction: LONG}, 100},
		{&Position{EntryPrice: 6000, ClosePrice: 5000, Margin: 0.05, Size: 0.4, Leverage: 8, Direction: SHORT}, 400},
	}

	marketHandler := newMarketHandler(USDFutures, 0.0002, 0.0004)

	for _, test := range tests {
		if unrealizedPNL := marketHandler.unrealizedPNL(test.inputPosition, test.inputPosition.ClosePrice); math.Abs(unrealizedPNL-test.expectedPNL) > maxErrorPrice {
			t.Errorf("The expected unrealized PNL was %f but the result is %f", test.expectedPNL, unrealizedPNL)
		}
	}
}

func TestLiquidationPriceUSDMargined(t *testing.T) {
	var tests = []struct {
		inputPosition *Position
		expectedPrice float64
	}{
		{&Position{EntryPrice: 10000, Margin: 1, Leverage: 50, Direction: LONG}, 9850},
		{&Position{EntryPrice: 8000, Margin: 1, Leverage: 40, Direction: SHORT}, 8160},
	}

	marketHandler := newMarketHandler(USDFutures, 0.0002, 0.0004)

	for _, test := range tests {
		if liquidationPrice := marketHandler.liquidationPrice(test.inputPosition); math.Abs(liquidationPrice-test.expectedPrice) > maxErrorPrice {
			t.Errorf("The expected liquidation price was %f but the result is %f", test.expectedPrice, liquidationPrice)
		}
	}
}
