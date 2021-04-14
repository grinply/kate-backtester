package main

import (
	"fmt"
	bt "quickbacktest/pkg"
)

type SimpleStrategy struct{}

func main() {
	data, err := bt.PricesFromCSV("testdata/complete.csv")
	if err != nil {
		panic("could`t load data." + err.Error())
	}

	backtester := bt.NewBacktester(newSimpleStrategy(), data)
	fmt.Printf("%+v", backtester.Run())
}

func newSimpleStrategy() *SimpleStrategy {
	return &SimpleStrategy{}
}

func (strategy *SimpleStrategy) ProcessNextPriceData(latestPrices []bt.DataPoint) *bt.OpenPositionEvt {
	latest := len(latestPrices) - 1
	if latestPrices[latest].Close > latestPrices[latest-1].Close &&
		latestPrices[latest-1].Close > latestPrices[latest-2].Close {
		return &bt.OpenPositionEvt{Direction: bt.LONG, Leverage: 30}
	}
	return nil
}

func (strategy *SimpleStrategy) SetStoploss(openPosition bt.Position) *bt.StoplossEvt {
	if openPosition.Stoploss <= 0 {
		return &bt.StoplossEvt{Price: openPosition.EntryPrice * 0.95}
	}
	return nil
}

func (strategy *SimpleStrategy) SetTakeProfit(openPosition bt.Position) *bt.TakeProfitEvt {
	if openPosition.TakeProfit <= 0 {
		return &bt.TakeProfitEvt{Price: openPosition.EntryPrice * 1.1}
	}
	return nil
}
