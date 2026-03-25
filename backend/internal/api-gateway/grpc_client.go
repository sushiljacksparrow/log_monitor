package apigateway

import (
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Query query.QueryServiceClient
}

func InitGRPC() *GRPCClient {
	conn, _ := grpc.NewClient("localhost:50051")
	return &GRPCClient{
		Query: query.NewQueryServiceClient(conn),
	}
}
