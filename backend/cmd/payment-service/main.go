// payment-service
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mahirjain10/logflow/backend/internal/utils"
	logmonitor "github.com/mahirjain10/logflow/backend/pkg/log_monitor"
)

func main() {
	logger, err := logmonitor.New("payment-service",
		logmonitor.WithLogDir("/var/log/logflow"),
	)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	for {
		for _, mockLog := range MockPaymentLog {
			time.Sleep(3 * time.Second)

			paymentId, _ := utils.GenerateUUID()
			orderId, _ := utils.GenerateUUID()

			level := mockLog["level"].(string)
			msg := mockLog["message"].(string)

			fields := map[string]interface{}{
				"payment_id": paymentId,
				"order_id":   orderId,
				"gateway":    mockLog["gateway"],
				"amount":     mockLog["amount"],
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
