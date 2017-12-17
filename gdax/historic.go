package gdax

import (
	"log"
	"math"
	
	//API Request
	"encoding/json"
	"net/http"
	"fmt"
	"time"

	//CSV creation
	"encoding/csv"
	"os"
	"strconv"
)

const baseURL = "https://api.gdax.com/products/ETH-USD/candles"

// UnmarshalJSON handles decomposing the returned array from GDAX into a series of record structs.
func (n *Record) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Time, &n.Low, &n.High, &n.Open, &n.Close, &n.Volume}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	return nil
}

func processFrame(sframe time.Time, eframe time.Time, gran int) []Record {
	var records []Record

	// Make request
	url := fmt.Sprintf("%s?granularity=%d&start=%s&end=%s",
		baseURL,
		gran,
		sframe.UTC().Format(time.RFC3339),
		eframe.UTC().Format(time.RFC3339))

	res, err := http.Get(url)
	if err != nil {
		log.Print("Failed Get: ", err)
		return records
	}

	defer res.Body.Close()

	// Parse the response
	if err := json.NewDecoder(res.Body).Decode(&records); err != nil {
		log.Print(err, " Are you over the request limit?")
		return records
	}

	return records
}

func (g *GDAX) Historic(startTime time.Time, endTime time.Time, gran int) []Record {
	var records []Record

	requests := math.Ceil(endTime.Sub(startTime).Seconds() / float64(gran))
	if requests > 200 {
		shortDuration := time.Duration(gran) * time.Second
		longDuration := time.Duration(200*gran) * time.Second
		sframe := startTime
		eframe := sframe.Add(longDuration)
		for eframe.Before(endTime) {
			records = append(records, processFrame(sframe, eframe, gran)...)
			sframe = eframe.Add(shortDuration)
			eframe = sframe.Add(longDuration)
			time.Sleep(500 * time.Millisecond)
		}
		if eframe.After(endTime) {
			records = append(records, processFrame(sframe, endTime, gran)...)
		}
	} else {
		records = append(records, processFrame(startTime, endTime, gran)...)
	}

	return records
}


func (g *GDAX) writeToCSV(records []Record) {
	f, err := os.Create("./output.csv")
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, obj := range records {
		var item []string
		item = append(item, strconv.FormatInt(int64(obj.Time), 10))
		item = append(item, strconv.FormatFloat(float64(obj.Low), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.High), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.Open), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.Close), 'f', -1, 32))
		item = append(item, strconv.FormatFloat(float64(obj.Volume), 'f', -1, 32))

		w.Write(item)
	}
	w.Flush()
}