package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	apigateway "github.com/mahirjain10/logflow/backend/internal/api-gateway"
	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/kafka"
	livews "github.com/mahirjain10/logflow/backend/internal/websocket"
)

func main() {
	envs, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	grpcClient, err := apigateway.InitGRPC()
	if err != nil {
		log.Fatalf("query grpc client failed: %v", err)
	}
	hub := livews.NewHub()
	go hub.Run()

	consumer, err := kafka.NewConsumer(envs.KafkaBrokers, "logs-live-1")
	if err != nil {
		log.Printf("could not initialize live-log consumer: %v\n", err)
	} else {
		go apigateway.StartLiveLogConsumer(context.Background(), consumer, hub)
	}

	log.Println("grpcClient:", grpcClient)
	log.Println("grpcClient.Query:", grpcClient.Query)
	apigateway.RegisterRoutes(router, grpcClient, hub)
	// ssss
	log.Printf("api-gateway running on port: 8000\n")
	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
