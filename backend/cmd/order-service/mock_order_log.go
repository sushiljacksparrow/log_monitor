package main

var MockOrderLog = []map[string]interface{}{
	{
		"service":   "order-service",
		"level":     "INFO",
		"message":   "Order created successfully",
		"requestId": "req-101",
		"userId":    "user-101",
		"orderId":   "order-1001",
		"timestamp": "",
	},
	{
		"service":   "order-service",
		"level":     "DEBUG",
		"message":   "Inventory check started",
		"requestId": "req-102",
		"orderId":   "order-1001",
		"productId": "prod-501",
		"timestamp": "",
	},
	{
		"service":   "order-service",
		"level":     "WARN",
		"message":   "Low stock detected",
		"requestId": "req-103",
		"productId": "prod-501",
		"stockLeft": 3,
		"timestamp": "",
	},
	{
		"service":   "order-service",
		"level":     "ERROR",
		"message":   "Order failed due to insufficient stock",
		"requestId": "req-104",
		"orderId":   "order-1002",
		"productId": "prod-999",
		"timestamp": "",
	},
	{
		"service":   "order-service",
		"level":     "INFO",
		"message":   "Order shipped successfully",
		"requestId": "req-105",
		"orderId":   "order-1003",
		"carrier":   "BlueDart",
		"timestamp": "",
	},
}
