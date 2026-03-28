package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	"github.com/mahirjain10/logflow/backend/internal/elasticsearch"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
	logprocessor "github.com/mahirjain10/logflow/backend/internal/log-processor"
)

func main() {
	indexMappings := map[string]string{
		constants.AUTH_SERVICE_LOGS_INDEX:    logprocessor.AuthServiceLogMapping,
		constants.ORDER_SERVICE_LOGS_INDEX:   logprocessor.OrderServiceLogMapping,
		constants.PAYMENT_SERVICE_LOGS_INDEX: logprocessor.PaymentServiceLogMapping,
	}
	bulkIndexers := &logprocessor.BulkIndexers{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error while initalizing config: %v", err)
	}
	esClient, _, err := elasticsearch.InitES(config)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range indexMappings {
		err = elasticsearch.EnsureIndex(esClient, k, v)
		if err != nil {
			log.Println(err)
		}
		bulkIndexer, err := logprocessor.InitBulkIndexer(esClient, k)
		if err != nil {
			log.Printf("couldnt build bulk indexer for %s: %v\n", k, err)
		}
		esBgCtx := context.Background()
		defer bulkIndexer.Close(esBgCtx)
		switch k {
		case constants.AUTH_SERVICE_LOGS_INDEX:
			bulkIndexers.AuthBulkIndexer = bulkIndexer
		case constants.ORDER_SERVICE_LOGS_INDEX:
			bulkIndexers.OrderBulkIndexer = bulkIndexer
		case constants.PAYMENT_SERVICE_LOGS_INDEX:
			bulkIndexers.PaymentBulkIndexer = bulkIndexer
		default:
		}

	}
	producer, err := kafka.NewProducer(config.KafkaBrokers)
	if err != nil {
		log.Printf("error while init new producer:%v\n", err)
	}
	defer producer.Close()
	consumer, err := kafka.NewConsumer(config.KafkaBrokers, config.KafkaTopicLogGroupID)
	if err != nil {
		log.Printf("error while initalizing a consumer with group ID: %s - %v\n", config.KafkaTopicLogGroupID, err)
	}
	fmt.Println("Consumption starting")

	topics := []string{constants.AUTH_SERVICE_LOGS_TOPIC, constants.ORDER_SERVICE_LOGS_TOPIC, constants.PAYMENT_SERVICE_LOGS_TOPIC}
	if err := consumer.Start(ctx, topics, func(msg *sarama.ConsumerMessage) error {
		return logprocessor.ConsumeKafka(msg, bulkIndexers, producer)
	}); err != nil {
		log.Println("failed to consumed message: ", err)
	}
}
