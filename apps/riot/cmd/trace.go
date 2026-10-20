package cmd

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/riot/internal/trace"
)

// getTraceCmd returns the `riot trace` command, which opens the full-screen
// tracer on the log file the running daemon writes to.
func getTraceCmd() *cli.Command {
	return &cli.Command{
		Name:  "trace",
		Usage: "open the full-screen tracer on the riot daemon's log",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   `render format: %t time, %l level, %m message (e.g. "%l | %t | %m")`,
				Value:   trace.DefaultFormat,
			},
			&cli.StringFlag{
				Name:    "needle",
				Aliases: []string{"n"},
				Usage:   "seed the filter; a leading / selects a regular expression",
			},
			&cli.StringFlag{
				Name:  "level",
				Usage: "restrict to one level: trace, info, warn, error, fatal",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"p"},
				Usage:   "log file to follow (default: the daemon's log)",
			},
			&cli.IntFlag{
				Name:  "tail",
				Usage: "replay the last N lines before following",
				Value: 0,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			path := c.String("file")
			if path == "" {
				path = trace.DefaultLogPath()
			}
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("riot trace: %s: %w", path, err)
			}

			return trace.Run(ctx, trace.Options{
				Path:   path,
				Format: c.String("format"),
				Needle: c.String("needle"),
				Level:  c.String("level"),
				Tail:   c.Int("tail"),
			})
		},
	}
}
