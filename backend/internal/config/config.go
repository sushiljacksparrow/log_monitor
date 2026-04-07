package config

import (
	"fmt"
	"os"
	"path/filepath"

	dotenv "github.com/joho/godotenv"
)

type Envs struct {
	KafkaBrokers         []string
	KafkaTopicAuthLog    string
	KafkaTopicOrderLog   string
	KafkaTopicPaymentLog string
	KafkaTopicLogGroupID string
	ElasticSearchHost    []string
}

func InitConfig() (Envs, error) {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		return Envs{}, err
	}
	fmt.Printf("Current working directory: %s\n", dir)

	envPath := filepath.Join(dir, "..", "..", ".env")

	fmt.Printf("Loading .env from: %s\n", envPath)

	if err := dotenv.Load(envPath); err != nil {
		fmt.Printf("Warning: could not load env file %s: %v\n", envPath, err)
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopicLogGroupID := os.Getenv("KAFKA_LOG_GROUP_ID")
	elasticSearchHost := os.Getenv("ELASTIC_SEARCH_HOST")
	if len(kafkaBrokers) == 0 || len(kafkaTopicLogGroupID) == 0 || len(elasticSearchHost) == 0 {
		return Envs{}, fmt.Errorf("KAFKA_BROKERS or KAFKA_LOG_GROUP_ID env not found")
	}
	fmt.Println(kafkaBrokers)
	return Envs{
		KafkaBrokers:         []string{kafkaBrokers},
		KafkaTopicLogGroupID: kafkaTopicLogGroupID,
		ElasticSearchHost:    []string{elasticSearchHost},
	}, nil
}
