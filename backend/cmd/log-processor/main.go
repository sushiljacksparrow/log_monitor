package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/mahirjain10/logflow/backend/internal/common/types"
	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/elasticsearch"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

// func initAuthServiceLogIndex(es *es.Client) (int, error) {
// 	mapping := `{
//   "mappings": {
//     "properties": {
//       "service":   { "type": "keyword" },
//       "level":     { "type": "keyword" },
//       "message":   { "type": "text" },
//       "requestId": { "type": "keyword" },
//       "userId":    { "type": "keyword" },
//       "ip":        { "type": "keyword" },
//       "timestamp": { "type": "date" }
//     }
// 		}
// 	}`

//		res, err := es.Indices.Create(
//			AUTH_SERVICE_LOGS,
//			es.Indices.Create.WithBody(strings.NewReader(mapping)),
//		)
//		if err != nil {
//			return 0, fmt.Errorf("failed to create auth-service-logs index: %w", err)
//		}
//		log.Println(res)
//		log.Println("auth-service-logs index created:", res.StatusCode)
//		return res.StatusCode, nil
//	}

func main() {
	var bulkIndexers *types.BulkIndexers
	bulkIndexers = &types.BulkIndexers{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error while initalizing config: %v", err)
	}
	esClient, err := elasticsearch.InitES(config)
	if err != nil {
		log.Fatal(err)
	}
	var bulkIndexer esutil.BulkIndexer
	for k, v := range IndexMappings {
		err = elasticsearch.EnsureIndex(esClient, k, v)
		if err != nil {
			log.Println(err)
		}
		bulkIndexer, err = InitBulkIndexer(esClient, k)
		if err != nil {
			log.Printf("couldnt build bulk indexer for %s: %v", k, err)
		}
		EsBgCtx := context.Background()
		defer bulkIndexer.Close(EsBgCtx)
		switch k {
		case AUTH_SERVICE_LOGS:
			bulkIndexers.AuthBulkIndexer = bulkIndexer
		case ORDER_SERVICE_LOGS:
			bulkIndexers.OrderBulkIndexer = bulkIndexer
		case PAYMENT_SERVICE_LOGS:
			bulkIndexers.PaymentBulkIndexer = bulkIndexer
		default:
		}

	}
	producer, err := kafka.NewProducer(config.KafkaBrokers)
	if err != nil {
		log.Printf("error while init new producer:%v ", err)
	}
	defer producer.Close()
	consumer, err := kafka.NewConsumer(config.KafkaBrokers, config.KafkaTopicLogGroupID)
	if err != nil {
		log.Printf("error while initalizing a consumer with group ID: %s - %v ", config.KafkaTopicLogGroupID, err)
	}
	fmt.Println("Consumption starting")
	topics := []string{config.KafkaTopicAuthLog, config.KafkaTopicOrderLog, config.KafkaTopicPaymentLog}
	if err := consumer.Start(ctx, topics, func(msg *sarama.ConsumerMessage) error {
		return ConsumeKafka(msg, bulkIndexers, producer)
	}); err != nil {
		log.Println("failed to consumed message: ", err)
	}
}
