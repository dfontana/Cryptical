package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/dfontana/Cryptical/computation"
	poloClient "github.com/dfontana/Cryptical/poloniex"
)

// Timeframe describes frame of time starting at date and lasting length, with
// interval increments within it. Ex: starting "12/10/17" at every "300
// seconds" for "24 hours"
type Timeframe struct {
	Date   time.Time
	Period int // in seconds
	Length time.Duration
}

// Simulate will perform a historical simulation of the given model over the given timeframe.
func Simulate(model computation.Computation, pf *Portfolio, tf Timeframe) {
	// 1. Fetch data to simulate over.
	records, err := poloClient.Historic("USDT_ETH", tf.Date, tf.Date.Add(tf.Length), tf.Period)
	if err != nil {
		log.Fatal(err)
	}

	// (used for logging)
	hodlCt := 0
	entry, err := pf.Latest()
	if err != nil {
		log.Fatal("simulation requires a portfolio with at least one entry")
		return
	}
	log.Println(fmt.Sprintf("0: Starting PF %.10f, $%.2f", entry.Crypto, entry.Pair))

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

		// 5.   Apply strategy
		action, err := model.Analyze()
		if err != nil {
			log.Fatal(err)
		}
		switch action.Type {
		case computation.Buy:
			pf.Update(action.Crypto, action.Value)
			model.Plot(fmt.Sprintf("./inference%d.png", i+1))
		case computation.Sell:
			pf.Update(-action.Crypto, action.Value)
			model.Plot(fmt.Sprintf("./inference%d.png", i+1))
		}

		// 6. Pretty Log
		latest, _ := pf.Latest()
		if action.Type == computation.Hodl {
			if hodlCt != 1 {
				log.Println(fmt.Sprintf("%d: %s (@$%.2f) PF: %.10f, $%.2f", i+1, action.Type, latest.CurrentValue, latest.Crypto, latest.Pair))
				hodlCt = 1
			}
		} else {
			log.Println(fmt.Sprintf("%d: %s %.10f (@$%.2f) PF: %.10f, $%.2f", i+1, action.Type, action.Crypto, action.Value, latest.Crypto, latest.Pair))
			hodlCt = 0
		}
	}

	log.Println("Done.")
}
