package log_monitor

import "encoding/json"

// LogEntry represents a single structured log line.
type LogEntry struct {
	Service   string                 `json:"service"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"request_id"`
	Fields    map[string]interface{} `json:"-"`
}

// MarshalJSON flattens Fields into top-level JSON alongside core fields.
// This matches the existing log format where service-specific fields
// (user_id, ip, amount, etc.) appear at the top level, not nested.
func (e LogEntry) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 5+len(e.Fields))
	m["service"] = e.Service
	m["level"] = e.Level
	m["message"] = e.Message
	m["timestamp"] = e.Timestamp
	m["request_id"] = e.RequestID
	for k, v := range e.Fields {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
	return json.Marshal(m)
}
