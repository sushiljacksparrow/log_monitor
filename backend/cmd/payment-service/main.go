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
		log.Printf("couldnt intialize kafka producer: %v", err)
	}
	defer producer.Close()
	for {
		for _, logValue := range MockPaymentLog {
			time.Sleep(3 * time.Second)
			logValue["timestamp"] = time.Now().UTC().Format(time.RFC3339)
			reqID, err := utils.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			paymentId, err := utils.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			logValue["requestId"] = reqID
			logValue["paymentId"] = paymentId
			logValue["ip"] = utils.RandomIP()
			logByte, err := json.Marshal(logValue)
			if err != nil {
				log.Printf("error while marshalling log %v", err)
			}
			producer.Publish(envs.KafkaTopicPaymentLog, logByte)
			fmt.Println(logValue)
		}
	}
}
