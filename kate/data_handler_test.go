package kate

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
		open:   746.25,
		high:   747.25,
		low:    746.2,
		close:  746.95,
		volume: 1045532,
	}
	if !reflect.DeepEqual(handler.prices[0], firstPrice) {
		t.Errorf("The first price in the csv is %v the expected was %v", handler.prices[0], firstPrice)
	}

	lastPrice := DataPoint{
		open:   744.4,
		high:   744.6,
		low:    744.4,
		close:  744.6,
		volume: 80597,
	}

	if !reflect.DeepEqual(handler.prices[len(handler.prices)-1], lastPrice) {
		t.Errorf("The last price in the csv is %v the expected was %v", handler.prices[len(handler.prices)-1], lastPrice)
	}

	middlePrice := DataPoint{
		open:   746.4,
		high:   746.45,
		low:    746.4,
		close:  746.45,
		volume: 5306,
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

	//Checking if a error is raised with a csv containing a unknown column
	for _, columnLine := range columnsCSV {
		unknownColumnCSV := createTempCSV()
		defer os.Remove(unknownColumnCSV.Name())
		unknownColumnCSV.WriteString(columnLine)

		if _, err := PricesFromCSV(unknownColumnCSV.Name()); err == nil {
			t.Errorf("A error was expected when loading a csv containing a unknown column")
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
