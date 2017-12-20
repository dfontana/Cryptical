package poloniex

/**
 * Websocket functionality for Poloniex using the undocumented
 * weboscket feed. This API has been modified from the original
 * author's code found here https://github.com/pharrisee/poloniex-api/
 */

import(
	"log"
	"net/http"
	"encoding/json"
	"time"
	"errors"
	"strconv"
	"math"
	"github.com/gorilla/websocket"
)

const (
	POLONIEX_WEBSOCKET_URL = "wss://api2.poloniex.com"
	POLONIEX_TICKER = "1002" //	Hard code :S
	POLONIEX_BTC = "121" 		 //	Get from "returnTicker" endpoint
)

func (p *Poloniex) Live() {
	// Connect
	var Dialer websocket.Dialer
	conn, resp, err := Dialer.Dial(POLONIEX_WEBSOCKET_URL, http.Header{})
	if err != nil {
		log.Printf("Unable to connect to Websocket. Error: %s\n", err)
		log.Printf("%+v\n", resp)
		return
	}

	// Subscribe to ticker and BTC
	subscribe(conn, POLONIEX_TICKER)
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
	channelID := toFloat(message[0])
	ticker,_ := strconv.ParseFloat(POLONIEX_TICKER, 64)
	if channelID < 1000 && channelID > 100 {
		// it's an orderbook
		orderbook, err := parseCurr(message)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%+v\n",orderbook)
	} else if channelID ==  ticker {
		// it's a ticker
		ticker, err := parseTicker(message)
		if err != nil {
			log.Printf("%s: (%s)\n", err, message)
			return
		}
		log.Printf("%+v\n",ticker)
	}
}

func parseTicker(raw []interface{}) (WSTicker, error) {
	wt := WSTicker{}
	var rawInner []interface{}
	if len(raw) <= 2 {
		return wt, errors.New("Not a ticker item.")
	}
	rawInner = raw[2].([]interface{}) 
	marketID := int64(toFloat(rawInner[0]))
	wt.Pair = "UNMAPPED"
	wt.PairID = marketID
	wt.Last = toFloat(rawInner[1])
	wt.Ask = toFloat(rawInner[2])
	wt.Bid = toFloat(rawInner[3])
	wt.PercentChange = toFloat(rawInner[4])
	wt.BaseVolume = toFloat(rawInner[5])
	wt.QuoteVolume = toFloat(rawInner[6])
	wt.IsFrozen = toFloat(rawInner[7]) != 0.0
	wt.DailyHigh = toFloat(rawInner[8])
	wt.DailyLow = toFloat(rawInner[9])
	return wt, nil
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
		POLONIEX_TICKER,
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