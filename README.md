
[![Go Report Card](https://goreportcard.com/badge/github.com/victorl2/quick-backtest?style=flat-square)](https://goreportcard.com/report/github.com/victorl2/kate-backtester)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=victorl2_quick-backtest&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=victorl2_quick-backtest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](/LICENSE.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/victorl2/kate-backtester.svg)](https://pkg.go.dev/github.com/victorl2/kate-backtester)
# Kate Backtester
A fast and simple backtest implementation for **algorithmic trading** focused on [cryptocurrencies](https://en.wikipedia.org/wiki/Cryptocurrency#:~:text=A%20cryptocurrency%2C%20crypto%20currency%20or,creation%20of%20additional%20coins%2C%20and) written in golang.

## Data
The price data used to run the backtests can be from any time interval, but it must contain a [**OHLCV**](https://en.wikipedia.org/wiki/Open-high-low-close_chart) structure _(**O**pen **H**igh **L**ow **C**lose **V**olume)_. It is possible to load data from **csv** files and the [**postgresql** database](https://www.postgresql.org/).

## Usage
To start using **kate backtester** you will need to implement the [**Strategy interface**](https://github.com/victorl2/kate-backtester/blob/main/pkg/strategy.go) and provide a **csv** a dataset with the following format:

| open      | high      | low       | close     | volume     
|:---------:|:---------:|:---------:|:---------:|:----------
| 7922.0700 | 7924.9900 | 7920.1600 | 7924.7500 | 9.90606700 
| 7923.4300 | 7929.1400 | 7920.8000 | 7922.9000 | 15.83760800
| 7923.1300 | 7934.0900 | 7922.9000 | 7932.2600 | 9.98577900

The Strategy interface contains 3 functions that describe how/when to trade:

+ OpenNewPosition
+ SetStoploss
+ SetTakeProfit

### OpenNewPosition

This function is responsible for opening new trade positions when there are none open already, the function is called with every new price data to check, a nil return denotes that no positions should be open yet. When opening a position a OpenPositionEvt is returned containing the **Direction** for the trade _(LONG/SHORT)_ and the desired [leverage](https://blog.earn2trade.com/leverage-trading/), a possible return would be `return &kate.OpenPositionEvt{Direction: kate.LONG, Leverage: 30}`

### SetStoploss
As the name already implies this function is responsible for setting the stoploss price for the **already open position**, the function is called when new price data is avaliable and a position is open. This function makes possible changing the **stoploss** dynamically as the position evolves, the updated PNL is avaliable for checking. A nil return denotes that no changes should be made, a example return would be `return &kate.StoplossEvt{Price: openPosition.EntryPrice * 0.995}` 

### SetTakeProfit
This function has the same behavior as **SetStoploss** but instead it manipulates the take profit price. A example return would be `return &kate.TakeProfitEvt{Price: openPosition.EntryPrice * 1.005}`

A example implementation where a strategy opens a long position every time the latest close price is higher than the last close would be: 

```go
package main

import (
	"fmt"

	kate "github.com/victorl2/kate-backtester/pkg"
)

type simpleStrategy struct{}

func main() {
	data, err := kate.PricesFromCSV("../../testdata/ETHUSD5.csv")
	kate.NewBacktester(&simpleStrategy{}, data)

	if err != nil {
		panic("could`t load data." + err.Error())
	}

	backtester := kate.NewBacktester(&simpleStrategy{}, data)
	fmt.Println(backtester.Run())
}

//ProcessNextPriceData process the next data point and checks if a position should be opened
func (stg *simpleStrategy) OpenNewPosition(latestPrices []kate.DataPoint) *kate.OpenPositionEvt {
	latest := len(latestPrices) - 1

	if latestPrices[latest].Close > latestPrices[latest-1].Close {
		return &kate.OpenPositionEvt{Direction: kate.LONG, Leverage: 30}
	}
	return nil
}

//SetStoploss defines a stoploss for the current open position
func (stg *simpleStrategy) SetStoploss(openPosition kate.Position) *kate.StoplossEvt {
	if openPosition.Direction == kate.LONG && openPosition.Stoploss <= 0 {
		return &kate.StoplossEvt{Price: openPosition.EntryPrice * 0.995}
	}
	return nil
}

//SetTakeProfit defines a takeprofit for the current open position
func (stg *simpleStrategy) SetTakeProfit(openPosition kate.Position) *kate.TakeProfitEvt {
	if openPosition.Direction == kate.LONG && openPosition.TakeProfit <= 0 {
		return &kate.TakeProfitEvt{Price: openPosition.EntryPrice * 1.005}
	}
	return nil
}
```
