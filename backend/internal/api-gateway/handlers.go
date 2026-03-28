package apigateway

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	errmap "github.com/mahirjain10/logflow/backend/internal/errmap"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthLogRequest struct {
	Service        *string    `json:"service"`
	Level          *string    `json:"level"`
	Message        *string    `json:"message"`
	RequestId      *string    `json:"request_id"`
	UserId         *string    `json:"user_id"`
	Ip             *string    `json:"ip"`
	StartTimestamp *time.Time `json:"start_timestamp"`
	EndTimestamp   *time.Time `json:"end_timestamp"`
}

type OrderLogRequest struct {
	Service        *string    `json:"service"`
	Level          *string    `json:"level"`
	Message        *string    `json:"message"`
	RequestId      *string    `json:"request_id"`
	UserId         *string    `json:"user_id"`
	OrderId        *string    `json:"order_id"`
	Carrier        *string    `json:"carrier"`
	ProductId      *string    `json:"product_id"`
	StartTimestamp *time.Time `json:"start_timestamp"`
	EndTimestamp   *time.Time `json:"end_timestamp"`
}

type PaymentLogRequest struct {
	Service        *string    `json:"service"`
	Level          *string    `json:"level"`
	Message        *string    `json:"message"`
	RequestId      *string    `json:"request_id"`
	OrderId        *string    `json:"order_id"`
	PaymentId      *string    `json:"payment_id"`
	Gateway        *string    `json:"gateway"`
	Amount         *string    `json:"amount"`
	StartTimestamp *time.Time `json:"start_timestamp"`
	EndTimestamp   *time.Time `json:"end_timestamp"`
}

func SearchAuthLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v", err)
			c.JSON(http.StatusBadRequest, "invalid request body")
			return
		}

		authReq := query.AuthLogsRequest{
			Service:   req.Service,
			Level:     req.Level,
			Message:   req.Message,
			RequestId: req.RequestId,
			UserId:    req.UserId,
			Ip:        req.Ip,
		}

		if req.StartTimestamp != nil {
			authReq.StartTimestamp = timestamppb.New(*req.StartTimestamp)
		}

		if req.EndTimestamp != nil {
			authReq.EndTimestamp = timestamppb.New(*req.EndTimestamp)
		}

		resp, err := grpcClient.Query.AuthLogs(ctx, &authReq)
		if err != nil {
			code, msg := errmap.GRPCToHTTP(err)
			c.JSON(code, msg)
			return
		}

		if len(resp.Logs) == 0 {
			c.JSON(http.StatusNotFound, "no logs found")
			return
		}

		c.JSON(http.StatusOK, resp.Logs)
	}
}

func SearchOrderLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v", err)
			c.JSON(http.StatusBadRequest, "invalid request body")
			return
		}

		orderReq := query.OrderLogsRequest{
			Service:   req.Service,
			Level:     req.Level,
			Message:   req.Message,
			RequestId: req.RequestId,
			UserId:    req.UserId,
			OrderId:   req.OrderId,
			Carrier:   req.Carrier,
			ProductId: req.ProductId,
		}

		if req.StartTimestamp != nil {
			orderReq.StartTimestamp = timestamppb.New(*req.StartTimestamp)
		}

		if req.EndTimestamp != nil {
			orderReq.EndTimestamp = timestamppb.New(*req.EndTimestamp)
		}

		resp, err := grpcClient.Query.OrderLogs(ctx, &orderReq)
		if err != nil {
			code, msg := errmap.GRPCToHTTP(err)
			c.JSON(code, msg)
			return
		}

		if len(resp.Logs) == 0 {
			c.JSON(http.StatusNotFound, "no logs found")
			return
		}

		c.JSON(http.StatusOK, resp.Logs)
	}
}

func SearchPaymentLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PaymentLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v", err)
			c.JSON(http.StatusBadRequest, "invalid request body")
			return
		}

		paymentReq := query.PaymentLogsRequest{
			Service:   req.Service,
			Level:     req.Level,
			Message:   req.Message,
			RequestId: req.RequestId,
			OrderId:   req.OrderId,
			PaymentId: req.PaymentId,
			Gateway:   req.Gateway,
			Amount:    req.Amount,
		}

		if req.StartTimestamp != nil {
			paymentReq.StartTimestamp = timestamppb.New(*req.StartTimestamp)
		}

		if req.EndTimestamp != nil {
			paymentReq.EndTimestamp = timestamppb.New(*req.EndTimestamp)
		}

		resp, err := grpcClient.Query.PaymentLogs(ctx, &paymentReq)
		if err != nil {
			code, msg := errmap.GRPCToHTTP(err)
			c.JSON(code, msg)
			return
		}

		if len(resp.Logs) == 0 {
			c.JSON(http.StatusNotFound, "no logs found")
			return
		}

		c.JSON(http.StatusOK, resp.Logs)
	}
}
