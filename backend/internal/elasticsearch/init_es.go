package elasticsearch

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/mahirjain10/logflow/backend/internal/config"
)

func InitES(config config.Envs) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: config.ElasticSearchHost,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	return es, nil
}
