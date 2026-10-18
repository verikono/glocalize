package cmd

import (
	"context"
	"fmt"
	"time"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

// getRestartCmd returns the `riot restart` command, which stops the running
// daemon and starts it again with the flags it was originally launched with.
func getRestartCmd() *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: "restart the service with the same flags",
		Action: func(_ context.Context, _ *cli.Command) error {
			flags, err := loadRuntimeFlags()
			if err != nil {
				return err
			}

			if err := stopAll(); err != nil {
				return err
			}
			if !waitPortFree(flags.Port, 5*time.Second) {
				return fmt.Errorf("riot: restart: port %d still in use", flags.Port)
			}

			if flags.Hot {
				devDir := flags.DevDir
				if devDir == "" {
					devDir, err = utils.ResolveHotDir()
					if err != nil {
						return err
					}
				}
				return startRiot(flags.Port, devDir)
			}
			return startRiot(flags.Port, "")
		},
	}
}
