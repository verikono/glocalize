package tracer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// followPollInterval is how often a file source re-checks for appended bytes.
const followPollInterval = 100 * time.Millisecond

// fromEnd is the offset sentinel meaning "start at the current end of file".
const fromEnd int64 = -1

// FileSource follows a log file, emitting each complete line as it is appended.
// It behaves like `tail -f` from the current end of the file, so a tracer shows
// what the process writes next rather than replaying history.
type FileSource struct {
	path string
	from int64
}

// NewFileSource follows path from its current end.
func NewFileSource(path string) *FileSource {
	return &FileSource{path: path, from: fromEnd}
}

// NewFileSourceFrom follows path from byte offset from, so a tracer can replay
// the tail of an existing log before following it. An offset of 0 replays the
// whole file.
func NewFileSourceFrom(path string, from int64) *FileSource {
	return &FileSource{path: path, from: from}
}

// Path returns the file being followed.
func (s *FileSource) Path() string {
	return s.path
}

// Follow opens the file and streams lines until ctx is cancelled. The returned
// channel is closed by the goroutine that owns it.
func (s *FileSource) Follow(ctx context.Context) (<-chan string, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("follow %s: %w", s.path, err)
	}

	offset := s.from
	if offset < 0 {
		info, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("follow %s: %w", s.path, err)
		}
		offset = info.Size()
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("follow %s: seek %d: %w", s.path, offset, err)
	}

	lines := make(chan string, 256)
	go s.pump(ctx, file, lines)
	return lines, nil
}

// pump reads complete lines and closes lines when it stops. It owns file.
func (s *FileSource) pump(ctx context.Context, file *os.File, lines chan<- string) {
	defer close(lines)
	defer file.Close()

	reader := bufio.NewReader(file)
	var pending string

	for {
		chunk, err := reader.ReadString('\n')
		if chunk != "" {
			pending += chunk
			for {
				index := indexByte(pending, '\n')
				if index < 0 {
					break
				}
				line := pending[:index]
				pending = pending[index+1:]
				if !send(ctx, lines, line) {
					return
				}
			}
		}
		if err != nil && err != io.EOF {
			return
		}
		if err == io.EOF {
			select {
			case <-ctx.Done():
				return
			case <-time.After(followPollInterval):
			}
		}
	}
}

// send delivers line unless ctx is cancelled first.
func send(ctx context.Context, lines chan<- string, line string) bool {
	select {
	case <-ctx.Done():
		return false
	case lines <- line:
		return true
	}
}

// indexByte returns the index of b in s, or -1.
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
