package pkg

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadCSVData(t *testing.T) {
	handler, _ := PricesFromCSV("../testdata/simple_example.csv")

	if len(handler.prices) != 100 {
		t.Errorf("The amount of prices is %v the expected amount is 100", len(handler.prices))
	}

	firstPrice := DataPoint{
		Open:   746.25,
		High:   747.25,
		Low:    746.2,
		Close:  746.95,
		Volume: 1045532,
	}
	if !reflect.DeepEqual(handler.prices[0], firstPrice) {
		t.Errorf("The first price in the csv is %v the expected was %v", handler.prices[0], firstPrice)
	}

	lastPrice := DataPoint{
		Open:   744.4,
		High:   744.6,
		Low:    744.4,
		Close:  744.6,
		Volume: 80597,
	}

	if !reflect.DeepEqual(handler.prices[len(handler.prices)-1], lastPrice) {
		t.Errorf("The last price in the csv is %v the expected was %v", handler.prices[len(handler.prices)-1], lastPrice)
	}

	middlePrice := DataPoint{
		Open:   746.4,
		High:   746.45,
		Low:    746.4,
		Close:  746.45,
		Volume: 5306,
	}
	if !reflect.DeepEqual(handler.prices[66], middlePrice) {
		t.Errorf("The middle price in the csv is %v the expected was %v", handler.prices[66], middlePrice)
	}

}

func TestLoadInvalidCSV(t *testing.T) {
	columnsCSV := []string{"ERROR,Open,High,Low,Close,Volume", "Open,ERROR,High,Low,Close,Volume",
		"Open,High,ERROR,Low,Close,Volume", "Open,High,Low,ERROR,Close,Volume", "Open,High,Low,Close,ERROR,Volume",
		"ERRORHigh,Low,Close,Volume", "Open,ERRORHigh,Low,Close,Volume", "Open,High,ERRORLow,Close,Volume",
		"Open,High,Low,ERRORClose,Volume", "Open,High,Low,Close,ERRORVolume"}

	//Checking if a error is raised with a csv containing a unkown column
	for _, columnLine := range columnsCSV {
		unkownColumnCSV := createTempCSV()
		defer os.Remove(unkownColumnCSV.Name())
		unkownColumnCSV.WriteString(columnLine)

		if _, err := PricesFromCSV(unkownColumnCSV.Name()); err == nil {
			t.Errorf("A error was expected when loading a csv containing a unkown column")
		}
	}

	//Checking if a error is raised with a csv containing a invalid parameter instead of a number
	invalidNumberCSV := createTempCSV()
	defer os.Remove(invalidNumberCSV.Name())
	invalidNumberCSV.WriteString("Open,High,Low,Close,Volume\n")
	invalidNumberCSV.WriteString("746.25,747.25,746.2,746.95,1045532\n746.95,xpto,746.8,747.05,351191")

	if _, err := PricesFromCSV(invalidNumberCSV.Name()); err == nil ||
		!strings.Contains(err.Error(), "xpto") {
		t.Errorf("A error was expected containing the invalid parameter 'xpto', provided error msg was:\n%v", err.Error())
	}
}

func createTempCSV() *os.File {
	file, err := ioutil.TempFile(".", "prices.*.csv")
	if err != nil {
		log.Fatal(err)
	}
	return file
}
