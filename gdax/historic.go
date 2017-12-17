package gdax

import (
	"math"
	"../common"
	"fmt"
	"log"
	"time"
	
	"strconv"
)

func processFrame(currency string, sframe time.Time, eframe time.Time, gran int) []Record {
	var records []Record

	// Make request
	url := fmt.Sprintf("https://api.gdax.com/products/%s/candles?granularity=%d&start=%s&end=%s",
		currency,
		gran,
		sframe.UTC().Format(time.RFC3339),
		eframe.UTC().Format(time.RFC3339))

	if err := common.SimpleGet(url, &records); err != nil {
		log.Print(err)
		return records
	}

	return records
}

func (g *GDAX) Historic(curr string, startTime time.Time, endTime time.Time, gran int) []Record {
	var records []Record

	requests := math.Ceil(endTime.Sub(startTime).Seconds() / float64(gran))
	if requests > 200 {
		shortDuration := time.Duration(gran) * time.Second
		longDuration := time.Duration(200*gran) * time.Second
		sframe := startTime
		eframe := sframe.Add(longDuration)
		for eframe.Before(endTime) {
			records = append(records, processFrame(curr, sframe, eframe, gran)...)
			sframe = eframe.Add(shortDuration)
			eframe = sframe.Add(longDuration)
			time.Sleep(500 * time.Millisecond)
		}
		if eframe.After(endTime) {
			records = append(records, processFrame(curr, sframe, endTime, gran)...)
		}
	} else {
		records = append(records, processFrame(curr, startTime, endTime, gran)...)
	}

	return records
}


func (g *GDAX) CSV(path string, records []Record) {
	var items [][]string
	for _, n := range records {
		var item []string
		item = append(item, strconv.FormatInt(int64(n.Time), 10))
		item = append(item, strconv.FormatFloat(float64(n.Low), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(n.High), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(n.Open), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(n.Close), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(n.Volume), 'f', -1, 32))
		items = append(items, item)
	}

	if err := common.WriteToCSV(path, items); err != nil {
		log.Print(err)
	}
}