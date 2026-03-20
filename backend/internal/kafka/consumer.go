package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// Custom Handler function should have this signature
// Accepts kafka stream messages and returns error if any
type MessageHandler func(msg *sarama.ConsumerMessage) error

// Wrapper around consumer group
type Consumer struct {
	group sarama.ConsumerGroup
}

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	group, err := sarama.NewConsumerGroup(brokers, groupID, NewConfig())
	if err != nil {
		return nil, err
	}
	return &Consumer{group: group}, nil
}

func (c *Consumer) Start(
	ctx context.Context,
	topics []string,
	handler MessageHandler,
) error {
	h := &consumerGroupHandler{
		handler: handler,
	}

	for {
		// Provided by saram itself
		if err := c.group.Consume(ctx, topics, h); err != nil {
			log.Println("Kafka consume error:", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	// This code passes our handler func to the abstraction to work with our handler code
	handler MessageHandler
}

// These are lifecycle methods.
// Currently returning nil as there are no errors

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim is called internally by sarama itself
func (h *consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {

		if err := h.handler(msg); err != nil {
			log.Println("Message handling failed:", err)
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
