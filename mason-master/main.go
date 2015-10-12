package main

import (
	"flag"
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
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
	flag.StringVar(&rethinkDatabase, "rethink-database", "mason-ci", "Rethink database name")
	flag.Parse()
	rethinkSession, err := r.Connect(r.ConnectOpts{
		Address:  rethinkAddress,
		Database: rethinkDatabase,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	Store = &RethinkDatastore{
		rethinkSession: rethinkSession,
	}
}

func main() {
	http.Handle("/", Router)
	err := http.ListenAndServe(":9535", nil)
	if err != nil {
		log.Panic(err.Error())
	}
}
