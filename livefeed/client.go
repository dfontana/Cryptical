package main

import(
	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"
)

// Main subscribes to the websocket feed
func main() {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}
	
	// Build our subscription message.
	subscribe := sub_message([]string{"level2"}, []string{"BTC-USD"})

	// Send the subscribe message to the server.
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}
	
	// Start reciving messages
	recv_messages(wsConn)
}

// handle_message subscribes to a channel, and is run as a goroutine.
// When a message arrives, its value is determined and considered against
// the running average. The user is notified if its an anomaly. Then the 
// value is weighted into the average for future consideration.
func handle_message(message gdax.Message){
	
}

// recv_messages is a looping function to recieve messages. Ideally
// this is run in a lone goroutine to handle receving of messages
// into a shared channel.
func recv_messages(wsConn ws.Conn){
	message:= gdax.Message{}
	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}
	
		if message.Type == "match" {
			println("Got a match")
			// TODO this will eventually be refactored out, instead subscribing
			// to a channel
			handle_message(message)
		}
	}
}

// sub_message builds a message for the given channels and products,
// creating out subscription service. 
func sub_message(channels []string, products []string) gdax.Message {
	chan_msgs := make([]gdax.MessageChannel, len(channels))
	for channel := range channels {
		msg := gdax.MessageChannel{
			Name: channel,
			ProductIds: products,
		}
		chan_msgs = append(chan_msgs, msg)
	}

	return gdax.Message{
		Type:      "subscribe",
		Channels: chan_msgs,
	}
}
