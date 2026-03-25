package elasticsearch

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/mahirjain10/logflow/backend/internal/config"
)

func InitES(config config.Envs) (*elasticsearch.Client, *elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: config.ElasticSearchHost,
	}
	es, err := elasticsearch.NewClient(cfg)
	esTypedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	return es, esTypedClient, nil
}
