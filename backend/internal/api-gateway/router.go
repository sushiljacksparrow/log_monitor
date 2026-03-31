package apigateway

import (
	"github.com/gin-gonic/gin"
	"github.com/mahirjain10/logflow/backend/internal/helper"
	livews "github.com/mahirjain10/logflow/backend/internal/websocket"
)

func RegisterRoutes(engine *gin.Engine, grpcClient *GRPCClient, hub *livews.Hub) {
	api := engine.Group("/api")

	search := api.Group("/search")
	{
		search.POST("/auth-service", SearchAuthLogs(grpcClient))
		search.POST("/order-service", SearchOrderLogs(grpcClient))
		search.POST("/payment-service", SearchPaymentLogs(grpcClient))
	}

	engine.GET("/health", func(c *gin.Context) {
		helper.SendResponse(c, 200, "gateway healthy", gin.H{"service": "api-gateway", "ok": true})
	})

	engine.GET("/ws", func(c *gin.Context) {
		livews.ServeWS(hub, c.Writer, c.Request)
	})
}
