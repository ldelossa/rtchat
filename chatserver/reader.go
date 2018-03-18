package chatserver

import (
	"log"

	"github.com/gorilla/websocket"
)

// StartReader is meant to be ran as a go routine and allows for asynchronous processing
// in our connect hander loop.
func StartReading(conn *websocket.Conn, msgChan chan map[string]interface{}) {
	for {
		var json map[string]interface{}
		err := conn.ReadJSON(&json)
		if err != nil {
			log.Printf("could not deserialize received websocket message: %s", err)
			continue
		}
		msgChan <- json
	}
}
