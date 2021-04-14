package pkg

import (
	"math"
	"testing"
)

const maxError = 0.000001

//Fees based on Binance VIP 0 without discount- https://www.binance.com/en/support/faq/360033544231
func TestCreateExchangeHandler(t *testing.T) {
	handler := NewExchangeHandler(USDFutures, 0.020, 0.040, 1)
	if handler.makerFee != 0.00020 {
		t.Errorf("The wrong value was generated for base calculations of maker fee")
	}

	if handler.takerFee != 0.00040 {
		t.Errorf("The wrong value was generated for base calculations of taker fee")
	}

	if handler.amountPerTrade != 0.01 {
		t.Errorf("The wrong amount was generated for base calculations of value used per trade")
	}

	if handler.market != USDFutures {
		t.Errorf("The wrong market was assigned to the handler")
	}
}

func TestOpenCloseUSDPosition(t *testing.T) {
	var tests = []struct {
		direction            Direction
		leverage             uint
		prices               []DataPoint
		makerFee, TakerFee   float64
		takeProfit, stoploss float64
		percentagePerTrade   float64
		expectedPosition     Position
	}{
		{LONG, 17, []DataPoint{CreateData(100), CreateData(110), CreateData(123)}, 0.020, 0.040, 120, 90, 1, Position{
			RealizedPNL: 390.9898, Size: 17, Margin: 1, TakeProfit: 120, Stoploss: 90, Direction: LONG, TotalFeePaid: 0.0102,
		}},
		{LONG, 0, []DataPoint{CreateData(90), CreateData(110), CreateData(130)}, 0.020, 0.040, 130, 80, 1, Position{
			RealizedPNL: 39.9994, Size: 1, Margin: 1, TakeProfit: 130, Stoploss: 80, Direction: LONG, TotalFeePaid: 0.0006,
		}},
		{SHORT, 0, []DataPoint{CreateData(100), CreateData(90), CreateData(75)}, 0.020, 0.040, 75, 120, 1, Position{
			RealizedPNL: 24.9994, Size: 1, Margin: 1, TakeProfit: 75, Stoploss: 120, Direction: SHORT, TotalFeePaid: 0.0006,
		}},
	}

	for _, test := range tests {
		handler := NewExchangeHandler(USDFutures, test.makerFee, test.TakerFee, test.percentagePerTrade)
		handler.SetBalance(100)
		handler.onPriceChange(test.prices[0])

		handler.OpenMarketOrder(test.direction, test.leverage)
		handler.SetTakeProfit(test.takeProfit)
		handler.SetStoploss(test.stoploss)

		for _, latestPrice := range test.prices[1:] {
			handler.onPriceChange(latestPrice)
		}

		if handler.openPosition != nil || len(handler.tradeHistory) != 1 {
			t.Errorf("The position didnt close properly")
		}

		resultPosition := handler.tradeHistory[0]

		if resultPosition.Size != test.expectedPosition.Size || resultPosition.Margin != test.expectedPosition.Margin ||
			!isDifferent(resultPosition.TotalFeePaid, test.expectedPosition.TotalFeePaid) ||
			resultPosition.Stoploss != test.expectedPosition.Stoploss ||
			resultPosition.RealizedPNL != test.expectedPosition.RealizedPNL ||
			resultPosition.TakeProfit != test.expectedPosition.TakeProfit ||
			resultPosition.Direction != test.expectedPosition.Direction {
			t.Errorf("The traded position finished containing wrong values\nThe expected position is:\n%+v The result was:\n%+v",
				test.expectedPosition, resultPosition)
		}
	}
}

func isDifferent(x, y float64) bool {
	return math.Abs(x-y) < maxError
}

func CreateData(value float64) DataPoint {
	return DataPoint{Open: value, High: value, Low: value, Close: value, Volume: value}
}
