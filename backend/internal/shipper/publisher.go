package shipper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

// Publisher wraps a Kafka producer and derives the topic from the log's service field.
type Publisher struct {
	producer *kafka.Producer
}

func NewPublisher(producer *kafka.Producer) *Publisher {
	return &Publisher{producer: producer}
}

// serviceExtractor is a minimal struct to parse just the service field from a JSON log line.
type serviceExtractor struct {
	Service string `json:"service"`
}

// Publish parses the service field from the JSON line and publishes to the correct Kafka topic.
// Topic derivation: "{service}-logs-topic" (matches constants in constants/topics.go).
func (p *Publisher) Publish(line []byte) error {
	var s serviceExtractor
	if err := json.Unmarshal(line, &s); err != nil {
		return fmt.Errorf("failed to parse service field: %w", err)
	}
	if s.Service == "" {
		return fmt.Errorf("log line missing 'service' field")
	}

	topic := s.Service + "-logs-topic"
	p.producer.Publish(topic, line)
	log.Printf("shipped log to topic=%s\n", topic)
	return nil
}
