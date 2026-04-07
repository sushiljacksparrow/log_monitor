package shipper

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// FilePosition tracks where the shipper left off in a log file.
type FilePosition struct {
	Inode    uint64    `json:"inode"`
	Offset   int64     `json:"offset"`
	LastRead time.Time `json:"last_read"`
}

// PositionTracker persists file read positions to survive restarts.
type PositionTracker struct {
	mu        sync.Mutex
	positions map[string]FilePosition
	filePath  string
}

func NewPositionTracker(filePath string) *PositionTracker {
	return &PositionTracker{
		positions: make(map[string]FilePosition),
		filePath:  filePath,
	}
}

func (pt *PositionTracker) Get(path string) (FilePosition, bool) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pos, ok := pt.positions[path]
	return pos, ok
}

func (pt *PositionTracker) Update(path string, offset int64, inode uint64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.positions[path] = FilePosition{
		Inode:    inode,
		Offset:   offset,
		LastRead: time.Now().UTC(),
	}
}

// Save atomically writes positions to disk (write to .tmp, then rename).
func (pt *PositionTracker) Save() error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	data, err := json.MarshalIndent(pt.positions, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := pt.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, pt.filePath)
}

// Load reads positions from disk.
func (pt *PositionTracker) Load() error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	data, err := os.ReadFile(pt.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no position file yet, start fresh
		}
		return err
	}

	return json.Unmarshal(data, &pt.positions)
}
