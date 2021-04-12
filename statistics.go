package main

type Statistics struct {
	ROI         float64
	SharpeRatio float64
	TotalTrades uint
	WinRate     float64 //Percentage of wins
	MaxDrawdown float64 //Percentage for the maximum drawdown after applying the strategy
}
