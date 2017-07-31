package main

// See this playground for using routines + a result channel to aggregate all the JSON data before outputting it to a CSV. https://play.golang.org/p/x7I8ikHFnR

// Also read up on this to learn how to properly do this: http://guzalexander.com/2013/12/06/golang-channels-tutorial.html

import (
    //"os"
    //"encoding/csv"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "time"
    //"strconv"
    "net/http"
    "net/url"
)

var (
    baseUrl = "https://api.gdax.com"
    historic = "/products/ETH-USD/candles"
    startTime = time.Date(2017, time.July, 23, 10 , 0, 0, 0, time.Local)
    endTime = time.Now()
    granularity = 1 //Seconds
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

// Grabs the given time frame from the API, returning the data in array of structs.
func request_frame(sframe time.Time, eframe time.Time) []Record {
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
    if err != nil { fmt.Errorf("Error: ", err) }

    return records
}

// Computes what time frames to request from the API, working within their limits.
func determine_frames(starttime time.Time, endtime time.Time) []time.Time {
    var frames []time.Time

    requests :=  math.Ceil(endtime.Sub(starttime).Seconds() / float64(granularity))
    if requests > 200 {
        shortDuration := time.Duration(granularity) * time.Second
        longDuration := time.Duration(200*granularity) * time.Second
        sframe := starttime
        eframe := sframe.Add(longDuration)
        for eframe.Before(endtime) {
            frames = append(frames, []time.Time{sframe, eframe}...)
            sframe = eframe.Add(shortDuration)
            eframe = sframe.Add(longDuration)
        }
        if eframe.After(endTime) {
            frames = append(frames, []time.Time{sframe, endtime}...)
        }
    }else{
        frames = append(frames, []time.Time{starttime, endtime}...)
    }

    return frames
}

func main() {
    //var records []Record

    frames := determine_frames(startTime, endTime)

    //f, err := os.Create("./output.csv")
    //if err != nil {
    //    fmt.Errorf("Error: ", err)
    //}
    //defer f.Close()

    //w := csv.NewWriter(f)
    //for _, obj := range records {
    //    var item []string
    //    item = append(item, strconv.FormatInt(int64(obj.time),10))
    //    item = append(item, strconv.FormatFloat(float64(obj.low), 'f', -1, 32))
    //    item = append(item, strconv.FormatFloat(float64(obj.high), 'f', -1, 32))
    //    item = append(item, strconv.FormatFloat(float64(obj.open), 'f', -1, 32))
    //    item = append(item, strconv.FormatFloat(float64(obj.close), 'f', -1, 32))
    //    item = append(item, strconv.FormatFloat(float64(obj.volume), 'f', -1, 32))

    //    w.Write(item)
    //}
    //w.Flush()
}
