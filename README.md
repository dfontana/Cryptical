# Cryptical

[![GoDoc](https://godoc.org/github.com/dfontana/Cryptical?status.svg)](https://godoc.org/github.com/dfontana/Cryptical)

## TODO

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
Work with ReactJS (for practice) to create a web interface served by the bot.

- Routing should be such that: `/` leads to home page and interface, while `/api/` leads to data sources.
- Design the layout & use cases of the simulation page
- Implement the simulation page in chunks.
- Move onto livefeed data, chartings, and bot monitoring

## How to run

See `_examples/` for ways to work with the API. Be sure to see the GoDoc, link provided in the badge at the top (and explore the subdirectories). To run these examples just `cd` into their folder and run `go run <file>.go`. Be sure you have already installed the dependencies with `go get`.

To run the frontend, you'll need to:

- `cd` into the `frontend` folder.
- Run: `npm install` to install dependencies
- Run: `npm run build` to compile the site into static files for serving.

Then you may proceed to start the server (`go run server.go`).

## Dependencies

- github.com/gorilla/websocket
- github.com/preichenberger/go-gdax
- github.com/wcharczuk/go-chart
- github.com/joho/godotenv

## Others

Keep an eye on [this topic for more inspiration](https://github.com/topics/trading-bot)
