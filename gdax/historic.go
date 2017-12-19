package gdax

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"time"
	"sort"

	"../common"
)

// Historic returns data in interval of gran (in seconds), for the specified currency pair, curr.
// This history is bounded by the start and endtime stamps provided. The result is a slice, which
// has been sorted in ascending (oldest first) order to guarentee order (since API does not)
func (g *GDAX) Historic(curr string, startTime time.Time, endTime time.Time, gran int) []Record {
	var records []Record

	requests := math.Ceil(endTime.Sub(startTime).Seconds() / float64(gran))
	if requests > 200 {
		frameLen := time.Duration(200*gran) * time.Second
		sframe := startTime
		eframe := startTime.Add(frameLen)
		for eframe.Before(endTime) {
			// Request the frame, move forward 1 frame (no overlap), wait 500ms to prevent lockout
			records = append(records, processFrame(curr, sframe, eframe, gran)...)
			sframe = eframe.Add(time.Duration(gran) * time.Second)
			eframe = sframe.Add(frameLen)
			time.Sleep(500 * time.Millisecond)
		}
		if eframe.After(endTime) {
			// The frame extends over the desired end boundary, so fill in
			records = append(records, processFrame(curr, sframe, endTime, gran)...)
		}
	} else {
		records = processFrame(curr, startTime, endTime, gran)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Time < records[j].Time
	})

	return records
}


func processFrame(currency string, sframe time.Time, eframe time.Time, gran int) []Record {
	var records []Record

	// Make request
	values := url.Values{}
	values.Set("start", sframe.Format(time.RFC822Z))
	values.Set("end", eframe.Format(time.RFC822Z))
	values.Set("granularity", strconv.Itoa(gran))
	fmtUrl := fmt.Sprintf("https://api.gdax.com/products/%s/candles?", currency) + values.Encode()

	if err := common.SimpleGet(fmtUrl, &records); err != nil {
		log.Println(err)
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
