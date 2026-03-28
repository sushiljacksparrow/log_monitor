package apigateway

import (
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Query query.QueryServiceClient
}

func InitGRPC() (*GRPCClient, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		Query: query.NewQueryServiceClient(conn),
	}, nil
}
