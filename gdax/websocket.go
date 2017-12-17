package gdax

import(
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
)

const (
	GDAX_WEBSOCKET_URL = "wss://ws-feed.gdax.com"
)

func (g *GDAX) WebsocketSubscribe(product string, conn *websocket.Conn) error {
	subscribe := WebsocketSubscribe{"subscribe", product}
	json, err := json.Marshal(subscribe)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		return err
	}
	return nil
}

func (g *GDAX) Live() {
	for g.Enabled {
		var Dialer websocket.Dialer
		conn, _, err := Dialer.Dial(GDAX_WEBSOCKET_URL, http.Header{})

		if err != nil {
			log.Printf("Unable to connect to Websocket. Error: %s\n", err)
			continue
		}

		log.Printf("Connected to Websocket.\n")

		for _, x := range g.Currencies {
			err = g.WebsocketSubscribe(x, conn)
			if err != nil {
				log.Printf("Websocket subscription error: %s\n", err)
				continue
			}
		}

		log.Printf("Subscribed to product messages.")

		for g.Enabled {
			_, resp, err := conn.ReadMessage();
			if err != nil {
				log.Println(err)
				break
			}

			type MsgType struct {
				Type string `json:"type"`
			}

			msgType := MsgType{}
			if err := json.Unmarshal(resp, &msgType); err != nil {
				log.Println(err)
				continue
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
			if msgType.Type == "error" {
				log.Println(string(resp))
				break
			}

			if msgType.Type == "match" {
				message := WebsocketMatch{}
				if err := json.Unmarshal(resp, &message); err != nil {
					log.Println(err)
					continue
				}
				log.Printf("%+v\n", message)
			}

		}
		conn.Close()
		log.Printf("Websocket client disconnected.")
	}
}