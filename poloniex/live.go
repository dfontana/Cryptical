package poloniex

/**
 * Websocket functionality for Poloniex using the undocumented
 * weboscket feed. This API has been modified from the original
 * author's code found here https://github.com/pharrisee/poloniex-api/
 */

import(
	"log"
	"encoding/json"
	"time"
	"strconv"
	"math"
	"net/http"
	"github.com/gorilla/websocket"
)

const (
	POLONIEX_WEBSOCKET_URL = "wss://api2.poloniex.com"
	POLONIEX_TICKER = "https://poloniex.com/public?command=returnTicker"
)

// Live streams data for the given currencies into the matches channel.
// The stream stops once a message is sent into the quit channel.
// TODO currently subscribes to BTC USD only, needs ID lookup logic
// implemented based on currency strings passed in.
func Live(currencies []string, matches chan WSOrderbook, quit chan bool) {
	// Connect
	var Dialer websocket.Dialer
	conn, resp, err := Dialer.Dial(POLONIEX_WEBSOCKET_URL, nil)
	if err != nil {
		log.Printf("Unable to connect to Websocket. Error: %s\n", err)
		log.Printf("%+v\n", resp)
		return
	}

	// Function for cleanup if something goes wrong
	cleanup := func(){
		conn.Close()
		log.Println("Websocket client disconnected.")
		close(matches)
	}

	// Lookup currency pair ids & subscribe
	markets, err := getTickers()
	if err != nil {
		log.Fatal(err)
	}
	for _,curr := range currencies {
		if v, ok := markets[curr]; ok {
			if err = subscribe(conn, v.ID); err != nil {
				log.Println("Failed to subscribe to "+curr)
			}else{
				log.Println("Subscribed to: ", curr, v.ID)
			}
		}else{
			log.Fatal("Invalid currency pair provided: " + curr)
		}
	}

	//Listen, quitting when told.
	for {
		select {
		case <- quit:
			cleanup()
			return
		default:
			// Get a message
			_, resp, err := conn.ReadMessage();
			if err != nil {
				log.Println(err)
				cleanup()
				return
			}

			// Determine message type
			message := []interface{}{}
			if err := json.Unmarshal(resp, &message); err != nil {
				log.Printf("%s: (%s)\n", err, string(resp))
				continue
			}
			chid := toFloat(message[0])
			if chid > 100.0 && chid < 1000.0 {
				// It's an orderbook message, which we want
				orderbook, err := parseCurr(message, markets)
				if err != nil {
					log.Println(err)
					continue
				}

				// Dump each item into the channel
				for _,item := range orderbook {
					matches <- item
				}
			}
		}
	}
}

// parseCurr decomposes the raw message into a currency's
// orderbook structure.
func parseCurr(raw []interface{}, markets map[string]TickerEntry) ([]WSOrderbook, error){
	trades := []WSOrderbook{}
	marketID := int64(toFloat(raw[0]))
	market := "UNMAPPED"
	for k,v := range markets {
		if v.ID == marketID {
			market = k
		}
	}

	for _, _v := range raw[2].([]interface{}) {
		v := _v.([]interface{})
		trade := WSOrderbook{}
		trade.Pair = market
		switch v[0].(string) {
		case "i":
		case "o":
			trade.Event = "modify"
			if t := toFloat(v[3]); t == 0.0 {
				trade.Event = "remove"
			}
			trade.Type = "ask"
			if t := toFloat(v[1]); t == 1.0 {
				trade.Type = "bid"
			}
			trade.Rate = toFloat(v[2])
			trade.Amount = toFloat(v[3])
			trade.TS = time.Now()
		case "t":
			trade.Event = "trade"
			trade.TradeID = int64(toFloat(raw[1]))
			trade.Type = "sell"
			if t := toFloat(v[2]); t == 1.0 {
				trade.Type = "buy"
			}
			trade.Rate = toFloat(v[3])
			trade.Amount = toFloat(v[4])
			trade.Total = trade.Rate * trade.Amount
			t := time.Unix(int64(toFloat(v[5])), 0)
			trade.TS = t
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// Subscribes to the given channel
func subscribe(conn *websocket.Conn, channel int64) error {
	message := struct {
		Command string `json:"command"`
		Channel int64 `json:"channel"`
	}{
		"subscribe",
		channel,
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if err = conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
		return err
	}
	return nil
}

func getTickers() (retval map[string]TickerEntry, err error) {
	res, err := http.Get(POLONIEX_TICKER)
	if err != nil {
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&retval)
	return
}

func toFloat(i interface{}) float64 {
	maxFloat := float64(math.MaxFloat64)
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return maxFloat
		}
		return a
	case float64:
		return i
	case int64:
		return float64(i)
	case json.Number:
		a, err := i.Float64()
		if err != nil {
			return maxFloat
		}
		return a
	}
	return maxFloat
}