package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math"
    "time"
    "net/http"
    "net/url"
)

var (
    baseUrl = "https://api.gdax.com"
    historic = "/products/ETH-USD/candles"
    startTime = time.Date(2017, time.July, 23, 10 , 0, 0, 0, time.Local)
    endTime = time.Now()
    granularity = 30 //Seconds
)


type Record struct {
    time    int     `json:",string"` //Start time
    low     float32 `json:",string"` //Low price for time
    high    float32 `json:",string"` //High price for time
    open    float32 `json:",string"` //First trade price
    close   float32 `json:",string"` //Last trade price
    volume  float32 `json:",string"` //Trading activity volume
}

// Custom Unmarshalling to handle array properly.
func (n *Record) UnmarshalJSON(buf []byte) error {
    tmp := []interface{}{&n.time, &n.low, &n.high, &n.open, &n.close, &n.volume}
    wantLen := len(tmp)
    if err := json.Unmarshal(buf, &tmp); err != nil { return err }
    if len(tmp) != wantLen { return fmt.Errorf("Size mismatch in record") }
    return nil
}

func process_frame(sframe time.Time, eframe time.Time) []Record {
    url := fmt.Sprintf("%s%s?granularity=%d&start=%s&end=%s",
        baseUrl,
        historic,
        granularity,
        url.QueryEscape(sframe.String()),
        url.QueryEscape(eframe.String()))

    res, err := http.Get(url)
    if err != nil {
        log.Fatal("Failed Get: ", err)
        return nil
    }
    defer res.Body.Close()

    var records []Record
    err = json.NewDecoder(res.Body).Decode(&records)
    if err != nil { log.Fatal(err) }

    return records
}

func main() {
    var records []Record

    requests :=  math.Ceil(endTime.Sub(startTime).Seconds() / float64(granularity))
    if requests > 200 {
        shortDuration := time.Duration(granularity) * time.Second
        longDuration := time.Duration(200*granularity) * time.Second
        sframe := startTime
        eframe := sframe.Add(longDuration)
        for eframe.Before(endTime) {
            records = append(records, process_frame(sframe, eframe)...)
            sframe = eframe.Add(shortDuration)
            eframe = sframe.Add(longDuration)
        }
        if eframe.After(endTime) {
            records = append(records, process_frame(sframe, endTime)...)
        }
    }else{
        records = append(records, process_frame(startTime, endTime)...)
    }

    fmt.Println(records)
}
