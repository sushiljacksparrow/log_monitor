package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the client
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client
	pongWait = 60 * time.Second

	// Send pings to client with this period - must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from client - prevents memory exhaustion attacks
	maxMessageSize = 512 * 1024 // 512KB
)

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	// Ensure cleanup happens when this goroutine exits
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	// Configure connection limits
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Reset the read deadline every time we receive a pong
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Main read loop
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// Check if this is a normal close
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		// Process the message - in this example, we broadcast to all clients
		hub.broadcast <- message
	}
}

// writePump sends messages to the WebSocket connection
func (c *Client) writePump() {
	// Create a ticker for sending periodic pings
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// Set write deadline for every write operation
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The hub closed the channel - send close message
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Get a writer for the next message
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Batch any queued messages into the same WebSocket frame
			// This improves performance when messages arrive faster than we send
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send periodic ping to keep connection alive
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
