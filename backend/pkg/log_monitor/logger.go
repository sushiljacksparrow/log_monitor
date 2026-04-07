package log_monitor

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/go-uuid"
)

// Logger is the public API for the log_monitor SDK.
// Create one per service with New(), then call Info/Warn/Error/Debug.
type Logger struct {
	config LoggerConfig
	writer *fileWriter
}

// New creates a new Logger for the given service.
//
//	logger, err := log_monitor.New("auth-service",
//	    log_monitor.WithLogDir("/tmp/logs"),
//	    log_monitor.WithLevel(log_monitor.INFO),
//	)
func New(serviceName string, opts ...Option) (*Logger, error) {
	cfg := newDefaultConfig(serviceName)
	for _, opt := range opts {
		opt(&cfg)
	}

	w, err := newFileWriter(cfg.LogDir, cfg.ServiceName, cfg.MaxFileSize)
	if err != nil {
		return nil, err
	}

	return &Logger{config: cfg, writer: w}, nil
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log(INFO, msg, fields)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log(WARN, msg, fields)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.log(ERROR, msg, fields)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.log(DEBUG, msg, fields)
}

// Close flushes and closes the underlying log file.
func (l *Logger) Close() error {
	return l.writer.Close()
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	if level < l.config.Level {
		return
	}

	reqID, _ := uuid.GenerateUUID()

	entry := LogEntry{
		Service:   l.config.ServiceName,
		Level:     level.String(),
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: reqID,
		Fields:    fields,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	// Append newline — one JSON line per log entry
	data = append(data, '\n')
	l.writer.Write(data)
}
