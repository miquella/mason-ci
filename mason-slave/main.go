package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/miquella/mason-ci/messages"
	"golang.org/x/net/websocket"
)

var (
	MasterURL string

	OriginURL    *url.URL
	SlavePath, _ = url.Parse("slave")
)

func init() {
	flag.StringVar(&MasterURL, "master-url", "", "URI for master node")
	flag.Parse()

	var err error
	OriginURL, err = url.Parse(MasterURL)
	if err != nil {
		log.Fatal(err)
	}
	switch OriginURL.Scheme {
	case "http":
		OriginURL.Scheme = "ws"
	case "https":
		OriginURL.Scheme = "wss"
	default:
		log.Fatalf("Invalid master url scheme: %s\n", OriginURL.Scheme)
	}
}

func main() {
	slaveUrl := OriginURL.ResolveReference(SlavePath)
	ws, err := websocket.Dial(slaveUrl.String(), "", OriginURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// register me yay
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	msg := &messages.Register{
		Frame:    messages.Frame{"Register"},
		Hostname: hostname,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("BAD: died - %s\n", err)
	}

	err = websocket.Message.Send(ws, b)
	if err != nil {
		log.Fatalf("BAD: died - %s\n", err)
	}
}
