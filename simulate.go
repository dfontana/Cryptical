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

// Portfolio describes the amount of currency and its USD value.
type Portfolio struct {
	Crypto float64
	USD    float64
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
	log.Println(fmt.Sprintf("0: Starting PF %.10f, $%.2f", pf.Crypto, pf.USD))

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
		model.Plot(fmt.Sprintf("./inference%d.png", i+1))

		// 5.   Apply strategy
		trade, err := model.Analyze()
		if err != nil {
			log.Fatal(err)
		}
		switch trade.Type {
		case computation.Buy:
			pf.Crypto += trade.Crypto
			pf.USD -= trade.Crypto * trade.USD
		case computation.Sell:
			pf.Crypto -= trade.Crypto
			pf.USD += trade.Crypto * trade.USD
			// Do nothing under Hodl
		}

		// 6. Pretty Log
		if trade.Type == computation.Hodl {
			if hodlCt != 1 {
				log.Println(fmt.Sprintf("%d: %s. PF: %.10f, $%.2f", i+1, trade.Type, pf.Crypto, pf.USD))
				hodlCt = 1
			}
		} else {
			var fmtStr string
			// if trade.Type == computation.Buy {
			// 	fmtStr = "%d: %s %.10f (-$%.2f) PF: %.10f, $%.2f"
			// } else {
			// 	fmtStr = "%d: %s %.10f (+$%.2f) PF: %.10f, $%.2f"
			// }
			fmtStr = "%d: %s %.10f (@$%.2f) PF: %.10f, $%.2f"
			log.Println(fmt.Sprintf(fmtStr, i+1, trade.Type, trade.Crypto, trade.USD, pf.Crypto, pf.USD))
			hodlCt = 0
		}
	}

	log.Println("Done.")
}
