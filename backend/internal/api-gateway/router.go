package apigateway

import (
	"github.com/gin-gonic/gin"
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

	engine.GET("/ws", func(c *gin.Context) {
		livews.ServeWS(hub, c.Writer, c.Request)
	})
}
