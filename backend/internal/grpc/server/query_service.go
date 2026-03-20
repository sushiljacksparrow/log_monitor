package server

import (
	"context"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/mahirjain10/logflow/backend/internal/elasticsearch/wrapper"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
)

type QueryServer struct {
	query.UnimplementedQueryServiceServer
	es *es.Client
}

func NewQueryServer(esClient *es.Client) *QueryServer {
	return &QueryServer{es: esClient}
}

func (qs *QueryServer) GetIndexesWithMapping(ctx context.Context, request *query.GetIndexesWithMappingRequest) (*query.GetIndexesWithMappingResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	res, err := wrapper.GetIndexes(ctxWithTimeout, qs.es)
	if err != nil {
		return nil, err
	}
	response := &query.GetIndexesWithMappingResponse{
		Indexes: ,
	}
	return
}
