package poloniex

import (
	"../common"
	"log"
	"time"
	"strconv"
	"net/http"
	"encoding/json"
	"fmt"
)

const (
	POLONIEX_API = "https://poloniex.com/public"
	CHART_DATA = "returnChartData"
)

// Historic returns data for the hiven currency between the given times. The
// data's interval is specified with the gran arugment - the number of seconds
// between data points. Valid entries are 300, 900, 1800, 7200, 14400, or 86400.
// Invalid entries will result in errors unmarshalling the response.
func Historic(curr string, startTime time.Time, endTime time.Time, gran int) (retVal []CandleStick, err error) {
	url := fmt.Sprintf(
		"%s?command=%s&currencyPair=%s&period=%d&start=%d&end=%d",
		POLONIEX_API,
		CHART_DATA,
		curr,
		gran,
		startTime.Unix(),
		endTime.Unix(),
	)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&retVal)
	if err != nil {
		log.Println("FUDGE")
	}
	return
}

// CSV creates a csv at the given path consisting of the given candlestick data.
func CSV(path string, records []CandleStick) {
	items := make(chan []string)
	errors := make(chan error)

	go common.WriteToCSV(path, items, errors)

	for _, obj := range records {
		select {
			case err := <-errors:
				log.Print(err)
				break; // Out of loop
			default:
				//Send next item
				var item []string
				item = append(item, obj.Date.Time.Format(time.RFC822Z))
				item = append(item, strconv.FormatFloat(float64(obj.High), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Low), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Open), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Close), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Volume), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.QuoteVolume), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.WeightedAverage), 'f', -1, 32))
				items <- item
		}		
	}
	close(items)
	<-errors
}