package apigateway

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
	livews "github.com/mahirjain10/logflow/backend/internal/websocket"
)

func StartLiveLogConsumer(ctx context.Context, consumer *kafka.Consumer, hub *livews.Hub) {
	topics := []string{constants.LIVE_LOGS_TOPIC}
	if err := consumer.Start(ctx, topics, func(msg *sarama.ConsumerMessage) error {
		log.Println("received live log", string(msg.Value))
		hub.Broadcast(msg.Value)
		return nil
	}); err != nil {
		log.Printf("live-log consumer stopped: %v\n", err)
	}
}
