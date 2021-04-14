package pkg

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
