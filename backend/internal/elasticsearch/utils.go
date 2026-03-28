package elasticsearch

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	es_api "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type SortCursor struct {
	SortedValue []types.FieldValue `json:"sorted_value"`
}

func EnsureIndex(es *es_api.Client, indexName string, mapping string) error {
	ctx := context.Background()

	res, err := es.Indices.Exists([]string{indexName}, es.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("index %s already exists, skipping creation\n", indexName)
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

func EncodeSortedValue(sortedValue []types.FieldValue) (string, error) {
	if len(sortedValue) == 0 {
		return "", errors.New("sort values cannot be empty")
	}
	sv := SortCursor{SortedValue: sortedValue}
	jsonBytes, err := json.Marshal(sv)
	if err != nil {
		return "", fmt.Errorf("error while marshalling sorted value: %w", err)
	}
	return base64.URLEncoding.EncodeToString(jsonBytes), nil
}

func DecodeSortedValue(cursorString string) ([]types.FieldValue, error) {
	if cursorString == "" {
		return nil, errors.New("cursor string is empty")
	}
	bytes, err := base64.URLEncoding.DecodeString(cursorString)
	if err != nil {
		return nil, fmt.Errorf("cursor decode: invalid base64: %w", err)
	}
	var sortedValue SortCursor
	if err := json.Unmarshal(bytes, &sortedValue); err != nil {
		return nil, fmt.Errorf("error while unmarshalling decoded bytes into sort cursor %w", err)
	}
	return sortedValue.SortedValue, nil
}
