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
	prices []DataPoint
}

type DataPoint struct {
	Open, High, Low, Close, Volume float64
}

type Position struct {
	EntryPrice           float64
	Stoploss, Takeprofit float64
	UnrealizedPNL        float64
}

//Required columns in the CSV file
var csvColumns = []string{"open", "high", "low", "close", "volume"}

//NextValue returns the next value in the stream of datapoints, a null return value denotes the end of the stream
func (handler *DataHandler) NextValue() *DataPoint {
	return &DataPoint{}
}

//LoadPricesFromCSV reads all csv data in the OHLCV format to the DataHandler and returns if a error ocurred
func (handler *DataHandler) LoadPricesFromCSV(csvFilePath string) error {
	csvFile, _ := os.Open(csvFilePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	//Reading first line header and validating the required columns
	if line, error := reader.Read(); error != nil || !isCSVHeaderValid(line) {
		return fmt.Errorf(`error reading header with columns in the csv.
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
				return err
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

	handler.prices = prices
	return nil
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

func isCSVHeaderValid(firstLine []string) bool {
	for i, column := range csvColumns {
		if strings.ToLower(firstLine[i]) != column {
			return false
		}
	}
	return true
}
