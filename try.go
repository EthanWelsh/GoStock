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

	arr := getStockInfo("GOOG", from, to)

	a := make(plotter.XYs, len(arr))
	for i := range a {
		a[i].X = float64(i)
		a[i].Y = float64(arr[i].open)
	}

	b := make(plotter.XYs, len(arr))
	for i := range b {
		b[i].X = float64(i)
		b[i].Y = float64(arr[i].close)
	}

	displayPlot("GOOG Open Close Chart", a, b)
	
	//printTable(arr)
	
	
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
