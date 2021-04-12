package main

import (
	"math"
	"testing"
)

const maxError = 0.0000001

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

func TestCloseLongPositionInProfit(t *testing.T) {
	var leverage uint = 5
	handler := NewExchangeHandler(USDFutures, 0.020, 0.040, 1)
	handler.SetBalance(100)
	handler.onPriceChange(100)
	handler.OpenMarketOrder(LONG, leverage)

	handler.SetTakeProfit(120)
	handler.SetStoploss(90)

	handler.onPriceChange(110)
	handler.onPriceChange(123)

	if handler.openPosition != nil || len(handler.tradeHistory) != 1 {
		t.Errorf("The position didnt close properly")
	}

	expectedPNL := 114.997
	tradedPosition := handler.tradeHistory[0]
	if handler.tradeHistory[0].RealizedPNL != expectedPNL {
		t.Errorf("The realized pnl is not the expected value, expected %v received %v",
			expectedPNL, handler.tradeHistory[0].RealizedPNL)
	}

	if tradedPosition.Size != 5 || tradedPosition.Margin != 1 || !isDifferent(tradedPosition.TotalFeePaid, 0.003) ||
		tradedPosition.Stoploss != 90 || tradedPosition.Takeprofit != 120 || tradedPosition.Direction != LONG {
		t.Errorf("The traded position finished containing wrong values")
	}
}

func TestCloseShortPositionInProfit(t *testing.T) {
	handler := NewExchangeHandler(USDFutures, 0.020, 0.040, 1)
	handler.SetBalance(100)
	handler.onPriceChange(100)
	handler.OpenMarketOrder(SHORT, 0)

	handler.SetTakeProfit(80)
	handler.SetStoploss(120)

	handler.onPriceChange(110)
	handler.onPriceChange(100)
	handler.onPriceChange(90)
	handler.onPriceChange(75)

	if handler.openPosition != nil || len(handler.tradeHistory) != 1 {
		t.Errorf("The position didnt close properly")
	}

	expectedPNL := 24.9994
	tradedPosition := handler.tradeHistory[0]
	if tradedPosition.RealizedPNL != expectedPNL {
		t.Errorf("The realized pnl is not the expected value, expected %v received %v",
			expectedPNL, tradedPosition.RealizedPNL)
	}

	if tradedPosition.Size != 1 || tradedPosition.Margin != 1 || !isDifferent(tradedPosition.TotalFeePaid, 0.0006) ||
		tradedPosition.Stoploss != 120 || tradedPosition.Takeprofit != 80 || tradedPosition.Direction != SHORT {
		t.Errorf("The traded position finished containing wrong values")
	}
}

func isDifferent(x, y float64) bool {
	return math.Abs(x-y) < maxError
}
