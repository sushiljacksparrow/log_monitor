package queryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"github.com/mahirjain10/logflow/backend/internal/constants"
	utils "github.com/mahirjain10/logflow/backend/internal/elasticsearch"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrLogsNotFound = errors.New("logs not found")

type Repository interface {
	SearchAuthLogs(ctx context.Context, filters AuthLogFilters, sortedValue *string, size int32) ([]*query.AuthLogs, *query.Pagination, error)
	SearchOrderLogs(ctx context.Context, filters OrderLogFilters, sortedValue *string, size int32) ([]*query.OrderLogs, *query.Pagination, error)
	SearchPaymentLogs(ctx context.Context, filters PaymentLogFilters, sortedValue *string, size int32) ([]*query.PaymentLogs, *query.Pagination, error)
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

func (r *repository) SearchAuthLogs(ctx context.Context, filters AuthLogFilters, sortedValue *string, size int32) ([]*query.AuthLogs, *query.Pagination, error) {
	queries := []types.Query{}
	searchAfter := []types.FieldValue{}
	sort := []types.SortCombinations{}
	appendTermQuery(&queries, "service", filters.Service)
	appendTermQuery(&queries, "level", filters.Level)
	appendMatchQuery(&queries, "message", filters.Message)
	appendTermQuery(&queries, "request_id", filters.RequestID)
	appendTermQuery(&queries, "user_id", filters.UserID)
	appendTermQuery(&queries, "ip", filters.IP)
	appendTimestampRange(&queries, filters.StartTimestamp, filters.EndTimestamp)
	if sortedValue != nil {
		decodedInteface, err := utils.DecodeSortedValue(*sortedValue)
		if err != nil {
			return nil, nil, fmt.Errorf("error while decoding sorted value %w", err)
		}
		appendSearchAfter(&searchAfter, decodedInteface)
		fmt.Println("SEARCH AFTER FINAL:", searchAfter)
		fmt.Printf("sorted value decoded: %+v\n", decodedInteface)
	}
	appendSortCombinations(&sort, "timestamp", sortorder.Desc)
	appendSortCombinations(&sort, "request_id", sortorder.Desc)
	fmt.Printf("query: %+v\n", queries)
	jbytes, _ := json.Marshal(queries)
	fmt.Printf("final query built: %+v\n", string(jbytes))
	return searchLogs[query.AuthLogs](ctx, r.esTypedClient, constants.AUTH_SERVICE_LOGS_INDEX, queries, searchAfter, sort, size)

}

func (r *repository) SearchOrderLogs(ctx context.Context, filters OrderLogFilters, sortedValue *string, size int32) ([]*query.OrderLogs, *query.Pagination, error) {
	queries := []types.Query{}
	searchAfter := []types.FieldValue{}
	sort := []types.SortCombinations{}
	appendTermQuery(&queries, "service", filters.Service)
	appendTermQuery(&queries, "level", filters.Level)
	appendMatchQuery(&queries, "message", filters.Message)
	appendTermQuery(&queries, "request_id", filters.RequestID)
	appendTermQuery(&queries, "user_id", filters.UserID)
	appendTermQuery(&queries, "order_id", filters.OrderID)
	appendTermQuery(&queries, "carrier", filters.Carrier)
	appendTermQuery(&queries, "product_id", filters.ProductID)
	appendTimestampRange(&queries, filters.StartTimestamp, filters.EndTimestamp)
	if sortedValue != nil {
		decodedInteface, err := utils.DecodeSortedValue(*sortedValue)
		if err != nil {
			return nil, nil, fmt.Errorf("error while decoding sorted value %w", err)
		}
		appendSearchAfter(&searchAfter, decodedInteface)
	}
	appendSortCombinations(&sort, "timestamp", sortorder.Desc)
	appendSortCombinations(&sort, "request_id", sortorder.Desc)
	fmt.Printf("final query built: %+v\n", queries)
	return searchLogs[query.OrderLogs](ctx, r.esTypedClient, constants.ORDER_SERVICE_LOGS_INDEX, queries, searchAfter, sort, size)
}

func (r *repository) SearchPaymentLogs(ctx context.Context, filters PaymentLogFilters, sortedValue *string, size int32) ([]*query.PaymentLogs, *query.Pagination, error) {
	queries := []types.Query{}
	searchAfter := []types.FieldValue{}
	sort := []types.SortCombinations{}
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
	if sortedValue != nil {
		decodedInteface, err := utils.DecodeSortedValue(*sortedValue)
		if err != nil {
			return nil, nil, fmt.Errorf("error while decoding sorted value %w", err)
		}
		appendSearchAfter(&searchAfter, decodedInteface)
		fmt.Printf("sorted value decoded: %+v\n", sortedValue)
	}
	appendSortCombinations(&sort, "timestamp", sortorder.Desc)
	appendSortCombinations(&sort, "request_id", sortorder.Desc)
	fmt.Printf("SORT STRUCT: %+v\n", sort)
	return searchLogs[query.PaymentLogs](ctx, r.esTypedClient, constants.PAYMENT_SERVICE_LOGS_INDEX, queries, searchAfter, sort, size)
}

func searchLogs[T any](ctx context.Context, client *es.TypedClient, indexName string, filters []types.Query, searchAfter []types.FieldValue, sort []types.SortCombinations, size int32) ([]*T, *query.Pagination, error) {
	if client == nil {
		return nil, nil, unsupportedRepositoryError()
	}
	pageSize := int(size)
	sizeInt := pageSize + 1
	req := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: filters,
			},
		},
		Sort: sort,
		Size: &sizeInt,
	}
	if len(searchAfter) > 0 {
		req.SearchAfter = searchAfter
	}
	jbytes, _ := json.Marshal(req)
	fmt.Printf("final query built: %+v\n", string(jbytes))
	res, err := client.Search().Index(indexName).Request(req).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Do(ctx)
	if err != nil {
		log.Printf("error while quering ES: %v\n", err)
		return nil, nil, fmt.Errorf("search %s: %w", indexName, err)
	}
	// data, _ := json.Marshal(res)
	// fmt.Println(string(data))
	resultCount := len(res.Hits.Hits)
	jb, _ := json.Marshal(res.Hits.Hits)
	fmt.Println("Printing all the hits : ", string(jb))
	logCount := resultCount
	if logCount > pageSize {
		logCount = pageSize
	}
	logs := make([]*T, 0, logCount)
	for i, hit := range res.Hits.Hits {
		if i >= pageSize {
			break
		}
		var entry T
		if err := json.Unmarshal(hit.Source_, &entry); err != nil {
			return nil, nil, fmt.Errorf("unmarshal %s hit: %w", indexName, err)
		}
		logs = append(logs, &entry)
	}
	var cursor types.Hit
	hasMore := resultCount > pageSize
	if hasMore && pageSize > 0 {
		// The cursor must point to the last hit returned to the client.
		cursor = res.Hits.Hits[pageSize-1]
	}
	var encodedString string
	if len(cursor.Sort) > 0 {
		encodedString, err = utils.EncodeSortedValue(cursor.Sort)
		if err != nil {
			return nil, nil, fmt.Errorf("error while encoding sorted value into string : %w", err)
		}
	}

	if len(logs) == 0 {
		return nil, nil, ErrLogsNotFound
	}
	fmt.Println("no error in repo")
	pagination := query.Pagination{
		HasMore:     hasMore,
		SortedValue: encodedString,
	}
	return logs, &pagination, nil
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

func appendSortCombinations(sort *[]types.SortCombinations, key string, sortOrder sortorder.SortOrder) {
	*sort = append(*sort, &types.SortOptions{SortOptions: map[string]types.FieldSort{
		key: {Order: &sortOrder},
	}})
}

func appendSearchAfter(searchAfter *[]types.FieldValue, value []types.FieldValue) {
	*searchAfter = append(*searchAfter, value...)
}
