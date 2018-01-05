# Cryptical
[![GoDoc](https://godoc.org/github.com/dfontana/Cryptical?status.svg)](https://godoc.org/github.com/dfontana/Cryptical)

## TODO:
- Test cases (for specific time periods from historic data, computation calculations, and simulation for historic period/model)
- Add fee retrieval for each exchange's package

Enhance the Simulation:
- Account for fees in each Trade (integrate into the Trade struct)
- Track number of trades made
- Track average trade size
- Track fees incurred
- Should portfolio be tracking sum, or should it really be tracking individual orders
- May need lightweight database backing for portfolio tracking (when history is large).
 
Enhance and Expand Exchanges:
- Binance (github.com/adshao/go-binance)
- Add authenticated endpoints for submitting trades

Automate:
- Adapt simulation code for realtime trading: you'll need to fetch historic data at startup, then work in realtime from there.
- Integrate database for book-keeping (if/when needed?)

Interface:
- Work with ReactJS (for practice) and D3JS to create a web interface served by the bot.
 - Start with running & displaying stats / logs for a simulation
 - Then display an interactive chart for that simulation (trade data + indicators)
 - Then add a livefeed page to display chart from WSS feed
 - Then consider adding indicators to livefeed
 - And finally, adding monitoring dashboard for bot activity (chart + indicators + bot's decisions)

## How to run:
See `_examples/` for ways to work with the API. Be sure to see the GoDoc, link provided in the badge at the top (and explore the subdirectories). To run these examples just `cd` into their folder and run `go run <file>.go`. Be sure you have already installed the dependencies with `go get`.

## Dependencies:
- github.com/gorilla/websocket
- github.com/preichenberger/go-gdax
- github.com/wcharczuk/go-chart
- github.com/joho/godotenv

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot
