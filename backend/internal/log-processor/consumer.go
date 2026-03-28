package logprocessor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

type BaseLog struct {
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

type AuthLog struct {
	BaseLog
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
}

type OrderLog struct {
	BaseLog
	OrderID   string `json:"order_id"`
	Carrier   string `json:"carrier"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	StockLeft int    `json:"stock_left"`
}

type PaymentLog struct {
	BaseLog
	PaymentID string  `json:"payment_id"`
	Gateway   string  `json:"gateway"`
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
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

func ConsumeKafka(msg *sarama.ConsumerMessage, bulkIndexers *BulkIndexers, producer *kafka.Producer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	retryMap := RetryMap{
		RetryMap: make(map[string]IndexError),
	}

	stringifiedMessage := string(msg.Value)
	log.Println("received message:", stringifiedMessage)
	log.Println("topic:", msg.Topic)

	switch msg.Topic {
	case constants.AUTH_SERVICE_LOGS_TOPIC:
		if err := AddIndex(ctx, bulkIndexers.AuthBulkIndexer, constants.AUTH_SERVICE_LOGS_INDEX, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	case constants.ORDER_SERVICE_LOGS_TOPIC:
		if err := AddIndex(ctx, bulkIndexers.OrderBulkIndexer, constants.ORDER_SERVICE_LOGS_INDEX, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	case constants.PAYMENT_SERVICE_LOGS_TOPIC:
		if err := AddIndex(ctx, bulkIndexers.PaymentBulkIndexer, constants.PAYMENT_SERVICE_LOGS_INDEX, stringifiedMessage, &retryMap, producer); err != nil {
			return err
		}
	default:
		log.Println("invalid topic")
	}

	return nil
}
