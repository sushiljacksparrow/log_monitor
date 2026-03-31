package apigateway

import (
	"os"

	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Query query.QueryServiceClient
}

func InitGRPC() (*GRPCClient, error) {
	queryAddr := os.Getenv("QUERY_SERVICE_ADDR")
	if queryAddr == "" {
		queryAddr = "localhost:50051"
	}

	conn, err := grpc.NewClient(queryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		Query: query.NewQueryServiceClient(conn),
	}, nil
}
