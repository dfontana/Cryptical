package poloniex

import (
	"log"
	"fmt"
	"time"
	"strconv"

	// Get requests
	"net/http"
	"encoding/json"
	"io/ioutil"

	//CSV creation
	"encoding/csv"
	"os"
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

	log.Print(url)

	if err := simpleGet(url, &records); err != nil {
		log.Print(err)
		return records
	}

	return records
}

func (p *Poloniex) CSV(path string, records []Record) {
	items := make(chan []string)
	errors := make(chan error)

	go common.writeToCSV(path, items, errors)

	for _, obj := range records {
		select {
			case err := <-errors:
				log.Print(err)
				break; // Out of loop
			default:
				//Send next item
				var item []string
				item = append(item, strconv.FormatInt(int64(obj.GlobalTradeID), 10))
				item = append(item, strconv.FormatInt(int64(obj.TradeID), 10))
				item = append(item, obj.Date)
				item = append(item, obj.Type)
				item = append(item, strconv.FormatFloat(float64(obj.Rate), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Amount), 'f', -1, 32))
				item = append(item, strconv.FormatFloat(float64(obj.Total), 'f', -1, 32))
				items <- item
		}		
	}
	close(items)
	<-errors
}

func simpleGet(url string, into interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		err, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("SimpleGet returned status %d, %s", res.StatusCode, string(err))
	}

	defer res.Body.Close()

	// Parse the response
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		return err
	}

	return nil
}

// writeToCSV is a routine that will write incoming items to a CSV
// at the given path. Should an error occur, it is sent into the given
// error channel and the routine terminates.
func writeToCSV(path string, items chan []string, errors chan error) {
	f, err := os.Create(path)

	// Terminate early, sending our error to caller channel
	if err != nil {
		errors <- err
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)
	for item := range items {
		w.Write(item)
	}
	w.Flush()
	close(errors)
}
