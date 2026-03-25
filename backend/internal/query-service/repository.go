package queryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrLogsNotFound = errors.New("logs not found")

type Repository interface {
	SearchAuthLogs(ctx context.Context, filters AuthLogFilters) ([]*query.AuthLogs, error)
	SearchOrderLogs(ctx context.Context, filters OrderLogFilters) ([]*query.OrderLogs, error)
	SearchPaymentLogs(ctx context.Context, filters PaymentLogFilters) ([]*query.PaymentLogs, error)
}

type repository struct {
	es            *es.Client
	esTypedClient *es.TypedClient
}

type AuthLogFilters struct {
	Service        *string
	Level          *string
	Message        *string
	RequestID      *string
	UserID         *string
	IP             *string
	StartTimestamp *timestamppb.Timestamp
	EndTimestamp   *timestamppb.Timestamp
}

type OrderLogFilters struct {
	Service        *string
	Level          *string
	Message        *string
	RequestID      *string
	UserID         *string
	OrderID        *string
	Carrier        *string
	ProductID      *string
	StartTimestamp *timestamppb.Timestamp
	EndTimestamp   *timestamppb.Timestamp
}

type PaymentLogFilters struct {
	Service        *string
	Level          *string
	Message        *string
	RequestID      *string
	OrderID        *string
	PaymentID      *string
	Gateway        *string
	Amount         *float64
	StartTimestamp *timestamppb.Timestamp
	EndTimestamp   *timestamppb.Timestamp
}

func NewRepository(esClient *es.Client, esTypedClient *es.TypedClient) Repository {
	return &repository{
		es:            esClient,
		esTypedClient: esTypedClient,
	}
}

func (r *repository) SearchAuthLogs(ctx context.Context, filters AuthLogFilters) ([]*query.AuthLogs, error) {
	queries := []types.Query{}
	appendTermQuery(&queries, "service", filters.Service)
	appendTermQuery(&queries, "level", filters.Level)
	appendMatchQuery(&queries, "message", filters.Message)
	appendTermQuery(&queries, "request_id", filters.RequestID)
	appendTermQuery(&queries, "user_id", filters.UserID)
	appendTermQuery(&queries, "ip", filters.IP)
	appendTimestampRange(&queries, filters.StartTimestamp, filters.EndTimestamp)

	return searchLogs[query.AuthLogs](ctx, r.esTypedClient, "auth-service-logs", queries)
}

func (r *repository) SearchOrderLogs(ctx context.Context, filters OrderLogFilters) ([]*query.OrderLogs, error) {
	queries := []types.Query{}
	appendTermQuery(&queries, "service", filters.Service)
	appendTermQuery(&queries, "level", filters.Level)
	appendMatchQuery(&queries, "message", filters.Message)
	appendTermQuery(&queries, "request_id", filters.RequestID)
	appendTermQuery(&queries, "user_id", filters.UserID)
	appendTermQuery(&queries, "order_id", filters.OrderID)
	appendTermQuery(&queries, "carrier", filters.Carrier)
	appendTermQuery(&queries, "product_id", filters.ProductID)
	appendTimestampRange(&queries, filters.StartTimestamp, filters.EndTimestamp)

	return searchLogs[query.OrderLogs](ctx, r.esTypedClient, "order-service-logs", queries)
}

func (r *repository) SearchPaymentLogs(ctx context.Context, filters PaymentLogFilters) ([]*query.PaymentLogs, error) {
	queries := []types.Query{}
	appendTermQuery(&queries, "service", filters.Service)
	appendTermQuery(&queries, "level", filters.Level)
	appendMatchQuery(&queries, "message", filters.Message)
	appendTermQuery(&queries, "request_id", filters.RequestID)
	appendTermQuery(&queries, "order_id", filters.OrderID)
	appendTermQuery(&queries, "payment_id", filters.PaymentID)
	appendTermQuery(&queries, "gateway", filters.Gateway)
	if filters.Amount != nil {
		queries = append(queries, types.Query{
			Term: map[string]types.TermQuery{
				"amount": {Value: *filters.Amount},
			},
		})
	}
	appendTimestampRange(&queries, filters.StartTimestamp, filters.EndTimestamp)

	return searchLogs[query.PaymentLogs](ctx, r.esTypedClient, "payment-service-logs", queries)
}

func searchLogs[T any](ctx context.Context, client *es.TypedClient, indexName string, filters []types.Query) ([]*T, error) {
	if client == nil {
		return nil, unsupportedRepositoryError()
	}

	res, err := client.Search().
		Index(indexName).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Filter: filters,
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("search %s: %w", indexName, err)
	}

	logs := make([]*T, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var entry T
		if err := json.Unmarshal(hit.Source_, &entry); err != nil {
			return nil, fmt.Errorf("unmarshal %s hit: %w", indexName, err)
		}
		logs = append(logs, &entry)
	}

	if len(logs) == 0 {
		return nil, ErrLogsNotFound
	}

	return logs, nil
}

func appendTermQuery(filters *[]types.Query, field string, value *string) {
	if value == nil || *value == "" {
		return
	}
	*filters = append(*filters, types.Query{
		Term: map[string]types.TermQuery{
			field: {Value: *value},
		},
	})
}

func appendMatchQuery(filters *[]types.Query, field string, value *string) {
	if value == nil || *value == "" {
		return
	}
	*filters = append(*filters, types.Query{
		Match: map[string]types.MatchQuery{
			field: {Query: *value},
		},
	})
}

func appendTimestampRange(filters *[]types.Query, start, end *timestamppb.Timestamp) {
	if start == nil || end == nil {
		return
	}

	startValue := start.AsTime().Format(time.RFC3339)
	endValue := end.AsTime().Format(time.RFC3339)

	*filters = append(*filters, types.Query{
		Range: map[string]types.RangeQuery{
			"timestamp": types.DateRangeQuery{
				Gte: &startValue,
				Lte: &endValue,
			},
		},
	})
}
