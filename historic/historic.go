package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

var (
	baseURL     = "https://api.gdax.com"
	historic    = "/products/ETH-USD/candles"
	startTime   = time.Date(2017, time.December, 1, 0, 0, 0, 0, time.Local)
	endTime     = time.Now()
	granularity = 600 //Seconds
)

type record struct {
	Time   int     `json:",string"` //Start time
	Low    float32 `json:",string"` //Low price for time
	High   float32 `json:",string"` //High price for time
	Open   float32 `json:",string"` //First trade price
	Close  float32 `json:",string"` //Last trade price
	Volume float32 `json:",string"` //Trading activity volume
}

// UnmarshalJSON handles decomposing the returned array from GDAX into a series of record structs.
func (n *record) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Time, &n.Low, &n.High, &n.Open, &n.Close, &n.Volume}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if len(tmp) != wantLen {
		return errors.New("Size mismatch in record")
	}
	return nil
}

func processFrame(sframe time.Time, eframe time.Time) []record {
	url := fmt.Sprintf("%s%s?granularity=%d&start=%s&end=%s",
		baseURL,
		historic,
		granularity,
		sframe.UTC().Format(time.RFC3339),
		eframe.UTC().Format(time.RFC3339))

	res, err := http.Get(url)
	if err != nil {
		log.Print("Failed Get: ", err)
		return nil
	}

	defer res.Body.Close()

	var records []record
	if err := json.NewDecoder(res.Body).Decode(&records); err != nil {
		log.Print(err, " Are you over the request limit?")
	}

	return records
}

func main() {
	var records []record

	requests := math.Ceil(endTime.Sub(startTime).Seconds() / float64(granularity))
	if requests > 200 {
		shortDuration := time.Duration(granularity) * time.Second
		longDuration := time.Duration(200*granularity) * time.Second
		sframe := startTime
		eframe := sframe.Add(longDuration)
		for eframe.Before(endTime) {
			records = append(records, processFrame(sframe, eframe)...)
			sframe = eframe.Add(shortDuration)
			eframe = sframe.Add(longDuration)
			time.Sleep(500 * time.Millisecond)
		}
		if eframe.After(endTime) {
			records = append(records, processFrame(sframe, endTime)...)
		}
	} else {
		records = append(records, processFrame(startTime, endTime)...)
	}

	f, err := os.Create("./output.csv")
	if err != nil {
		log.Print(err)
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
