package logprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	"github.com/mahirjain10/logflow/backend/internal/ids"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
)

type DLQMessage struct {
	DocumentID string `json:"documentId"`
	IndexName  string `json:"indexName"`
	Status     int    `json:"status"`
	Reason     string `json:"reason"`
	Body       string `json:"body"`
}

func AddIndex(ctx context.Context, esIndexer esutil.BulkIndexer, indexName, body string, retryMap *RetryMap, producer *kafka.Producer) error {
	initialBackoff := time.Second
	uuid, err := ids.GenerateUUID()
	if err != nil {
		return fmt.Errorf("error while generating UUID for index: %s", indexName)
	}

	err = esIndexer.Add(ctx, esutil.BulkIndexerItem{
		Index:      indexName,
		Action:     "index",
		DocumentID: uuid,
		Body:       strings.NewReader(body),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			log.Printf("Indexed: %s for index: %s", item.DocumentID, indexName)
			producer.Publish(constants.LIVE_LOGS_TOPIC, []byte(body))
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				log.Printf("Transport error for %s: %v", item.DocumentID, err)
				dlq := DLQMessage{
					DocumentID: item.DocumentID,
					IndexName:  indexName,
					Status:     res.Status,
					Reason:     res.Error.Reason,
					Body:       body,
				}
				data, _ := json.Marshal(dlq)
				producer.Publish(constants.LOGS_DLQ_TOPIC, data)
			}

			switch res.Status {
			case 429, 503, 507:
				retryMap.mu.Lock()
				v, exists := retryMap.RetryMap[item.DocumentID]
				if !exists {
					v = IndexError{IndexName: indexName, RetryCount: 0, Body: body}
				}
				v.RetryCount++
				retryMap.RetryMap[item.DocumentID] = v
				count := v.RetryCount
				retryMap.mu.Unlock()

				if count <= 3 {
					backoff := initialBackoff * time.Duration(1<<count)
					jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
					go func() {
						select {
						case <-time.After(backoff + jitter):
							AddIndex(ctx, esIndexer, indexName, body, retryMap, producer)
						case <-ctx.Done():
						}
					}()
				} else {
					log.Printf("DLQ: %s exhausted retries", item.DocumentID)
				}
			case 400, 404:
				dlq := DLQMessage{
					DocumentID: item.DocumentID,
					IndexName:  indexName,
					Status:     res.Status,
					Reason:     res.Error.Reason,
					Body:       body,
				}
				data, _ := json.Marshal(dlq)
				producer.Publish(constants.LOGS_DLQ_TOPIC, data)
				log.Printf("DLQ: %s status=%d reason=%s", item.DocumentID, res.Status, res.Error.Reason)
			}
		},
	})
	if err != nil {
		return fmt.Errorf("error adding document to bulk indexer: %v", err)
	}

	return nil
}

func InitBulkIndexer(esClient *es.Client, indexName string) (esutil.BulkIndexer, error) {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     esClient,
		Index:      indexName,
		NumWorkers: 4,
		FlushBytes: 5_000_000,
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing bulk indexer: %v", err)
	}

	return indexer, nil
}
