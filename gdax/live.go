package gdax

import(
  ws "github.com/gorilla/websocket"
  gdax "github.com/preichenberger/go-gdax"
  "log"
  "encoding/json"
)

const (
  GDAX_WEBSOCKET_URL = "wss://ws-feed.gdax.com"
)

// Live streams data for the given currencies into the matches channel. 
// To stop the stream, a message should be sent into the quit channel.
func Live(currencies []string, matches chan WsMatch, quit chan bool) {
  // Create a connection
  var Dialer ws.Dialer
  conn, _, err := Dialer.Dial(GDAX_WEBSOCKET_URL, nil)
  if err != nil {
    log.Printf("Unable to connect to Websocket. Error: %s\n", err)
    close(matches)
    return
  }
  log.Println("Connected to Websocket.")

  // Function for cleanup if something goes wrong
  cleanup := func(){
    conn.Close()
    log.Println("Websocket client disconnected.")
    close(matches)
  }

  // Subscribe to all our desired currencies
  subscribe := gdax.Message{
    Type:      "subscribe",
    Channels: []gdax.MessageChannel{
      gdax.MessageChannel{
        Name: "full",
        ProductIds: currencies,
      },
    },
  }
  if err := conn.WriteJSON(subscribe); err != nil {
    log.Printf("Subscription error: %s\n", err)
    cleanup()
    return
  }
  log.Println("Subscribed to product messages.")

  // Poll the quit channel. If you are told to quit, close messages and return.
  for {
    select {
      case <- quit:
        cleanup()
        return
      default:
        // Determine message type
        message := struct {
          Type string `json:"type"`
        }{}
        _, resp, err := conn.ReadMessage();
        if err != nil {
          log.Println(err)
          cleanup()
          return
        }
        if err := json.Unmarshal(resp, &message); err != nil {
          log.Println(err)
          cleanup()
          return
        }

    
        /**
         * Types:
         *  Error: Something went wrong with your request, fix it.
         * 	Recieved: Indicates a valid order is now active in exchange. Ignored.
         * 	Open:	Order is open and not immediately filled. The remaining_size
         * 				indicates how much is left.
         *  Done: Order is out of the book. Can be sent for a cancelled or
         * 				filled order. Remaining_size will be 0 for filled orders.
         * 	Match: Trade occured between two orders. Taker takes from an existing
         * 				 order (maker). Side indicates what side the maker was on (buy
         * 				 or sell). Sell = uptick. Buy = downtick.
         * 	Change: Order has changed in size or funds. Ignored.
         */
        if message.Type == "match" {
          match := WsMatch{}
          if err = json.Unmarshal(resp, &match); err != nil {
            continue;
          }
          matches <- match
        }
    }
  }
}