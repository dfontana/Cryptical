package poloniex

import (
	"../common"
	"log"
	"fmt"
	"time"
	"strconv"
)

const (
	BASE 			= "https://poloniex.com/public?command="
	HISTORY 	= "returnTradeHistory"
)

// Max of 50,000 records per request, so requests one day at a time to ensure most data points returned.
func (p *Poloniex) Historic(curr string, startTime time.Time, endTime time.Time) []Record {
	var records []Record
	
	sframe := startTime

	for sframe.Before(endTime) {
		nextRequest := time.Now().Add(167 * time.Millisecond)

		eframe := sframe.Add(1 * time.Hour * 24)
		records = append(records, processFrame(curr, sframe, eframe)...)
		sframe = eframe

		// Sleeps the remainder of the duration to meet rate limit
		diff := nextRequest.Sub(time.Now())
		if(diff > 0){
			time.Sleep(diff)
		}
	}

	return records
}

func processFrame(currency string, sframe time.Time, eframe time.Time) []Record {
	records := []Record{}

	// Make request
	url := fmt.Sprintf("%s%s&currencyPair=%s&start=%d&end=%d",
		BASE,
		HISTORY,
		currency,
		sframe.Unix(),
		eframe.Unix())

	if err := common.SimpleGet(url, &records); err != nil {
		log.Print(err)
		return records
	}

	return records
}

func (p *Poloniex) CSV(path string, records []Record) {
	var items [][]string
	for _, obj := range records {
		var item []string
		item = append(item, strconv.FormatInt(int64(obj.GlobalTradeID), 10))
		item = append(item, strconv.FormatInt(int64(obj.TradeID), 10))
		item = append(item, obj.Date)
		item = append(item, obj.Type)
		item = append(item, strconv.FormatFloat(float64(obj.Rate), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.Amount), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.Total), 'f', -1, 32))

		items = append(items, item)
	}

	if err := common.WriteToCSV(path, items); err != nil {
		log.Print(err)
	}
}