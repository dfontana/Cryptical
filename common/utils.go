package common

import (
	"encoding/json"
	"net/http"
	"fmt"

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
		return fmt.Errorf("SimleGet returned status %d", res.StatusCode)
	}

	defer res.Body.Close()

	// Parse the response
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		return err
	}

	return nil
}

func WriteToCSV(path string, records [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, item := range records {
		w.Write(item)
	}
	w.Flush()

	return nil
}