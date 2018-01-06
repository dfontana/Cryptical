package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/dfontana/Cryptical/computation"
	poloClient "github.com/dfontana/Cryptical/poloniex"
)

// Simulate will perform a historical simulation of the given model for the given day, lasting
// for the given duration with data retrieved in the given interval. Ideally, interval should match
// the interval the model was computed at to obtain meaningful results.
func Simulate(model computation.Computation, date time.Time, interval int, duration time.Duration) {
	// 1. Fetch data to simulate over.
	records, err := poloClient.Historic("USDT_ETH", date, date.Add(duration), interval)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Loop over each entry in the new data:
	for i, val := range records {
		// 3.   Append the entry to the old data
		point := computation.TimeSeries{
			val.Date.Time,
			val.High,
		}
		if err := model.AddPoint(point); err != nil {
			log.Fatal(err)
		}

		// 4.   Recompute model
		if err := model.Compute(); err != nil {
			log.Fatal(err)
		}

		// 5.   Apply strategy, log result
		trade, err := model.Analyze()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(fmt.Sprintf("%d: %s %.5f for $%.2f", i, trade.Type, trade.Crypto, trade.USD))
	}

	log.Println("Done.")
}
