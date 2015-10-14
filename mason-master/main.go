package main

import (
	"flag"
	"github.com/miquella/mason-ci/datastore"
	"github.com/miquella/mason-ci/datastore/drivers/rethink"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var (
	Router = mux.NewRouter()
	Store  *datastore.Datastore
)

func init() {
	rethinkAddress := flag.String("rethink-address", "localhost:28015", "Rethink address")
	rethinkDatabase := flag.String("rethink-database", "mason_ci", "Rethink database name")
	flag.Parse()

	log.Printf("Connecting to datastore %s (%s)", *rethinkAddress, *rethinkDatabase)
	driver, err := rethink.NewRethinkDriver(*rethinkAddress, *rethinkDatabase)
	if err != nil {
		log.Fatalf("Failed to start rethinkdb driver: %s", err)
	}
	Store = datastore.NewDatastore(driver)
}

func main() {
	log.Printf("Listening on :9535")
	http.Handle("/", Router)
	err := http.ListenAndServe(":9535", nil)
	if err != nil {
		log.Panic(err.Error())
	}
}
