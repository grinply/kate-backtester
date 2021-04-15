package main

import (
	"fmt"
	bt "quickbacktest/pkg"
	"time"
)

type simpleStrategy struct{}

func main() {
	data, err := bt.PricesFromCSV("testdata/complete.csv")
	if err != nil {
		panic("could`t load data." + err.Error())
	}

	backtester := bt.NewBacktester(newSimpleStrategy(), data)
	backtester.SetFixedTradeAmount(5)
	start := time.Now()
	result := backtester.Run()
	fmt.Printf("\nIt took %v to run the backtest\n", time.Since(start))
	fmt.Printf("%+v", result)
}

func newSimpleStrategy() *simpleStrategy {
	return &simpleStrategy{}
}

//ProcessNextPriceData process the next data point and checks if a position should be opened
func (strategy *simpleStrategy) ProcessNextPriceData(latestPrices []bt.DataPoint) *bt.OpenPositionEvt {
	latest := len(latestPrices) - 1
	if latestPrices[latest].Close > latestPrices[latest-1].Close &&
		latestPrices[latest-1].Close > latestPrices[latest-2].Close {
		return &bt.OpenPositionEvt{Direction: bt.LONG, Leverage: 30}
	}
	return nil
}

//SetStoploss defines a stoploss for the current open position
func (strategy *simpleStrategy) SetStoploss(openPosition bt.Position) *bt.StoplossEvt {
	if openPosition.Direction == bt.LONG && openPosition.Stoploss <= 0 {
		return &bt.StoplossEvt{Price: openPosition.EntryPrice * 0.97}
	}
	return nil
}

//SetTakeProfit defines a takeprofit for the current open position
func (strategy *simpleStrategy) SetTakeProfit(openPosition bt.Position) *bt.TakeProfitEvt {
	if openPosition.Direction == bt.LONG && openPosition.TakeProfit <= 0 {
		return &bt.TakeProfitEvt{Price: openPosition.EntryPrice * 1.18}
	}
	return nil
}
