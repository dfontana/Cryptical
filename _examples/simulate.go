package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dfontana/Cryptical/computation"
	poloClient "github.com/dfontana/Cryptical/poloniex"
)

// Just for the time being while testing
func main() {

	// The below code is the outline of a simulation one might write on historic data
	// to evaluate their trading model strategy. Here, we evaluate an MACD model.

	/** ============================================
	  =============== MODEL PREP STEP ================
	  ============================================ **/
	// Day in history I'd like to infer for, -1 day since Historic is inclusive)
	// So if I want to infer for the 21st, I must ask for up to the 20th's of data
	endDate := time.Date(2017, time.December, 20, 23, 59, 59, 0, time.Local)

	// These are in 5 min periods meaining we look back 12, 5, and 3.5 hours
	slow := 24 * 6
	fast := 10 * 6
	sign := 7 * 6
	gran := 300 // Seconds

	// Will probably want ~3 * Slow data points to ensure we have enough.
	startHist := endDate.Add(time.Duration(-3*slow*gran) * time.Second)
	records, err := poloClient.Historic("USDT_ETH", startHist, endDate, gran)
	if err != nil {
		log.Fatal(err)
	}

	// Convert said data to the needed format to run our MACD computation.
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Date.Time,
			val.High,
		}
	}

	// Compute the inital model and plot it for visual inspection
	comp := computation.MACD{}
	if err := comp.Compute(hist, fast, slow, sign); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./inf/inference0.png")

	/** ============================================
		  =============== STRATEGY DEFINITION STEP =======
	    ============================================ **/

	strategy := func(m computation.MACD) computation.Trade {
		// This is just an example and will perform really poorly
		// Look at last entry and take action
		val := m.Hist[len(m.Hist)-1]
		var action computation.Trade
		switch {
		case val > 2:
			action = computation.Trade{
				computation.Sell,
				0.05,  // amount in crypto
				200.0, // amount in USD
			}
		case val < -2:
			action = computation.Trade{
				computation.Buy,
				0.05,
				200.0,
			}
		default:
			action = computation.Trade{
				computation.Hodl,
				0.0,
				0.0,
			}
		}
		return action
	}

	/** ============================================
	  =============== SIMULATION STEP ================
	  ============================================ **/
	// 1. Fetch historic data from endDate onwards (this wil be 1 day's worth)
	records, err = poloClient.Historic("USDT_ETH", endDate, endDate.Add(24*time.Hour), gran)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Loop over each entry in the new data:
	sim := hist[:]
	for i, val := range records {
		// 3.   Append the entry to the old data
		sim = append(sim, computation.TimeSeries{
			val.Date.Time,
			val.High,
		})

		// 4.   Recompute model & plot it
		if err := comp.Compute(sim, fast, slow, sign); err != nil {
			log.Fatal(err)
		}
		comp.Plot(fmt.Sprintf("./inf/inference%d.png", i))

		// 5.   Apply strategy, log result
		trade := strategy(comp)
		log.Println(fmt.Sprintf("%d: %s %.5f for $%.2f", i, trade.Type, trade.Crypto, trade.USD))
	}

	log.Println("Done.")
}
