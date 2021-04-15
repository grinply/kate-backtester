package pkg

import "fmt"

//Statistics are the results based on trades executed on a backtest run
type Statistics struct {
	ROIPercentage float64
	NetProfit     float64
	SharpeRatio   float64
	TotalTrades   int
	WinRate       float64 //Percentage of wins
	MaxDrawdown   float64 //Percentage for the maximum drawdown after applying the strategy
}

//calculateStatistics calculates metrics based on a trade history
func calculateStatistics(initialBalance float64, tradeHistory []*Position) *Statistics {
	wins, totalProfit, peakProfit, bottomProfit := 0, 0.0, 0.0, 0.0

	for _, position := range tradeHistory {
		if position.RealizedPNL >= 0 {
			wins++
		}
		fmt.Printf("open: %f | close: %f | liquidation: %f | margin: %f | size: %f | PNL %f\n", position.EntryPrice,
			position.ClosePrice, position.LiquidationPrice, position.Margin*position.EntryPrice, position.Size, position.RealizedPNL)

		totalProfit += position.RealizedPNL

		if peakProfit < totalProfit {
			peakProfit = totalProfit
			bottomProfit = totalProfit
		} else if bottomProfit > totalProfit {
			bottomProfit = totalProfit
		}
	}

	return &Statistics{
		ROIPercentage: 100 * (totalProfit / initialBalance),
		NetProfit:     totalProfit,
		TotalTrades:   len(tradeHistory),
		WinRate:       float64(wins) / float64(len(tradeHistory)),
		MaxDrawdown:   (peakProfit - bottomProfit) / peakProfit,
	}
}
