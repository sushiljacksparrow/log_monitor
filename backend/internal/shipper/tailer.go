package shipper

import (
	"bufio"
	"context"
	"log"
	"os"
	"syscall"
	"time"
)

// Tailer reads new lines from a log file and publishes them to Kafka.
type Tailer struct {
	path      string
	file      *os.File
	scanner   *bufio.Scanner
	offset    int64
	inode     uint64
	position  *PositionTracker
	publisher *Publisher
}

func NewTailer(path string, pos *PositionTracker, pub *Publisher) (*Tailer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	inode := getInode(f)

	// Resume from saved position if the file hasn't been rotated (same inode).
	var offset int64
	if saved, ok := pos.Get(path); ok && saved.Inode == inode {
		offset = saved.Offset
	}

	if _, err := f.Seek(offset, 0); err != nil {
		f.Close()
		return nil, err
	}

	return &Tailer{
		path:      path,
		file:      f,
		scanner:   bufio.NewScanner(f),
		offset:    offset,
		inode:     inode,
		position:  pos,
		publisher: pub,
	}, nil
}

// Run continuously tails the file until the context is cancelled.
func (t *Tailer) Run(ctx context.Context) {
	log.Printf("tailing %s from offset=%d\n", t.path, t.offset)

	for {
		select {
		case <-ctx.Done():
			t.file.Close()
			return
		default:
		}

		if t.scanner.Scan() {
			line := t.scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			// Make a copy since scanner reuses the buffer.
			lineCopy := make([]byte, len(line))
			copy(lineCopy, line)

			if err := t.publisher.Publish(lineCopy); err != nil {
				log.Printf("failed to publish line from %s: %v\n", t.path, err)
				continue
			}

			t.offset += int64(len(line)) + 1 // +1 for newline
			t.position.Update(t.path, t.offset, t.inode)
			continue
		}

		// No more lines — check for rotation or truncation, then wait.
		if t.checkRotation() {
			continue
		}

		select {
		case <-ctx.Done():
			t.file.Close()
			return
		case <-time.After(250 * time.Millisecond):
			// Reset scanner to pick up new data appended since last read.
			t.scanner = bufio.NewScanner(t.file)
		}
	}
}

// checkRotation detects if the file was rotated (new inode) or truncated (size < offset).
// Returns true if the file was reopened and the tailer should continue reading.
func (t *Tailer) checkRotation() bool {
	info, err := os.Stat(t.path)
	if err != nil {
		return false
	}

	currentInode := getInodeFromInfo(info)

	// File was replaced (rotation): new inode.
	if currentInode != t.inode {
		log.Printf("detected rotation for %s (inode %d -> %d)\n", t.path, t.inode, currentInode)
		return t.reopen()
	}

	// File was truncated: size is smaller than our offset.
	if info.Size() < t.offset {
		log.Printf("detected truncation for %s (size=%d, offset=%d)\n", t.path, info.Size(), t.offset)
		return t.reopen()
	}

	return false
}

func (t *Tailer) reopen() bool {
	t.file.Close()

	f, err := os.Open(t.path)
	if err != nil {
		log.Printf("failed to reopen %s: %v\n", t.path, err)
		return false
	}

	t.file = f
	t.scanner = bufio.NewScanner(f)
	t.offset = 0
	t.inode = getInode(f)
	t.position.Update(t.path, 0, t.inode)
	log.Printf("reopened %s from offset=0\n", t.path)
	return true
}

func getInode(f *os.File) uint64 {
	info, err := f.Stat()
	if err != nil {
		return 0
	}
	return getInodeFromInfo(info)
}

func getInodeFromInfo(info os.FileInfo) uint64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}
	return stat.Ino
}
