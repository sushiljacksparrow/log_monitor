package queryservice

import (
	"context"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v9"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	query.UnimplementedQueryServiceServer
	repo Repository
}

func New(esClient *es.Client, esTypedClient *es.TypedClient) *Service {
	return &Service{
		repo: NewRepository(esClient, esTypedClient),
	}
}

func (s *Service) AuthLogs(ctx context.Context, req *query.AuthLogsRequest) (*query.AuthLogsResponse, error) {
	fmt.Printf("auth request %+v\n", req)
	logs, page, err := s.repo.SearchAuthLogs(ctx, AuthLogFilters{
		Service:        req.Service,
		Level:          req.Level,
		Message:        req.Message,
		RequestID:      req.RequestId,
		UserID:         req.UserId,
		IP:             req.Ip,
		StartTimestamp: req.StartTimestamp,
		EndTimestamp:   req.EndTimestamp,
	}, req.SortedValue, req.Size)
	if err != nil {
		return nil, toStatusError("auth logs", err)
	}
	fmt.Println("no error")
	response := &query.AuthLogsResponse{Logs: logs, BaseResponse: page}
	fmt.Printf("response from auhtlog: %+v\n", response)
	return response, nil
}

func (s *Service) PaymentLogs(ctx context.Context, req *query.PaymentLogsRequest) (*query.PaymentLogsResponse, error) {
	logs, page, err := s.repo.SearchPaymentLogs(ctx, PaymentLogFilters{
		Service:        req.Service,
		Level:          req.Level,
		Message:        req.Message,
		RequestID:      req.RequestId,
		OrderID:        req.OrderId,
		PaymentID:      req.PaymentId,
		Gateway:        req.Gateway,
		Amount:         req.Amount,
		StartTimestamp: req.StartTimestamp,
		EndTimestamp:   req.EndTimestamp,
	}, req.SortedValue, req.Size)
	if err != nil {
		return nil, toStatusError("payment logs", err)
	}

	return &query.PaymentLogsResponse{Logs: logs, BaseResponse: page}, nil
}

func (s *Service) OrderLogs(ctx context.Context, req *query.OrderLogsRequest) (*query.OrderLogsResponse, error) {
	logs, page, err := s.repo.SearchOrderLogs(ctx, OrderLogFilters{
		Service:        req.Service,
		Level:          req.Level,
		Message:        req.Message,
		RequestID:      req.RequestId,
		UserID:         req.UserId,
		OrderID:        req.OrderId,
		Carrier:        req.Carrier,
		ProductID:      req.ProductId,
		StartTimestamp: req.StartTimestamp,
		EndTimestamp:   req.EndTimestamp,
	}, req.SortedValue, req.Size)
	if err != nil {
		return nil, toStatusError("order logs", err)
	}

	return &query.OrderLogsResponse{Logs: logs, BaseResponse: page}, nil
}

func toStatusError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if err == ErrLogsNotFound {
		return status.Error(codes.NotFound, "logs not found with given filters")
	}
	return status.Errorf(codes.Internal, "error querying %s: %v", operation, err)
}

func unsupportedRepositoryError() error {
	return status.Error(codes.Internal, "query repository is not configured")
}
