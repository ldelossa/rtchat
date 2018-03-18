package chatserver

import (
	"log"
	"net/http"

	"github.com/ldelossa/rtchat/pubnub"

	"github.com/gorilla/websocket"
	js "github.com/ldelossa/rtchat/jsonerror"
)

func NewHTTPServer(addr string) *http.Server {
	// Create server with listening port
	s := &http.Server{
		Addr: addr,
	}

	// Create new serve mux
	m := http.NewServeMux()

	// Attach handlers
	m.HandleFunc("/chat/connect", ConnectHandler)

	// Attach mux to server
	s.Handler = m

	return s
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	// Switch statement to lock down requqest to post
	switch r.Method {
	case "GET":
		// Grab name of group chat. This will become a channel we subscribe to
		r.ParseForm()
		channel := r.FormValue("group")
		if channel == "" {
			js.Error(w,
				&js.Response{
					Message: "group query parameter was not provided in url",
				},
				http.StatusBadRequest)
			return
		}

		// Create connection to PubNub
		pn, err := pubnub.NewPubNub()
		if err != nil {
			js.Error(w,
				&js.Response{
					Message: "failed to create connection to PubNub service",
				},
				http.StatusInternalServerError)
			return
		}

		// Allocate go channels to receive messages and errors
		successChannel := make(chan []byte)
		errorChannel := make(chan []byte)

		// Subscribe to group channel
		go pn.Subscribe(channel, "", successChannel, false, errorChannel)

		// Upgrade to websockets
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			log.Printf("upgrade to websocket failed: %s", err)
			js.Error(w,
				&js.Response{
					Message: "upgrade to websocket failed",
				},
				http.StatusInternalServerError)
			return
		}

		// Start reading off conn for published ws messages
		WSReadChan := make(chan ValidMessage)
		go StartReading(conn, WSReadChan)

		// Enter message routing loop
		for {
			select {
			// Received message on subscribed pubnub channel, send over websocket
			case response := <-successChannel:
				// Extract json
				json, err := pubnub.ExtractMessage(response)
				if err != nil {
					log.Printf("received non published message: %s", err)
					continue
				}
				log.Printf("received published message: %v", json)

				// Write message to websocket
				if err = conn.WriteJSON(json); err != nil {
					log.Printf("could not write received pubnub message to websocket: %s", err)
					continue
				}
			// Received messsage on subscribed error pubnub channel, log and continue
			case err := <-errorChannel:
				log.Printf(string(err))
			// Recieved message on websocket read channel, publish to PubNub to update group chat
			case json := <-WSReadChan:

				// Publish received json to PubNub
				err = pubnub.PublishMessage(channel, &json, pn)
				if err != nil {
					log.Printf("publishing websocket messsage to channel %s failed: %s", channel, err)
				}
				log.Printf("published message: %v to channel: %s", json, channel)
			}
		}
	default:
		log.Printf("unsupported HTTP method")
		js.Error(w,
			&js.Response{Message: "unsupported http method"},
			http.StatusBadRequest)
		return
	}
}
