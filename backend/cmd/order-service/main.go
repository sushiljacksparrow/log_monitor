// order-service
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mahirjain10/logflow/backend/internal/utils"
	logmonitor "github.com/mahirjain10/logflow/backend/pkg/log_monitor"
)

func main() {
	logger, err := logmonitor.New("order-service",
		logmonitor.WithLogDir("/var/log/logflow"),
	)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	for {
		for _, mockLog := range MockOrderLog {
			time.Sleep(3 * time.Second)

			orderId, _ := utils.GenerateUUID()
			ip := utils.RandomIP()

			level := mockLog["level"].(string)
			msg := mockLog["message"].(string)

			fields := map[string]interface{}{
				"user_id":    mockLog["user_id"],
				"order_id":   orderId,
				"product_id": mockLog["product_id"],
				"stock_left": mockLog["stock_left"],
				"carrier":    mockLog["carrier"],
				"ip":         ip.String(),
			}

			switch level {
			case "INFO":
				logger.Info(msg, fields)
			case "WARN":
				logger.Warn(msg, fields)
			case "ERROR":
				logger.Error(msg, fields)
			case "DEBUG":
				logger.Debug(msg, fields)
			}

			fmt.Printf("logged: level=%s message=%s\n", level, msg)
		}
	}
}
