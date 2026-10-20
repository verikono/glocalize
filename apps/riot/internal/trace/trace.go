package trace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/verikono/glocalize/packages/go/tracer"
)

// DefaultFormat is the render format used when the caller supplies none.
const DefaultFormat = tracer.DefaultFormat

// logFileName is the file `riot serve` redirects the daemon's output into.
const logFileName = "riot.log"

// Options configures one trace session.
type Options struct {
	Path   string // log file to follow
	Format string // render format, e.g. "%l | %t | %m"
	Needle string // initial filter needle
	Level  string // initial level restriction
	Tail   int    // lines to replay before following
}

// DefaultLogPath returns the log file the riot daemon writes to. It mirrors the
// path `riot serve` opens for the detached child.
func DefaultLogPath() string {
	return filepath.Join(os.TempDir(), logFileName)
}

// Run opens the tracer on opts.Path and blocks until the user quits.
func Run(ctx context.Context, opts Options) error {
	source, err := newSource(opts)
	if err != nil {
		return err
	}

	options := []tracer.OptionerFn{
		tracer.WithApp("riot"),
		tracer.WithSource(source),
	}
	if opts.Format != "" {
		options = append(options, tracer.WithFormat(opts.Format))
	}
	if opts.Needle != "" {
		options = append(options, tracer.WithNeedle(opts.Needle))
	}
	if opts.Level != "" {
		options = append(options, tracer.WithLevel(tracer.LogLevel(opts.Level)))
	}

	viewer, err := tracer.NewTracer(options...)
	if err != nil {
		return fmt.Errorf("riot trace: %w", err)
	}
	return viewer.Run(ctx)
}

// newSource builds the log source, replaying the tail of the file when asked.
func newSource(opts Options) (tracer.TracerLogSource, error) {
	if opts.Tail <= 0 {
		return tracer.NewFileSource(opts.Path), nil
	}
	offset, err := tailOffset(opts.Path, opts.Tail)
	if err != nil {
		return nil, err
	}
	return tracer.NewFileSourceFrom(opts.Path, offset), nil
}

// tailOffset returns the byte offset at which the last n lines of path begin.
func tailOffset(path string, n int) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("riot trace: read %s: %w", path, err)
	}

	offset := len(data)
	if offset > 0 && data[offset-1] == '\n' {
		offset-- // the trailing newline terminates the last line; it does not start one
	}
	for i := 0; i < n; i++ {
		index := lastIndexByte(data[:offset], '\n')
		if index < 0 {
			return 0, nil
		}
		offset = index
	}
	if offset > 0 {
		offset++ // start after the newline that precedes the tail
	}
	return int64(offset), nil
}

// lastIndexByte returns the index of the last b in data, or -1.
func lastIndexByte(data []byte, b byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == b {
			return i
		}
	}
	return -1
}
