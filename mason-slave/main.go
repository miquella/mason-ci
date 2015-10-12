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
	URL       *url.URL
)

func init() {
	flag.StringVar(&MasterURL, "master-url", "", "URI for master node")
	flag.Parse()

	u, err := url.Parse(MasterURL)
	if err != nil {
		log.Fatal(err)
	}

	URL = &url.URL{
		Host:   u.Host,
		Scheme: "ws",
		Path:   "/slave",
	}

}

func main() {
	origin := &url.URL{
		Host:   URL.Host,
		Scheme: "http",
		Path:   "/",
	}
	println(URL.String())
	println(origin.String())
	ws, err := websocket.Dial(URL.String(), "", origin.String())
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
