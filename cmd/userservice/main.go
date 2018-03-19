package main

import (
	"log"
	"os"

	us "github.com/ldelossa/rtchat/userservice"

	pg "github.com/ldelossa/rtchat/userservice/postgres"
)

const (
	Addr = "0.0.0.0:8080"
)

// Default db connection string expecting postgres to be running on your localhost at default port.
// Use POSTGRESCONNSTR env variable to overwrite
var ConnString = "user=postgres dbname=userservice password=dev host=localhost sslmode=disable"

func main() {

	// Overwrite ConnStr if environemnt varible "POSTGRESCONNSTR" exists
	if os.Getenv("POSTGRESCONNSTR") != "" {
		ConnString = os.Getenv("POSTGRESCONNSTR")
	}

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
