package main

import "log"

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients - using a map for O(1) lookups
	clients map[*Client]bool

	// Inbound messages from clients to broadcast
	broadcast chan []byte

	// Register requests from new clients
	register chan *Client

	// Unregister requests from disconnecting clients
	unregister chan *Client
}

// newHub creates a new Hub instance
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// run starts the hub's main loop - call this in a goroutine
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// Add new client to the map
			h.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			// Remove client if it exists
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			// Send message to all connected clients
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message queued successfully
				default:
					// Client's send buffer is full - assume it's dead
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
