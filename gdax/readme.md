# Package: gdax

Create a new client by declaring a gdax struct: `g := GDAX{true, []string{"ETH-USD"}}`.

Properties:
- Enabled: boolean, declaring if it is currently enabled or not (to stop websocket feed)
- Currencies: string list of currencies you'd like to track in the websocket feed.

Methods:
- `g.Live()` starts the websocket feed in the current thread. You'll want to start this in a new thread to prevent blocking. Stop by setting the enabled property to false.
- `g.Historic()` gets historic rates for a given time period and granularity with `g.Historic`