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
- github.com/dfontana/Cryptical (technically...)

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot


## Temporary Simulation notes:
This function aims to provide the means to test strategies against historical data feeds, in effort to learn how well those strategies could perform. While there are sandbox environments out there, those would be better suited as a "systems" test - can the bot function under time pressure. This simulator aims to not evaluate that, but how well the bot gain *achieve*.

#### The Algorithm
1. Pick a day in time to evaluate from.
2. Create a "live" feed channel.
3. Start detection routine with that date and livefeed (so it can prepare historical computations if needed, like in MACD).
4. Request historical data for that choosen data in as fine of granularity as possible (30 mins)
5. Feed this data into the live feed at a timed rate (ie 1/sec)
6. Simulate:
 - A) The routine simply logs if it would buy or sell at a given point, later this log is feed through a process function that would determine how much it earned/lost
 - B) Facade the "buy and sell" routines for the bot, simulating buying and selling (so the bot requests to make a buy, you do so with necessary fees, & ditto for sell)
7. Log the final evaluation (profits, number of trades, when and what for it traded).