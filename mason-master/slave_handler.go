package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/miquella/mason-ci/messages"
	"golang.org/x/net/websocket"
)

func slaveWebsocketHandler(ws *websocket.Conn) {
	var data []byte
	for {
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			if err != io.EOF {
				log.Printf("Slave: websocket error (%s)\n", err)
			}
			return
		}

		var frame messages.Frame
		err = json.Unmarshal(data, &frame)
		if err != nil {
			log.Printf("Slave: invalid message (%s)\n", err)
			continue
		}

		if frame.MessageType == "" {
			log.Print("Slave: unknown message received")
			continue
		}

		log.Printf("Slave: Message received: %s", frame.MessageType)
	}
}
