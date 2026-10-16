package cmd

import (
	"context"

	cli "github.com/urfave/cli/v3"
)

func getParseCmd() *cli.Command {

	return &cli.Command{
		Name:  "parse",
		Usage: "merge prisma.schema files into .artifacts/data/schemas/launchboard.schema, then run phi generate",
		Args: []cli.ArgUsage{
			{
				Name:  "source",
				Usage: "source file to parse from",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "to",
				Usage: "target format",
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {

			source := c.Args().Get(0)
			target := c.String("to")

			return nil
		},
	}
}
