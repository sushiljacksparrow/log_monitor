package wrapper

import (
	"context"
	"encoding/json"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v9"
	// low level API
	esapi "github.com/elastic/go-elasticsearch/v9/esapi"
)

type Response struct {
	Index     string `json:"index"`
	Health    string `json:"health"`
	DocsCount string `json:"docs.count"`
	DocSize   string `json:"docs.size"`
}

func GetIndexes(ctx context.Context, esClient *es.Client) (*Response, error) {
	// Using low level API
	res, err := esapi.CatIndicesRequest{
		Format: "json",                                                 // JSON output
		H:      []string{"index", "health", "docs.count", "docs.size"}, // Only index names
	}.Do(context.Background(), esClient)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve indexes")
	}
	defer res.Body.Close()
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("couldnt decode body of ES response %v", err)
	}
	return &response, nil
}

func GetMappingViaIndexName(indexNames []string) {

}
