package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	//CSV creation
	"encoding/csv"
	"os"
)

func SimpleGet(url2 string, into interface{}) error {
	// res, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }

	// if res.StatusCode != 200 {
	// 	return fmt.Errorf("SimpleGet returned status %d", res.StatusCode)
	// }

	// defer res.Body.Close()

	// // Parse the response
	// if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
	// 	return err
	// }

	// return nil
	var response *http.Response
	var request *http.Request

	// values := url.Values{}
	// values.Set("start", "2017-11-18T07:03:18Z")
	// values.Set("end", "2017-12-18T07:03:18Z")
	// values.Set("granularity", "86400")
	// url2 = "https://api.gdax.com/products/ETH-USD/candles?" + values.Encode()
	log.Println(url2)
	request, err := http.NewRequest("GET", url2, nil)
	if err == nil {
		request.Header.Add("User-Agent", "curl/7.54.0")
		request.Header.Add("Accept", "*/*")
		debug(httputil.DumpRequestOut(request, true))
		response, err = (&http.Client{}).Do(request)
	}

	if err == nil {
		defer response.Body.Close()
		debug(httputil.DumpResponse(response, true))
		_, err = ioutil.ReadAll(response.Body)
	}

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	return nil
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// WriteToCSV is a routine that will write incoming items to a CSV
// at the given path. Should an error occur, it is sent into the given
// error channel and the routine terminates.
func WriteToCSV(path string, items chan []string, errors chan error) {
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
