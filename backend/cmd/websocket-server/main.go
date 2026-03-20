package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

var hub *Hub

func main() {
	ctx := context.Background()
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error while initalizing config: %v", err)
	}
	// Create and start the hub
	hub = newHub()
	go hub.run()

	consumer, err := kafka.NewConsumer(config.KafkaBrokers, "logs-live-1")
	if err != nil {
		log.Println("error while intializing consumer")
	}
	topics := []string{"logs-live"}
	go consumer.Start(ctx, topics, func(msg *sarama.ConsumerMessage) error {
		msgVal := msg.Value
		log.Println("recieved message", string(msgVal))
		hub.broadcast <- msgVal
		log.Println("message sent successfully to the clients")
		return nil
	})
	// Set up HTTP routes
	http.HandleFunc("/ws", handleWebSocket)

	// Create the HTTP server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in a goroutine so we can handle shutdown
	go func() {
		log.Printf("WebSocket server starting on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close all WebSocket connections gracefully
	for client := range hub.clients {
		// Send close message to each client
		client.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"))
		client.conn.Close()
	}

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
