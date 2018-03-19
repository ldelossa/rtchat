package main

import (
	"log"

	us "github.com/ldelossa/rtchat/userservice"

	pg "github.com/ldelossa/rtchat/userservice/postgres"
)

var Addr = "localhost:8081"
var ConnString = "user=postgres dbname=userservice password=dev host=localhost sslmode=disable"

func main() {

	// Create DataStore
	ds, err := pg.NewDatastore(ConnString)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// Create HTTPServer
	s, err := us.NewHTTPServer(Addr, ds)
	if err != nil {
		log.Fatalf("could not create HTTP server")
	}

	// Listen and Server
	log.Printf("lauching http server on %s", Addr)
	s.ListenAndServe()
}
