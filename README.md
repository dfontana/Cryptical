# Cryptical
[![GoDoc](https://godoc.org/github.com/dfontana/Cryptical?status.svg)](https://godoc.org/github.com/dfontana/Cryptical)

## TODO:
- Test cases.

Simulation:
- Integrate fees into the simulation code (Meaning each exchange client will need to provide their fees, and trades should store the fee they incur)

Automation:
- Convert the simulation code into an actual routine that plugs into live websocket feeds.
- Add authenticated endpoints to GDax Client that will make trades from our trade data.

Data:
- Interface: Make a D3JS webpage for better interaction with Computation plots. (Perhaps can use react to get some practice)
- Make computation plots capable of updating with realtime data.
- Experiment with combining multiple Computation inferences, for a higher order strategy.

## How to run:
See `_examples/` for ways to work with the API. Also be sure to see the GoDoc, link provided in the badge at the top. To run these examples just `cd` into their folder and run `go run main.go`. Be sure you have already installed the dependencies with `go get`.

## Dependencies:
- github.com/gorilla/websocket
- github.com/preichenberger/go-gdax
- github.com/wcharczuk/go-chart
- github.com/joho/godotenv

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot