package shipper

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ShipperConfig struct {
	WatchDir     string
	PositionFile string
	KafkaBrokers []string
	PollInterval time.Duration
}

func LoadConfig() (ShipperConfig, error) {
	watchDir := os.Getenv("LOGFLOW_WATCH_DIR")
	if watchDir == "" {
		watchDir = "/var/log/logflow"
	}

	positionFile := os.Getenv("LOGFLOW_POSITION_FILE")
	if positionFile == "" {
		positionFile = "/var/log/logflow/.positions.json"
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return ShipperConfig{}, fmt.Errorf("KAFKA_BROKERS env var is required")
	}

	return ShipperConfig{
		WatchDir:     watchDir,
		PositionFile: positionFile,
		KafkaBrokers: strings.Split(brokers, ","),
		PollInterval: 30 * time.Second,
	}, nil
}
