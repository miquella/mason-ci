package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var (
	Router = mux.NewRouter()
	Store  Datastore
)

func init() {
	Router.Path("/slave").Headers("Upgrade", "websocket").Handler(websocket.Handler(slaveWebsocketHandler))

	var rethinkAddress string
	var rethinkDatabase string
	flag.StringVar(&rethinkAddress, "rethink-address", "localhost:28015", "Rethink address")
	flag.StringVar(&rethinkDatabase, "rethink-database", "mason_ci", "Rethink database name")
	flag.Parse()

	var err error
	Store, err = NewRethinkDatastore(rethinkAddress, rethinkDatabase)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.Handle("/", Router)
	err := http.ListenAndServe(":9535", nil)
	if err != nil {
		log.Panic(err.Error())
	}
}
