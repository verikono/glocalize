package cmd

import (
	"context"
	"fmt"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/riot/internal"
	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

// getServeCmd returns the `riot serve` (alias `riot up`) command. Unflagged it
// starts the server in the background; --daemon runs it in-process for the
// detached child.
func getServeCmd() *cli.Command {
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"up"},
		Usage:   "start the riot MCP + GraphQL server (default: in the background)",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "daemon", Usage: "run the server in-process (internal)", Hidden: true},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   fmt.Sprintf("port to serve on (default %d)", utils.DefaultPort),
				Value:   utils.ResolvePort(),
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			if c.Bool("daemon") {
				server, err := internal.NewRiotServer(internal.WithPort(c.Int("port")))
				if err != nil {
					return err
				}
				return server.Start()
			}
			return StartRiot(c.Int("port"))
		},
	}
}
