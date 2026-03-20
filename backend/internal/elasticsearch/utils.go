package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"strings"

	es_api "github.com/elastic/go-elasticsearch/v9"
)

func EnsureIndex(es *es_api.Client, indexName string, mapping string) error {
	ctx := context.Background()

	res, err := es.Indices.Exists([]string{indexName}, es.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("index %s already exists, skipping creation", indexName)
		return nil
	}

	createRes, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithContext(ctx),
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("error response while creating index %s: %s", indexName, createRes.String())
	}

	log.Printf("index %s created successfully", indexName)
	return nil
}
