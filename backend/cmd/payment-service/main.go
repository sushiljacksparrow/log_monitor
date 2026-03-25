package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	"github.com/mahirjain10/logflow/backend/internal/ids"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
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
			reqID, err := ids.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			paymentId, err := ids.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			orderId, err := ids.GenerateUUID()
			if err != nil {
				log.Println(err)
			}
			logValue["request_id"] = reqID
			logValue["payment_id"] = paymentId
			logValue["order_id"] = orderId
			logByte, err := json.Marshal(logValue)
			if err != nil {
				log.Printf("error while marshalling log %v", err)
			}
			producer.Publish(constants.PAYMENT_SERVICE_LOGS_TOPIC, logByte)
			fmt.Println(logValue)
		}
	}
}
