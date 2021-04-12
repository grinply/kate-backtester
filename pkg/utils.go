package main

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000404726-P-L-calculations-Inverse-Contracts-#h_92ba55e9-4bbc-4879-a354-bc62eaa57d4d
func CoinMarginedUnrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * ((1 / position.EntryPrice) - (1 / lastTradedPrice))
	}
	return float64(position.Size) * ((1 / lastTradedPrice) - (1 / position.EntryPrice))
}

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000630066-P-L-calculations-USDT-Contract#Unrealized_P&L
func USDMarginedUnrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * (lastTradedPrice - position.EntryPrice)
	}
	return float64(position.Size) * (position.EntryPrice - lastTradedPrice)
}
