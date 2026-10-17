package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	cli "github.com/urfave/cli/v3"
)

// getDownCmd returns the `riot down` command, which stops every running riot
// process (SIGTERM, then SIGKILL for stragglers).
func getDownCmd() *cli.Command {
	return &cli.Command{
		Name:  "down",
		Usage: "stop every running riot process",
		Action: func(_ context.Context, _ *cli.Command) error {
			return stopAll()
		},
	}
}

func stopAll() error {
	pids, err := findRiotPids()
	if err != nil {
		return fmt.Errorf("riot down: %w", err)
	}
	if len(pids) == 0 {
		fmt.Println("riot: not running")
		return nil
	}

	for _, pid := range pids {
		if proc, err := os.FindProcess(pid); err == nil {
			_ = proc.Signal(syscall.SIGTERM)
		}
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		allGone := true
		for _, pid := range pids {
			if isProcessAlive(pid) {
				allGone = false
				break
			}
		}
		if allGone {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	for _, pid := range pids {
		if isProcessAlive(pid) {
			if proc, err := os.FindProcess(pid); err == nil {
				_ = proc.Signal(syscall.SIGKILL)
			}
		}
	}

	fmt.Println("riot: stopped")
	return nil
}
