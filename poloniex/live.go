package poloniex

import(
	"log"
	// "net/http"
	// "encoding/json"
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
)

const (
	POLONIEX_WEBSOCKET_URL = "wss://api.poloniex.com"
	POLONIEX_REALM = "realm1"
)

func (p *Poloniex) Live() {
	for p.Enabled {
		// Connect subscriber session.
		sub, err := client.ConnectNet(POLONIEX_WEBSOCKET_URL, client.ClientConfig{
			Realm:  POLONIEX_REALM,
		})
		if err != nil {
			log.Print(err)
			continue
		}
		defer sub.Close()

		log.Print("Connected")

		handleEvent := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
			log.Println("Received event")
			if len(args) != 0 {
				log.Println("  Event Message:", args[0])
			}
			/**
			 * 1. Find type field
			 * 2. If its "newTrade":
			 * 3.		Find message data
			 * 4.		trade := WebsocketMatch{
			 * 				Type 		= data.Type,
			 * 				TradeID = strconv.ParseInt(data.TradeID, 10, 64),
								Rate 		= strconv.ParseFloat(data.Rate, 64),
								Amount	= strconv.ParseFloat(data.Amount, 64),
								Date		= data.Date,
								Total		= strconv.ParseFloat(data.Total, 64),
			 * 			}
			 * 5. Do something with the data.
			 */
		}

		// Subscribe to currencies
		for _,currency := range p.Currencies {
			if err := sub.Subscribe(currency, handleEvent, nil); err != nil {
				log.Printf("Error subscribing to %s channel: %s\n", currency, err)
			}
		}

		log.Print("Subscribed")

		<-sub.Done()
		log.Println("Websocket Client Disconnected")
	}
}