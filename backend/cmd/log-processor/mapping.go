package main

// ES index name
const (
	AUTH_SERVICE_LOGS    = "auth-service-logs"
	ORDER_SERVICE_LOGS   = "order-service-logs"
	PAYMENT_SERVICE_LOGS = "payment-service-logs"
)

var IndexMappings = map[string]string{
	AUTH_SERVICE_LOGS:    AuthServiceLogMapping,
	ORDER_SERVICE_LOGS:   OrderServiceLogMapping,
	PAYMENT_SERVICE_LOGS: PaymentServiceLogMapping,
}

const AuthServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":   { "type": "keyword" },
      "level":     { "type": "keyword" },
      "message":   { "type": "text" },
      "requestId": { "type": "keyword" },
      "userId":    { "type": "keyword" },
      "ip":        { "type": "keyword" },
      "timestamp": { "type": "date" }
    }
		}
	}`

const OrderServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":   { "type": "keyword" },
      "level":     { "type": "keyword" },
      "message":   { "type": "text" },
      "requestId": { "type": "keyword" },
      "orderId":    { "type": "keyword" },
      "carrier":        { "type": "keyword" },
      "timestamp": { "type": "date" }
    }
		}
	}`

const PaymentServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":   { "type": "keyword" },
      "level":     { "type": "keyword" },
      "message":   { "type": "text" },
      "requestId": { "type": "keyword" },
      "paymentId":    { "type": "keyword" },
      "gateway":        { "type": "keyword" },
      "timestamp": { "type": "date" }
    }
		}
	}`
