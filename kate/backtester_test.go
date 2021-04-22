package kate

import (
	"testing"

	"github.com/go-test/deep"
)

type simpleStrategy struct {
	lastPrice    *DataPoint
	currentPrice *DataPoint
}

func TestSimpleTradingStrategyRun(t *testing.T) {
	var tests = []struct {
		filePath       string
		amountPerTrade float64
		initialBalance float64
		expectedResult *Statistics
	}{
		{"../testdata/ETHUSD1.csv", 1, 100, &Statistics{TotalDataPoints: 1757, TotalTrades: 43, MaxDrawdown: 0.016799672003129037,
			NetProfit: -0.7494000000001364, ROIPercentage: -0.7494000000001364, SharpeRatio: -0.3117941697762629, WinRate: 0.5116279069767442}},
		{"../testdata/ETHUSD2.csv", 20, 1000, &Statistics{TotalDataPoints: 1264, TotalTrades: 21, NetProfit: -11.872800000000666,
			SharpeRatio: -0.24420303389734746, WinRate: 0.47619047619047616, MaxDrawdown: 0.014473997331444076, ROIPercentage: -1.1872800000000665}},
		{"../testdata/ETHUSD3.csv", 5, 2000, &Statistics{TotalDataPoints: 2946, TotalTrades: 71, MaxDrawdown: 0.00538909887956826, WinRate: 0.5211267605633803,
			SharpeRatio: -0.1811079374911933, NetProfit: -5.155349999998634, ROIPercentage: -0.2577674999999317}},
		{"../testdata/ETHUSD4.csv", 7, 300, &Statistics{TotalDataPoints: 21265, TotalTrades: 133, MaxDrawdown: 0.06771563065859225, WinRate: 0.5338345864661654,
			SharpeRatio: -1.5822095979921864, NetProfit: -9.900870000000737, ROIPercentage: -3.300290000000246}},
		{"../testdata/ETHUSD5.csv", 10, 1000, &Statistics{TotalDataPoints: 43200, TotalTrades: 856, MaxDrawdown: 0.21341551110001716, WinRate: 0.49182242990654207,
			SharpeRatio: -3.2570713833447655, NetProfit: -208.15910000001497, ROIPercentage: -20.815910000001498}},
		{"../testdata/mockdata.csv", 5, 1000, &Statistics{TotalDataPoints: 22, TotalTrades: 4, WinRate: 0.75, MaxDrawdown: 0.0008679817866541949,
			ROIPercentage: 0.11098500000000514, NetProfit: 1.1098500000000513, SharpeRatio: 0.003966684810396764}},
	}

	for _, test := range tests {
		data, err := PricesFromCSV(test.filePath)
		if err != nil {
			t.Fatal("could`t load data." + err.Error())
		}

		backtester := NewBacktester(newSimpleStrategy(), data)
		backtester.SetBalance(test.initialBalance)

		if test.amountPerTrade > 0 {
			backtester.SetFixedTradeAmount(test.amountPerTrade)
		}

		result := backtester.Run()
		if diff := deep.Equal(result, test.expectedResult); diff != nil {
			t.Error("the result from the backtest with file (", test.filePath,
				") execution does not match the expected value.\nThe Diff is", diff)
		}
	}

}

func newSimpleStrategy() *simpleStrategy {
	return &simpleStrategy{}
}

//PreProcessIndicators nothing to do
func (strategy *simpleStrategy) PreProcessIndicators(latestPrice DataPoint) {
	strategy.lastPrice = strategy.currentPrice
	strategy.currentPrice = &latestPrice
}

//OpenNewPosition process the next data point and checks if a position should be opened
func (strategy *simpleStrategy) OpenNewPosition(latestPrices DataPoint) *OpenPositionEvt {
	if strategy.lastPrice != nil && strategy.currentPrice.Close() > strategy.lastPrice.Close() {
		return &OpenPositionEvt{Direction: LONG, Leverage: 30}
	}
	return nil
}

//SetStoploss defines a stoploss for the current open position
func (strategy *simpleStrategy) SetStoploss(openPosition Position) *StoplossEvt {
	if openPosition.Direction == LONG && openPosition.Stoploss <= 0 {
		return &StoplossEvt{Price: openPosition.EntryPrice * 0.995}
	}
	return nil
}

//SetTakeProfit defines a takeprofit for the current open position
func (strategy *simpleStrategy) SetTakeProfit(openPosition Position) *TakeProfitEvt {
	if openPosition.Direction == LONG && openPosition.TakeProfit <= 0 {
		return &TakeProfitEvt{Price: openPosition.EntryPrice * 1.005}
	}
	return nil
}
