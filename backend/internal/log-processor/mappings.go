package logprocessor

const AuthServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":    { "type": "keyword" },
      "level":      { "type": "keyword" },
      "message":    { "type": "text" },
      "request_id": { "type": "keyword" },
      "user_id":    { "type": "keyword" },
      "ip":         { "type": "keyword" },
      "timestamp":  { "type": "date" }
    }
  }
}`

const OrderServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":    { "type": "keyword" },
      "level":      { "type": "keyword" },
      "message":    { "type": "text" },
      "request_id": { "type": "keyword" },
      "order_id":   { "type": "keyword" },
      "carrier":    { "type": "keyword" },
      "user_id":    { "type": "keyword" },
      "product_id": { "type": "keyword" },
      "stock_left": { "type": "integer" },
      "timestamp":  { "type": "date" }
    }
  }
}`

const PaymentServiceLogMapping = `{
  "mappings": {
    "properties": {
      "service":    { "type": "keyword" },
      "level":      { "type": "keyword" },
      "message":    { "type": "text" },
      "request_id": { "type": "keyword" },
      "payment_id": { "type": "keyword" },
      "gateway":    { "type": "keyword" },
      "order_id":   { "type": "keyword" },
      "amount":     { "type": "scaled_float,"scaling_factor":100 },
      "timestamp":  { "type": "date" }
    }
  }
}`
