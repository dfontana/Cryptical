package gdax

import (
	gdax "github.com/preichenberger/go-gdax"
	"log"
	"math"
	"time"
	"sort"
	"strconv"

	"../common"
)

// Historic returns data in interval of gran (in seconds), for the specified
// currency pair, curr. This history is bounded by the start and endtime stamps
// provided. The result is a slice, which has been sorted in ascending (oldest
// first) order to guarentee order (since API does not)
func Historic(curr string, startTime time.Time, endTime time.Time, gran int) []gdax.HistoricRate {
	//Build client
	client := gdax.NewClient("", "", "")

	// Holds all our return data
	var records []gdax.HistoricRate

	// Since this API is limited to 200 returned results per request and 6 calls
	// per second, break up the interval into smaller calls with a sleep
	numExpected := math.Ceil(endTime.Sub(startTime).Seconds() / float64(gran)) + 1
	if numExpected > 200 {
		frameLen := time.Duration(200*gran) * time.Second
		sframe := startTime
		eframe := startTime.Add(frameLen)
		for eframe.Before(endTime) {
			params := gdax.GetHistoricRatesParams { sframe, eframe, gran }
			records = append(records, processFrame(client, curr, params)...)
			sframe = eframe.Add(time.Duration(gran) * time.Second)
			eframe = sframe.Add(frameLen)
			time.Sleep(500 * time.Millisecond)
		}
		if eframe.After(endTime) {
			// The frame extends over the desired end boundary, so fill in
			params := gdax.GetHistoricRatesParams { sframe, endTime, gran }
			records = append(records, processFrame(client, curr, params)...)
		}
	} else {
		// Don't need to break up the call
		params := gdax.GetHistoricRatesParams { startTime, endTime, gran }
		records = processFrame(client, curr, params)
	}

	// Ensures data is returned in historical order (oldest first)
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.Before(records[j].Time)
	})

	return records
}

// Handles a slice in history. On error returns an empty slice, logging the error.
func processFrame(client *gdax.Client, currency string, params gdax.GetHistoricRatesParams) []gdax.HistoricRate {
	rates, err := client.GetHistoricRates(currency, params)
	if err != nil {
		log.Println(err)
	}
	return rates
}

// CSV creates a CSV saved at the given path made from the given array of
// historical rate structs.
func CSV(path string, records []gdax.HistoricRate) {
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
			item = append(item, strconv.FormatInt(obj.Time.Unix(), 10))
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
