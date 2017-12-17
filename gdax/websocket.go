package gdax

import(
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"fmt"
)

const (
	GDAX_WEBSOCKET_URL = "wss://ws-feed.gdax.com"
)

func (g *GDAX) WebsocketSubscribe(product string, conn *websocket.Conn) error {
	json, err := json.Marshal(WebsocketSubscribe{"subscribe", product})
	if err != nil {
		return err
	}

	if err = conn.WriteMessage(websocket.TextMessage, json); err != nil {
		return err
	}

	return nil
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
func handleEvent(resp []byte) (WebsocketMatch, error) {
	message := WebsocketMatch{}

	type MsgType struct {
		Type string `json:"type"`
	}

	msgType := MsgType{}
	if err := json.Unmarshal(resp, &msgType); err != nil {
		return message, err;
	}

	if msgType.Type == "error" {
		return message, fmt.Errorf("%s", string(resp))
	}

	if msgType.Type == "match" {
		if err := json.Unmarshal(resp, &message); err != nil {
			return message, err
		}
		return message, nil;
	}

	return message, fmt.Errorf("Unsupported message recived")
}

func (g *GDAX) Live(messages chan WebsocketMatch, quit chan bool) {
	var Dialer websocket.Dialer
	conn, _, err := Dialer.Dial(GDAX_WEBSOCKET_URL, http.Header{})
	if err != nil {
		log.Printf("Unable to connect to Websocket. Error: %s\n", err)
		close(messages)
		return
	}

	cleanup := func(){
		conn.Close()
		log.Println("Websocket client disconnected.")
		close(messages)
	}

	log.Println("Connected to Websocket.")

	for _, x := range g.Currencies {
		if err = g.WebsocketSubscribe(x, conn); err != nil {
			log.Printf("Websocket subscription error: %s\n", err)
			cleanup()
			return
		}
	}

	log.Println("Subscribed to product messages.")

	// Poll the quit channel. If you are told to quit, close messages and return.
	for {
		select {
			case <- quit:
				cleanup()
				return
			default:
				_, resp, err := conn.ReadMessage();
				if err != nil {
					log.Println(err)
					cleanup()
					return
				}
		
				data, err := handleEvent(resp); 
				if err != nil {
					continue
				}
		
				messages <- data
		}
	}
}