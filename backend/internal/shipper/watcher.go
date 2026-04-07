package shipper

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a directory for .log files and starts Tailers for each one.
type Watcher struct {
	dir       string
	position  *PositionTracker
	publisher *Publisher
	tailers   map[string]context.CancelFunc
	mu        sync.Mutex
	pollInterval time.Duration
}

func NewWatcher(cfg ShipperConfig, pos *PositionTracker, pub *Publisher) *Watcher {
	return &Watcher{
		dir:          cfg.WatchDir,
		position:     pos,
		publisher:    pub,
		tailers:      make(map[string]context.CancelFunc),
		pollInterval: cfg.PollInterval,
	}
}

// Run starts watching the directory. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	// Start tailers for all existing .log files.
	w.scanDirectory(ctx)

	// Set up fsnotify watcher.
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fsWatcher.Close()

	if err := fsWatcher.Add(w.dir); err != nil {
		return err
	}

	log.Printf("watching directory %s for .log files\n", w.dir)

	// Periodic fallback scan (fsnotify can miss events on Docker volumes).
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.stopAllTailers()
			return ctx.Err()

		case event, ok := <-fsWatcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Create == fsnotify.Create && isLogFile(event.Name) {
				w.startTailer(ctx, event.Name)
			}

		case err, ok := <-fsWatcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("fsnotify error: %v\n", err)

		case <-ticker.C:
			w.scanDirectory(ctx)
		}
	}
}

func (w *Watcher) scanDirectory(ctx context.Context) {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		log.Printf("failed to scan directory %s: %v\n", w.dir, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(w.dir, entry.Name())
		if isLogFile(path) {
			w.startTailer(ctx, path)
		}
	}
}

func (w *Watcher) startTailer(ctx context.Context, path string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Don't start a duplicate tailer.
	if _, exists := w.tailers[path]; exists {
		return
	}

	tailer, err := NewTailer(path, w.position, w.publisher)
	if err != nil {
		log.Printf("failed to create tailer for %s: %v\n", path, err)
		return
	}

	tailerCtx, cancel := context.WithCancel(ctx)
	w.tailers[path] = cancel

	go tailer.Run(tailerCtx)
	log.Printf("started tailer for %s\n", path)
}

func (w *Watcher) stopAllTailers() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for path, cancel := range w.tailers {
		cancel()
		delete(w.tailers, path)
		log.Printf("stopped tailer for %s\n", path)
	}
}

func isLogFile(path string) bool {
	name := filepath.Base(path)
	return strings.HasSuffix(name, ".log") && !strings.HasPrefix(name, ".")
}
