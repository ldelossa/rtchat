package chatserver

import (
	"log"

	"github.com/gorilla/websocket"
)

// ValidMessage is the json structure we will accept over the websocket. If the schema does not match
// ValidMessage the message is rejected.
type ValidMessage struct {
	// Type of message. Direct or Channel. If Direct To must be specified.
	Type string `json:"type"`
	// From indicates the user this message is from
	From string `json:"from"`
	// Channel tells us where to route this message too. This is mainly used when type
	// is "direct" and we are sending to a user specific channel
	Channel string `json:channel`
	// Text is the text payload of our message
	Text string `json:"text"`
}

// StartReader is meant to be ran as a go routine and allows for asynchronous processing of websocket writes
// in our connect hander loop.
func StartReading(conn *websocket.Conn, msgChan chan ValidMessage) {
	for {
		var msg ValidMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("could not deserialize received websocket message: %s", err)
			continue
		}

		// Validate message
		if msg.Type == "" {
			log.Printf("message received without type field")
			continue
		}
		if msg.From == "" {
			log.Printf("message received without from field")
			continue
		}
		if msg.Channel == "" {
			log.Printf("message received without channel field")
			continue
		}
		if msg.Text == "" {
			log.Printf("message received without text field")
			continue
		}

		msgChan <- msg
	}
}
