package main

import (
	"log"

	"github.com/ldelossa/rtchat/chatserver"
)

var Addr = "localhost:8080"

func main() {
	// Get our HTTPServer

	s := chatserver.NewHTTPServer(Addr)

	log.Printf("lauching HTTP server on %s", Addr)
	s.ListenAndServe()
}
