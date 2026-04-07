package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahirjain10/logflow/backend/internal/kafka"
	"github.com/mahirjain10/logflow/backend/internal/shipper"
)

func main() {
	cfg, err := shipper.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Load saved positions from disk.
	pos := shipper.NewPositionTracker(cfg.PositionFile)
	if err := pos.Load(); err != nil {
		log.Printf("warning: could not load position file: %v\n", err)
	}

	// Create Kafka producer (reuses existing internal/kafka package).
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	pub := shipper.NewPublisher(producer)
	watcher := shipper.NewWatcher(cfg, pos, pub)

	// Periodically save positions to disk.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := pos.Save(); err != nil {
					log.Printf("failed to save positions: %v\n", err)
				}
			}
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("shutting down shipper...")
		cancel()
	}()

	log.Printf("logflow-shipper started: watching=%s, brokers=%v\n", cfg.WatchDir, cfg.KafkaBrokers)

	if err := watcher.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("watcher error: %v", err)
	}

	// Final position save before exit.
	if err := pos.Save(); err != nil {
		log.Printf("failed to save positions on shutdown: %v\n", err)
	}
	log.Println("shipper stopped")
}
