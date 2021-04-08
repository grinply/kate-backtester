package main

type DataHandler struct {
}

type DataPoint struct {
	Open, High, Low, Close float64
}

type Position struct {
	EntryPrice           float64
	Stoploss, Takeprofit float64
	UnrealizedPNL        float64
}

//NextValue returns the next value in the stream of datapoints, a null return value denotes the end of the stream
func (handler *DataHandler) NextValue() *DataPoint {
	return &DataPoint{}
}
