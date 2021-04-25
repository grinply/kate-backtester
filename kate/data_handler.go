package kate

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

//DataHandler is a wrapper that packages the required data for running backtesting simulation.
type DataHandler struct {
	Prices []DataPoint
}

//Position is the representation of a traded position
type Position struct {
	Direction              Direction
	Size                   float64 //total size of the position including leverage
	Leverage               uint    //the multiplier for increasing the total traded position
	Margin                 float64 //the amount of collateral in COIN that is backing the position
	EntryPrice, ClosePrice float64
	Stoploss, TakeProfit   float64
	UnrealizedPNL          float64
	RealizedPNL            float64
	TotalFeePaid           float64
	LiquidationPrice       float64
}

//Required columns in the CSV file
var csvColumns = []string{"open", "high", "low", "close", "volume"}

//newDataHandler creates and initializes a DataHandler with pricing data and executes the required setup
func newDataHandler(prices []DataPoint) *DataHandler {
	return &DataHandler{
		Prices: prices,
	}
}

//PricesFromCSV reads all csv data in the OHLCV format to the DataHandler and returns if a error occurred
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
			value, err := strToFloat(line[i])
			if err != nil {
				return nil, err
			}
			numbers[i] = value

		}

		prices = append(prices, DataPoint{
			open:   numbers[0],
			high:   numbers[1],
			low:    numbers[2],
			close:  numbers[3],
			volume: numbers[4],
		})
	}

	return newDataHandler(prices), nil
}

//strToFloat converts a string value to float64, in case of error Panic
func strToFloat(str string) (float64, error) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return number, nil
	}
	return -1, fmt.Errorf(`invalid parameter '%v' was found in the provided csv. 
		Make sure the csv contain only valid float numbers`, str)
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
