package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

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

func main() {

	var from date
	var to date

	from.day = 10
	from.month = 4
	from.year = 2014		

	to.day = 15
	to.month = 4
	to.year = 2014

	arr := getStockInfo("APL", from, to)

	for i:=1; i < len(arr); i++ {
		fmt.Println(arr[i])
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

func getStockInfo(symbol string, from date, to date) []record {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=d",
		symbol,
		from.month - 1, from.day, from.year,
	    to.month - 1, to.day, to.year)

	fmt.Println(url)
	
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

	ret := make([]record, len(records))

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
		
		ret[i] = temp
	} 
		
	return ret
}
