# Simulator:
This package aims to provide the means to test strategies against historical data feeds, in effort to learn how well those strategies could perform. While there are sandbox environments out there, those would be better suited as a "systems" test - can the bot function under time pressure. This simulator aims to not evaluate that, but how well the bot gain *achieve*.

## The Algorithm
1. Pick a day in time to evaluate from.
2. Create a "live" feed channel.
3. Start detection routine with that date and livefeed (so it can prepare historical computations if needed, like in MACD).
4. Request historical data for that choosen data in as fine of granularity as possible (30 mins)
5. Feed this data into the live feed at a timed rate (ie 1/sec)
6. Simulate:
 - A) The routine simply logs if it would buy or sell at a given point, later this log is feed through a process function that would determine how much it earned/lost
 - B) Facade the "buy and sell" routines for the bot, simulating buying and selling (so the bot requests to make a buy, you do so with necessary fees, & ditto for sell)
7. Log the final evaluation (profits, number of trades, when and what for it traded).