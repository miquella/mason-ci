package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var (
	Router = mux.NewRouter()
)

func init() {
	Router.Path("/slave").Headers("Upgrade", "websocket").Handler(websocket.Handler(slaveWebsocketHandler))
}

func main() {
	http.Handle("/", Router)
	err := http.ListenAndServe(":9535", nil)
	if err != nil {
		log.Panic(err.Error())
	}
}
