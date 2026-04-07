package log_monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// fileWriter handles thread-safe file writing with size-based rotation.
type fileWriter struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	logDir      string
	serviceName string
	maxFileSize int64
	currentSize int64
}

func newFileWriter(logDir, serviceName string, maxFileSize int64) (*fileWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}

	filePath := filepath.Join(logDir, serviceName+".log")

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &fileWriter{
		file:        f,
		filePath:    filePath,
		logDir:      logDir,
		serviceName: serviceName,
		maxFileSize: maxFileSize,
		currentSize: info.Size(),
	}, nil
}

// Write writes data to the log file, rotating if the size limit is exceeded.
func (w *fileWriter) Write(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(data)) > w.maxFileSize {
		if err := w.rotate(); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	n, err := w.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}
	w.currentSize += int64(n)
	return nil
}

// rotate renames the current file and opens a new one.
func (w *fileWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	rotatedName := fmt.Sprintf("%s-%s.log", w.serviceName, time.Now().UTC().Format("20060102T150405.000"))
	rotatedPath := filepath.Join(w.logDir, rotatedName)

	if err := os.Rename(w.filePath, rotatedPath); err != nil {
		return err
	}

	f, err := os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.file = f
	w.currentSize = 0
	return nil
}

// Close closes the underlying file.
func (w *fileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}
