package main

import (
	"log"
	"net"

	"github.com/mahirjain10/logflow/backend/internal/config"
	"github.com/mahirjain10/logflow/backend/internal/elasticsearch"
	pb "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	queryservice "github.com/mahirjain10/logflow/backend/internal/query-service"
	"google.golang.org/grpc"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error while initalizing config: %v", err)
	}
	esClient, esTypedClient, err := elasticsearch.InitES(config)
	if err != nil {
		log.Fatal(err)
	}
	// 1. open a TCP port
	//
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. create the gRPC server
	grpcServer := grpc.NewServer()

	// 3. plug your impl into the generated stub
	pb.RegisterQueryServiceServer(grpcServer, queryservice.New(esClient, esTypedClient))

	// 4. start serving
	log.Println("query-service gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
