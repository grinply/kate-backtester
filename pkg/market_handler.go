package pkg

//MarketHandler describes market expecific functionality
type MarketHandler interface {
	createPosition(tradeDirection Direction, currentPrice, balance, amountPerTrade float64, leverage uint) *Position
	unrealizedPNL(position *Position, lastTradedPrice float64) float64
	liquidationPrice(position *Position) float64
	marketFee(position *Position) float64
	limitFee(position *Position) float64
	liquidationFee(position *Position) float64
}

func newMarketHandler(market MarketType, makerFee, takerFee float64) MarketHandler {
	switch market {
	case USDFutures:
		return &USDMarket{Market: market, MakerFee: makerFee, TakerFee: takerFee}
	case CoinMarginedFutures:
		return &CoinMarket{Market: market, MakerFee: makerFee, TakerFee: takerFee}
	default:
		return nil
	}
}
