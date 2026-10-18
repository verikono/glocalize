package cmd

import (
	cli "github.com/urfave/cli/v3"
)

// Version is the riot version, injected at build time.
var Version = "dev"

// NewApp returns the root `riot` command with its subcommands.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "riot",
		Usage: "riot MCP + GraphQL service",
		Commands: []*cli.Command{
			getServeCmd(),
			getRestartCmd(),
			getStatusCmd(),
			getDownCmd(),
		},
	}
}
