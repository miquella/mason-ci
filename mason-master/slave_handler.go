package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/miquella/mason-ci/messages"
	"github.com/miquella/mason-ci/pool"
	"golang.org/x/net/websocket"
)

var (
	SlavePool = pool.New()

	messageHandlers = make(map[string]func(ws *websocket.Conn, messageType string, data []byte))
)

func init() {
	Router.Path("/slave").Headers("Upgrade", "websocket").Handler(websocket.Handler(slaveWebsocketHandler))
}

func slaveWebsocketHandler(ws *websocket.Conn) {
	var id string
	defer func() {
		if id != "" {
			SlavePool.Unregister(id)
		}
	}()

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

		// Handle "Register" specially
		if frame.MessageType == "Register" {
			var registerMsg messages.Register
			err := json.Unmarshal(data, &registerMsg)
			if err != nil {
				log.Printf("Slave: failed to unmarshal register message: %s\n", err)
				return
			}

			if registerMsg.Hostname == "" {
				ws.Close()
				log.Print("Slave: invalid hostname, closing socket")
				return
			}

			id = registerMsg.Hostname
			SlavePool.Register(id)
		} else if callback, exists := messageHandlers[frame.MessageType]; exists {
			callback(ws, frame.MessageType, data)
		} else {
			log.Printf("Slave: unhandled message: %s", frame.MessageType)
		}
	}
}
