package kate

//DataPoint is a unit that encapsulates OHLCV price data
type DataPoint struct {
	Event
	open, high, low, close, volume float64
}

//OHLCV - Represents a datapoint in candle format that contain Open,High, Low, Close prices and Volume data
type OHLCV interface {
	//Close is the finish price when a candlestick has concluded
	Close() float64

	//Open is the starting price for a candlestick
	Open() float64

	//High is the highest price reached between the time a candlestick is open and closed
	High() float64

	//Low is the lowest price reached between the time a candlestick is open and closed
	Low() float64

	//Volume is the amount of assets traded in the timeframe for the current candlestick
	Volume() float64
}

//Open is the starting price for a candlestick
func (candle DataPoint) Open() float64 {
	return candle.open
}

//High is the highest price reached between the time a candlestick is open and closed
func (candle DataPoint) High() float64 {
	return candle.high
}

//Low is the lowest price reached between the time a candlestick is open and closed
func (candle DataPoint) Low() float64 {
	return candle.low
}

//Close is the finish price when a candlestick has concluded
func (candle DataPoint) Close() float64 {
	return candle.close
}

//Volume is the amount of assets traded in the timeframe for the current candlestick
func (candle DataPoint) Volume() float64 {
	return candle.volume
}
