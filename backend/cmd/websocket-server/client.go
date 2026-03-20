package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader with production-ready settings
var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create a new client with a buffered send channel
	// Buffer size of 256 prevents slow clients from blocking broadcasts
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	// Register the client with our hub (we'll build this next)
	hub.register <- client

	// Start the read and write pumps in separate goroutines
	// These handle all communication with this specific client
	go client.writePump()
	go client.readPump()
}
