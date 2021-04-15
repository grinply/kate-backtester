package pkg

import (
	"math"
	"testing"
)

const maxError = 0.001

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
}

func TestOpenCloseUSDPosition(t *testing.T) {
	var tests = []struct {
		direction            Direction
		leverage             uint
		openPrice            DataPoint
		closePrice           DataPoint
		makerFee, TakerFee   float64
		takeProfit, stoploss float64
		percentagePerTrade   float64
		expectedPosition     Position
	}{
		{LONG, 32, CreateData(250), CreateData(267.7), 0.020, 0.040, 267.7, 248, 20, Position{
			RealizedPNL: 44.9189, Size: 2.56, Margin: 0.08, TakeProfit: 267.7, Stoploss: 248, Direction: LONG, TotalFeePaid: 0.39306,
		}},

		{LONG, 25, CreateData(6000), CreateData(5910), 0.020, 0.040, 6200, 5910, 72, Position{
			RealizedPNL: -28.4292, Size: 0.3, Margin: 0.012, TakeProfit: 6200, Stoploss: 5910, Direction: LONG, TotalFeePaid: 1.4292,
		}},

		{LONG, 40, CreateData(1000), CreateData(970), 0.020, 0.040, 1100, 0, 20, Position{
			RealizedPNL: -16.9472, Size: 0.8, Margin: 0.02, TakeProfit: 1100, Stoploss: 0, Direction: LONG, TotalFeePaid: 0.9472,
		}},

		{SHORT, 50, CreateData(1200), CreateData(1175), 0.020, 0.040, 1175, 1210, 81, Position{
			RealizedPNL: 81.961, Size: 3.375, Margin: 0.0675, TakeProfit: 1175, Stoploss: 1210, Direction: SHORT, TotalFeePaid: 2.413,
		}},

		//800 * ( 1 + (1/17) - 0,005)
		//0,03478 * (2300 - 2430) - (2300* 0,0004 * 0,03478) - (2430 * 0,0004 * 0,03478)
		{SHORT, 8, CreateData(2300), CreateData(2430), 0.020, 0.040, 2290, 2430, 10, Position{
			RealizedPNL: -4.58754, Size: 0.0347826, Margin: 0.004347, TakeProfit: 2290, Stoploss: 2430, Direction: SHORT, TotalFeePaid: 0.06580,
		}},

		{SHORT, 17, CreateData(800), CreateData(1000), 0.020, 0.040, 790, 1210, 20, Position{
			RealizedPNL: -18.7226, Size: 0.425, Margin: 0.025, TakeProfit: 790, Stoploss: 1210, Direction: SHORT, TotalFeePaid: 0.4226,
		}},
	}

	for _, test := range tests {
		handler := NewExchangeHandler(USDFutures, test.makerFee, test.TakerFee, test.percentagePerTrade)
		handler.SetBalance(100)
		handler.onPriceChange(test.openPrice)

		handler.OpenMarketOrder(test.direction, test.leverage)
		handler.SetTakeProfit(test.takeProfit)
		handler.SetStoploss(test.stoploss)

		handler.onPriceChange(test.closePrice)

		if handler.openPosition != nil || len(handler.tradeHistory) != 1 {
			t.Errorf("The position didnt close properly")
		}

		resultPosition := handler.tradeHistory[0]

		if !isEqual(resultPosition.Size, test.expectedPosition.Size) ||
			!isEqual(resultPosition.Margin, test.expectedPosition.Margin) ||
			!isEqual(resultPosition.TotalFeePaid, test.expectedPosition.TotalFeePaid) ||
			resultPosition.Stoploss != test.expectedPosition.Stoploss ||
			!isEqual(resultPosition.RealizedPNL, test.expectedPosition.RealizedPNL) ||
			resultPosition.TakeProfit != test.expectedPosition.TakeProfit ||
			resultPosition.Direction != test.expectedPosition.Direction {
			t.Errorf("The traded position finished containing wrong values\nThe expected position is:\n%+v The result was:\n%+v",
				test.expectedPosition, resultPosition)
		}
	}
}

func isEqual(x, y float64) bool {
	return math.Abs(x-y) < maxError
}

func CreateData(value float64) DataPoint {
	return DataPoint{Open: value, High: value, Low: value, Close: value, Volume: value}
}
