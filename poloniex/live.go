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
	"github.com/gorilla/websocket"
)

const (
	POLONIEX_WEBSOCKET_URL = "wss://api2.poloniex.com"
	POLONIEX_BTC = "121" 		 //	Get from "returnTicker" endpoint
)

func Live() {
	// Connect
	var Dialer websocket.Dialer
	conn, resp, err := Dialer.Dial(POLONIEX_WEBSOCKET_URL, nil)
	if err != nil {
		log.Printf("Unable to connect to Websocket. Error: %s\n", err)
		log.Printf("%+v\n", resp)
		return
	}

	// Subscribe to BTC
	subscribe(conn, POLONIEX_BTC)

	//Listen
	for {
		_, resp, err := conn.ReadMessage();
		if err != nil {
			log.Println(err)
			return
		}

		handleEvent(resp); 
	}
}

func handleEvent(resp []byte) {
	message := []interface{}{}
	if err := json.Unmarshal(resp, &message); err != nil {
		log.Fatal(err)
	}
	chid := toFloat(message[0])
	if chid > 100.0 && chid < 1000.0 {
		orderbook, err := parseCurr(message)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%+v\n",orderbook)
	}
}

func parseCurr(raw []interface{}) ([]WSOrderbook, error){
	trades := []WSOrderbook{}
	//marketID := int64(toFloat(raw[0]))

	for _, _v := range raw[2].([]interface{}) {
		v := _v.([]interface{})
		trade := WSOrderbook{}
		trade.Pair = "UNMAPPED"
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
func subscribe(conn *websocket.Conn, channel string) error {
	message := struct {
		Command string `json:"command"`
		Channel string `json:"channel"`
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