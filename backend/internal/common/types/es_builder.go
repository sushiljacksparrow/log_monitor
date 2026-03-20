package types

import "github.com/elastic/go-elasticsearch/v9/esutil"

type BulkIndexers struct {
	AuthBulkIndexer    esutil.BulkIndexer
	OrderBulkIndexer   esutil.BulkIndexer
	PaymentBulkIndexer esutil.BulkIndexer
}
