package cmd

import (
	"context"
	"fmt"

	cli "github.com/urfave/cli/v3"
)

var Version = "dev"

func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "parse",
		Usage: "glocalize development tooling",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "show launch version",
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			if c.Bool("version") {
				showVersion()
				return nil
			}
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			getParseCmd(),
		},
	}
}

// showVersion prints the launch name and version to the CLI.
func showVersion() {
	fmt.Printf("glock\nversion: %s\n", Version)
}
