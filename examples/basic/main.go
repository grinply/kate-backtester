package main

import (
	"fmt"

	kate "github.com/victorl2/kate-backtester/pkg"
)

type SimpleStrategy struct{}

func main() {
	data, err := kate.PricesFromCSV("../../testdata/ETHUSD5.csv")

	if err != nil {
		panic("could`t load data." + err.Error())
	}

	kate.NewBacktester(&SimpleStrategy{}, data)
	backtester := kate.NewBacktester(&SimpleStrategy{}, data)
	fmt.Println(backtester.Run())
}

//OpenNewPosition process the next data point and checks if a position should be opened
func (stg *SimpleStrategy) OpenNewPosition(latestPrices []kate.DataPoint) *kate.OpenPositionEvt {
	latest := len(latestPrices) - 1

	if latestPrices[latest].Close > latestPrices[latest-1].Close {
		return &kate.OpenPositionEvt{Direction: kate.LONG, Leverage: 30}
	}
	return nil
}

//SetStoploss defines a stoploss for the current open position
func (stg *SimpleStrategy) SetStoploss(openPosition kate.Position) *kate.StoplossEvt {
	if openPosition.Direction == kate.LONG && openPosition.Stoploss <= 0 {
		return &kate.StoplossEvt{Price: openPosition.EntryPrice * 0.995}
	}
	return nil
}

//SetTakeProfit defines a takeprofit for the current open position
func (stg *SimpleStrategy) SetTakeProfit(openPosition kate.Position) *kate.TakeProfitEvt {
	if openPosition.Direction == kate.LONG && openPosition.TakeProfit <= 0 {
		return &kate.TakeProfitEvt{Price: openPosition.EntryPrice * 1.005}
	}
	return nil
}
