package apigateway

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	errmap "github.com/mahirjain10/logflow/backend/internal/errmap"
	query "github.com/mahirjain10/logflow/backend/internal/grpc/gen"
	"github.com/mahirjain10/logflow/backend/internal/helper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthLogRequest struct {
	Service        *string    `json:"service"`
	Level          *string    `json:"level"`
	Message        *string    `json:"message"`
	RequestId      *string    `json:"request_id"`
	UserId         *string    `json:"user_id"`
	Ip             *string    `json:"ip"`
	SortedValue    *string    `json:"sorted_value"`
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
	SortedValue    *string    `json:"sorted_value"`
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
	SortedValue    *string    `json:"sorted_value"`
	StartTimestamp *time.Time `json:"start_timestamp"`
	EndTimestamp   *time.Time `json:"end_timestamp"`
}

func SearchAuthLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var size int32
		if v := c.Query("size"); len(v) == 0 {
			size = 10
		} else {
			s, err := strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "internal server error")
				return
			}
			size = int32(s)
			if size > 100 {
				c.JSON(http.StatusBadRequest, "size can't be more than 100")
				return
			}
		}
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v\n", err)
			c.JSON(http.StatusBadRequest, "invalid request body")
			return
		}

		authReq := query.AuthLogsRequest{
			Service:     req.Service,
			Level:       req.Level,
			Message:     req.Message,
			RequestId:   req.RequestId,
			UserId:      req.UserId,
			Ip:          req.Ip,
			Size:        size,
			SortedValue: req.SortedValue,
		}

		if req.StartTimestamp != nil {
			authReq.StartTimestamp = timestamppb.New(*req.StartTimestamp)
		}

		if req.EndTimestamp != nil {
			authReq.EndTimestamp = timestamppb.New(*req.EndTimestamp)
		}
		log.Println("grpcClient:", grpcClient)
		log.Println("grpcClient.Query:", grpcClient.Query)
		resp, err := grpcClient.Query.AuthLogs(ctx, &authReq)
		if err != nil {
			log.Printf("error while calling grpc call %v\n", err)
			code, msg := errmap.GRPCToHTTP(err)
			helper.SendResponse(c, code, msg, nil)
			return
		}

		// if len(resp.Logs) == 0 {

		// 	c.JSON(http.StatusNotFound, "no logs found")
		// 	return
		// }
		helper.SendResponse(c, http.StatusOK, "data retrieved successfully", resp)
	}
}

func SearchOrderLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var size int32
		if v := c.Query("size"); len(v) == 0 {
			size = 10
		} else {
			s, err := strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "internal server error")
				return
			}
			size = int32(s)
			if size > 100 {
				c.JSON(http.StatusBadRequest, "size can't be more than 100")
				return
			}
		}
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v\n", err)
			helper.SendResponse(c, http.StatusBadRequest, "invalid request body", nil)
			return
		}

		orderReq := query.OrderLogsRequest{
			Service:     req.Service,
			Level:       req.Level,
			Message:     req.Message,
			RequestId:   req.RequestId,
			UserId:      req.UserId,
			OrderId:     req.OrderId,
			Carrier:     req.Carrier,
			ProductId:   req.ProductId,
			SortedValue: req.SortedValue,
			Size:        size,
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
			helper.SendResponse(c, code, msg, nil)
			return
		}

		// if len(resp.Logs) == 0 {
		// 	c.JSON(http.StatusNotFound, "no logs found")
		// 	return
		// }

		helper.SendResponse(c, http.StatusOK, "data retrieved successfully", resp)
	}
}

func SearchPaymentLogs(grpcClient *GRPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PaymentLogRequest
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var size int32
		if v := c.Query("size"); len(v) == 0 {
			size = 10
		} else {
			s, err := strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "internal server error")
				return
			}
			size = int32(s)
			if size > 100 {
				c.JSON(http.StatusBadRequest, "size can't be more than 100")
				return
			}
		}

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Printf("error while decoding json: %v\n", err)
			c.JSON(http.StatusBadRequest, "invalid request body")
			return
		}

		paymentReq := query.PaymentLogsRequest{
			Service:     req.Service,
			Level:       req.Level,
			Message:     req.Message,
			RequestId:   req.RequestId,
			OrderId:     req.OrderId,
			PaymentId:   req.PaymentId,
			Gateway:     req.Gateway,
			Amount:      req.Amount,
			Size:        size,
			SortedValue: req.SortedValue,
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
			helper.SendResponse(c, code, msg, nil)
			return
		}

		// if len(resp.Logs) == 0 {
		// 	c.JSON(http.StatusNotFound, "no logs found")
		// 	return
		// }

		helper.SendResponse(c, http.StatusOK, "data retrieved successfully", resp)
	}
}
