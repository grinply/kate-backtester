package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type DataHandler struct {
	counter, windowSize int
	prices              []DataPoint
}

type DataPoint struct {
	Event
	Open, High, Low, Close, Volume float64
}

type Position struct {
	Direction            Direction
	Size                 float64 //total size of the position including leverage
	Leverage             uint
	Margin, EntryPrice   float64
	ClosePrice           float64
	Stoploss, TakeProfit float64
	UnrealizedPNL        float64
	RealizedPNL          float64
	TotalFeePaid         float64
}

type AggregatedDataPoints struct {
	Event
	datapoints []DataPoint
}

//Required columns in the CSV file
var csvColumns = []string{"open", "high", "low", "close", "volume"}

//NextValues returns AggregatedDataPoints with the next values in the stream of datapoints (containing the lastest windowSize of values).
//a nil return denotes the end for the stream
func (handler *DataHandler) nextValues() *AggregatedDataPoints {
	if handler.counter < len(handler.prices) {
		return &AggregatedDataPoints{
			datapoints: handler.prices[handler.counter-handler.windowSize : handler.counter],
		}
	}
	return nil
}

//LoadPricesFromCSV reads all csv data in the OHLCV format to the DataHandler and returns if a error ocurred
func PricesFromCSV(csvFilePath string) (*DataHandler, error) {
	csvFile, _ := os.Open(csvFilePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	//Reading first line header and validating the required columns
	if line, error := reader.Read(); error != nil || !isCSVHeaderValid(line) {
		return nil, fmt.Errorf(`error reading header with columns in the csv.
				Make sure the CSV has the columns Open, High, Low, Close, Volume`)
	}

	var prices []DataPoint
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		//Checking each OHLCV value in the csv
		var numbers [5]float64
		for i := 0; i < 5; i++ {
			if value, err := strToFloat(line[i]); err != nil {
				return nil, err
			} else {
				numbers[i] = value
			}
		}

		prices = append(prices, DataPoint{
			Open:   numbers[0],
			High:   numbers[1],
			Low:    numbers[2],
			Close:  numbers[3],
			Volume: numbers[4],
		})
	}

	handler := &DataHandler{}
	handler.initData(prices)
	return handler, nil
}

//initData initializes the DataHandler with pricing data and executes the required setup
func (handler *DataHandler) initData(prices []DataPoint) {
	handler.prices = prices
	handler.counter = handler.windowSize
}

//strToFloat converts a string value to float64, in case of error Panic
func strToFloat(str string) (float64, error) {
	if number, err := strconv.ParseFloat(str, 64); err == nil {
		return number, nil
	} else {
		return -1, fmt.Errorf(`invalid parameter '%v' was found in the provided csv. 
		Make sure the csv contain only valid float numbers`, str)
	}
}

//Check if the first line with columns of the csv are in the valid format
func isCSVHeaderValid(firstLine []string) bool {
	for i, column := range csvColumns {
		if strings.ToLower(firstLine[i]) != column {
			return false
		}
	}
	return true
}
