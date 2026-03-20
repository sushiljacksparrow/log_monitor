package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
	"github.com/mahirjain10/logflow/backend/internal/utils"
)

func main() {
	envs, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	producer, err := kafka.NewProducer(envs.KafkaBrokers)
	if err != nil {
		log.Println("couldnt intialize kafka producer")
	}
	defer producer.Close()
	for {
		for _, logValue := range MockOrderLog {
			time.Sleep(3 * time.Second)
			logValue["timestamp"] = time.Now().UTC().Format(time.RFC3339)
			reqID, err := utils.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			orderId, err := utils.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			logValue["requestId"] = reqID
			logValue["orderId"] = orderId
			logValue["ip"] = utils.RandomIP()
			logByte, err := json.Marshal(logValue)
			if err != nil {
				log.Printf("error while marshalling log %v", err)
			}
			producer.Publish(envs.KafkaTopicOrderLog, logByte)
			fmt.Println(logValue)
		}
	}
}
