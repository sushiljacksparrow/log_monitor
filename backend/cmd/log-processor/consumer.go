package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/mahirjain10/logflow/backend/internal/common/types"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

type BaseLog struct {
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	RequestID string    `json:"requestId"`
	Timestamp time.Time `json:"timestamp"`
}

type AuthLog struct {
	BaseLog
	UserID string `json:"userId"`
	IP     string `json:"ip"`
}

type OrderLog struct {
	BaseLog
	OrderID string `json:"orderId"`
	Carrier string `json:"carrier"`
}

type PaymentLog struct {
	BaseLog
	PaymentID string `json:"paymentId"`
	Gateway   string `json:"gateway"`
}

type IndexError struct {
	IndexName  string
	RetryCount int
	Error      error
	Body       string
}
type RetryMap struct {
	RetryMap map[string]IndexError
	mu       sync.Mutex
}

func ConsumeKafka(msg *sarama.ConsumerMessage, bulkIndexer *types.BulkIndexers, producer *kafka.Producer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	retryMap := RetryMap{
		RetryMap: make(map[string]IndexError),
	}
	log.Println("I got the messaged baby: ", string(msg.Value))
	log.Println(msg.Topic)
	stringifiedMessage := string(msg.Value)
	switch msg.Topic {
	case "logs-auth-service":
		if err := AddIndex(ctx, bulkIndexer.AuthBulkIndexer, AUTH_SERVICE_LOGS, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	case "logs-order-service":
		if err := AddIndex(ctx, bulkIndexer.AuthBulkIndexer, ORDER_SERVICE_LOGS, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	case "logs-payment-service":
		if err := AddIndex(ctx, bulkIndexer.AuthBulkIndexer, PAYMENT_SERVICE_LOGS, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	default:
		log.Println("invalid topic")
	}
	return nil
}
