package pkg

//Statistics are the results based on trades executed on a backtest run
type Statistics struct {
	ROIPercentage   float64
	NetProfit       float64
	SharpeRatio     float64
	WinRate         float64 //Percentage of wins
	MaxDrawdown     float64 //Percentage for the maximum drawdown after applying the strategy
	TotalTrades     int
	TotalDataPoints int
}

//calculateStatistics calculates metrics based on a trade history
func (bt *Backtester) calculateStatistics(initialBalance float64) *Statistics {
	tradeHistory := bt.exchangeHandler.tradeHistory
	wins, balance, peakProfit, bottomProfit := 0, initialBalance, 0.0, 0.0
	balanceHistory := []float64{initialBalance}

	for _, position := range tradeHistory {
		if position.RealizedPNL >= 0 {
			wins++
		}
		//fmt.Printf("open: %f | close: %f | liquidation: %f | margin: %f | size: %f | PNL %f\n", position.EntryPrice,
		//	position.ClosePrice, position.LiquidationPrice, position.Margin*position.EntryPrice, position.Size, position.RealizedPNL)

		balance += position.RealizedPNL
		balanceHistory = append(balanceHistory, balance)

		if peakProfit < balance {
			peakProfit = balance
			bottomProfit = balance
		} else if bottomProfit > balance {
			bottomProfit = balance
		}
	}
	stats := &Statistics{
		ROIPercentage:   100 * ((balance - initialBalance) / initialBalance),
		NetProfit:       balance - initialBalance,
		TotalTrades:     len(tradeHistory),
		WinRate:         float64(wins) / float64(len(tradeHistory)),
		MaxDrawdown:     (peakProfit - bottomProfit) / peakProfit,
		TotalDataPoints: len(bt.dataHandler.prices),
	}

	stats.SharpeRatio = sharpe(stats.NetProfit, 0.0, stdDev(balanceHistory))
	return stats
}
