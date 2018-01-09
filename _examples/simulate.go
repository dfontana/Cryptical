package main

import (
	"log"
	"time"

	"github.com/dfontana/Cryptical"
	"github.com/dfontana/Cryptical/computation"
	"github.com/dfontana/Cryptical/poloniex"
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
	slow := 26
	fast := 12
	sign := 9
	gran := 300 // Seconds

	// Will probably want ~3 * Slow data points to ensure we have enough.
	startHist := endDate.Add(time.Duration(-3*slow*gran) * time.Second)
	records, err := poloniex.Historic("USDT_ETH", startHist, endDate, gran)
	if err != nil {
		log.Fatal(err)
	}

	// Convert said data to the needed format to run our MACD computation.
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Date.Time,
			val.WeightedAverage,
		}
	}

	// Compute the inital model and plot it for visual inspection
	comp := computation.MACD{
		Data: hist,
		Fast: fast,
		Slow: slow,
		Sign: sign,
	}
	if err := comp.Compute(); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./inference0.png")

	/** ============================================
	=============== STRATEGY DEFINITION STEP =======
	============================================ **/

	pf := bot.Portfolio{CurrencyPair: "ETH-USD"}
	pf.Initialize(0.0, 500.0, 800.0)

	var lastAction computation.Action
	strategy := func(m *computation.MACD) computation.Action {
		// This is just an example and will perform really poorly
		// Look at last entry and take action

		action := computation.Action{
			computation.Hodl,
			0.0,
			0.0,
		}

		latest, _ := pf.Latest()

		prevValM := m.MACD[len(m.Hist)-2]
		todaValM := m.MACD[len(m.Hist)-1]

		prevValS := m.Signal[len(m.Hist)-2]
		todaValS := m.Signal[len(m.Hist)-1]

		// MACD goes above signal -> buy
		if prevValS > prevValM && todaValM > todaValS {
			// Dont repeat actions
			if lastAction.Type != computation.Buy {
				usdVal := 0.75 * latest.Pair // Buy with 75% your funds
				if usdVal < 0.3 {
					// Too small, bail.
					action = computation.Action{
						computation.Hodl,
						0.0,
						0.0,
					}
				} else {
					action = computation.Action{
						computation.Buy,
						usdVal / m.Data[len(m.Data)-1].Data,
						m.Data[len(m.Data)-1].Data,
					}
				}
				lastAction = action
			}
			// MACD goes below signal -> sell
		} else if prevValS < prevValM && todaValM < todaValS {
			// Dont repeat actions
			if lastAction.Type != computation.Sell {
				// Sell 75% of current portfolio
				sellAmt := 0.75 * latest.Crypto
				action = computation.Action{
					computation.Sell,
					sellAmt,
					m.Data[len(m.Data)-1].Data,
				}
				lastAction = action
			}
		}
		// otherwise we hodl
		return action
	}
	comp.Strategy = strategy

	/** ============================================
	=============== SIMULATION STEP ================
	============================================ **/
	tf := bot.Timeframe{endDate, gran, 24 * time.Hour}
	bot.Simulate(&comp, &pf, tf)
}
