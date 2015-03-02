package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"math"
)

func main() {

	var from date
	var to date

	from.day = 1
	from.month = 1
	from.year = 2014		

	to.day = 1
	to.month = 1
	to.year = 2015

/*	
    goog := getOpenPoints("GOOG", from, to)
	uhs := getOpenPoints("UHS", from, to)
	aapl := getOpenPoints("AAPL", from, to)
		
	displayPlot("Google & USH", goog, uhs, aapl)*

	goog := getPattern("goog", from, to, 4)

	fmt.Printf("%d %d\n", len(goog), len(goog[0]))

	for i := range goog {
		for j := range goog[0] {
			fmt.Printf("%f\t ", goog[i][j])
		}
		fmt.Println()
	}
*/

	sym := []string{"UHS", "GOOG", "AAPL"}
	getPatterns(sym, from, to, 10)

	
	
}

func getPatterns(symbol []string, from date, to date, patternSize int) [][]float32 {

	patternArrays := make([][][]float32, len(symbol))

	numberOfTotalPatterns := 0
	
	for i := range symbol {
		patternArrays[i] = getPattern(symbol[i], from, to, patternSize)
		numberOfTotalPatterns += len(patternArrays[i])
	}

	ret := make([][]float32, numberOfTotalPatterns)

	index := 0
	
	for stock := range patternArrays {
		for _, pat := range patternArrays[stock] {
			ret[index] = pat
			index++
		}
	}

	fmt.Printf("%d %d", len(ret), numberOfTotalPatterns)

	return ret
	

	
}

func getPattern(symbol string, from date, to date, patternSize int) [][]float32 {


	data := getPercentChangeData(symbol, from, to)
	
	pattern := make([][]float32, len(data) - (patternSize - 1))
	
	for i := range pattern {

		pattern[i] = make([]float32, patternSize)		
	}

	index := 0
	
	for i := range pattern {
		for j := 0; j < patternSize; j++ {
			pattern[i][j] = data[index + j]
		}
		index++
	}

	return pattern
}

func getPercentChangeData(symbol string, from date, to date) []float32 {

	arr := getStockInfo(symbol, from, to)
	a := make([]float32, len(arr))
	for i := range a {
		a[i] = arr[i].open
	}
	
	return toPercentChangeArray(a)	
}

func getOpenPoints(symbol string, from date, to date) plotter.XYs {
	pct := getPercentChangeData(symbol, from, to)
	points := getXYs(pct)
	return points
}


func getXYs(a []float32) plotter.XYs {
	points := make(plotter.XYs, len(a))
	for i := range a {
		points[i].X = float64(i)
		points[i].Y = float64(a[i])
	}
	return points
}


func getStockInfo(symbol string, from date, to date) []record {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=d",
		symbol,
		from.month - 1, from.day, from.year,
	    to.month - 1, to.day, to.year)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error making an HTTP request for stock %s.\n", symbol)
		return nil
	}

	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err2 := csvReader.ReadAll()

    if err2 != nil {
		fmt.Printf("Error parsing CSV values for stock %s.\n", symbol)
		return nil
	}

	ret := make([]record, len(records) - 1)

	var temp record
	for i := 1; i < len(records); i++ {
		temp.date = getDate(records[i][0])

		n, _:= strconv.ParseFloat(records[i][1], 32)
		temp.open = float32(n)
		
		n, _= strconv.ParseFloat(records[i][2], 32)
		temp.high = float32(n)

		n, _ = strconv.ParseFloat(records[i][3], 32)
		temp.low = float32(n)

		n, _ = strconv.ParseFloat(records[i][4], 32)
		temp.close = float32(n)

		nn, _ := strconv.ParseInt(records[i][5], 10, 32)
		temp.volume = uint32(nn)
		
		n, _ = strconv.ParseFloat(records[i][6], 32)
		temp.adjClose = float32(n)
		
		ret[i-1] = temp
	} 
		
	return ret
}


func toPercentChangeArray(a []float32) (ret []float32) {

	ret = make([]float32, len(a))

	for i:=0; i < len(a) - 1; i++ {
		ret[i] = percentChange(a[i], a[i+1])
	}

	return
}


func percentChange(old float32, new float32) (ret float32) {
	return float32(math.Abs(float64(((old-new)/old)*100.0)))
}


func displayPlot(name string, toPlot ...plotter.XYs) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	
	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"

	for i, set := range toPlot {
		err = plotutil.AddLines(p, string(i), set)
	}
	
	if err != nil {
		panic(err)
	}
	
	// Save the plot to a PNG file.
	if err := p.Save(4, 4, name + ".png"); err != nil {
		panic(err)
	}
}


func getDate(s string) (ret date) {

	a := strings.Split(s, "-")

	n, _ := strconv.ParseInt(a[0], 10, 32)
	ret.month = int(n)

	n, _ = strconv.ParseInt(a[1], 10, 32)
	ret.day = int(n)

	n, _= strconv.ParseInt(a[2], 10, 32)
	ret.year = int(n)
	
	return
}


func printTable(arr []record) {
	fmt.Println("Date \t\t Open \t\t High \t\t Low \t\t Close \t\t Volume \t adjClose")
	for i:=0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}

//*****************************************************************************************//

type date struct {
	day int
	month int
	year int
}


func (d date) String() string {
	return fmt.Sprintf("%d/%d/%d", d.day, d.month, d.year)
}


//*****************************************************************************************//


type record struct {

	date date
	open float32
	high float32
	low float32
	close float32
	volume uint32
	adjClose float32
}


func (r record) String() string {

	return fmt.Sprintf("%s \t %.6f \t %.6f \t %.6f \t %.6f \t %d \t %.6f", r.date, r.open, r.high, r.low, r.close, r.volume, r.adjClose)
}
