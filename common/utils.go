package common

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	//CSV creation
	"encoding/csv"
	"os"
)

func SimpleGet(url string, into interface{}) error {
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
