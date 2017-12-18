package gdax

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"time"

	"../common"
)

func processFrame(currency string, sframe time.Time, eframe time.Time, gran int) []Record {
	var records []Record

	// Make request
	values := url.Values{}
	sUTC := sframe
	eUTC := eframe
	sstring := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", sUTC.Year(), sUTC.Month(), sUTC.Day(), sUTC.Hour(), sUTC.Minute(), sUTC.Second())
	estring := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", eUTC.Year(), eUTC.Month(), eUTC.Day(), eUTC.Hour(), eUTC.Minute(), eUTC.Second())
	values.Set("start", sstring)
	values.Set("end", estring)
	values.Set("granularity", strconv.Itoa(gran))
	fmtUrl := fmt.Sprintf("https://api.gdax.com/products/%s/candles?", currency) + values.Encode()

	if err := common.SimpleGet(fmtUrl, &records); err != nil {
		log.Println(err)
	}
	log.Println(records)
	return nil
}

// Historic returns data in interval of gran (in seconds), for the specified currency pair, curr.
// This history is bounded by the start and endtime stamps provided. The result is an array.
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
	items := make(chan []string)
	errors := make(chan error)

	go common.WriteToCSV(path, items, errors)

	for _, obj := range records {
		select {
		case err := <-errors:
			log.Print(err)
			break // Out of loop
		default:
			//Send next item
			var item []string
			item = append(item, strconv.FormatInt(int64(obj.Time), 10))
			item = append(item, strconv.FormatFloat(float64(obj.Low), 'f', -1, 32))
			item = append(item, strconv.FormatFloat(float64(obj.High), 'f', -1, 32))
			item = append(item, strconv.FormatFloat(float64(obj.Open), 'f', -1, 32))
			item = append(item, strconv.FormatFloat(float64(obj.Close), 'f', -1, 32))
			item = append(item, strconv.FormatFloat(float64(obj.Volume), 'f', -1, 32))
			items <- item
		}
	}
	close(items)
	<-errors
}
