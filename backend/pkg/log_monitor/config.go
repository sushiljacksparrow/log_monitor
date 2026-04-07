package log_monitor

// Level represents log severity.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

const (
	defaultLogDir      = "/var/log/logflow"
	defaultMaxFileSize = 100 * 1024 * 1024 // 100MB
)

// LoggerConfig holds configuration for the logger.
type LoggerConfig struct {
	ServiceName string
	LogDir      string
	MaxFileSize int64
	Level       Level
}

// Option is a functional option for configuring the logger.
type Option func(*LoggerConfig)

func WithLogDir(dir string) Option {
	return func(c *LoggerConfig) {
		c.LogDir = dir
	}
}

func WithMaxFileSize(size int64) Option {
	return func(c *LoggerConfig) {
		c.MaxFileSize = size
	}
}

func WithLevel(level Level) Option {
	return func(c *LoggerConfig) {
		c.Level = level
	}
}

func newDefaultConfig(serviceName string) LoggerConfig {
	return LoggerConfig{
		ServiceName: serviceName,
		LogDir:      defaultLogDir,
		MaxFileSize: defaultMaxFileSize,
		Level:       DEBUG,
	}
}
