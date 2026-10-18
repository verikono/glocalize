package cmd

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/glo/internal"
)

func getParseCmd() *cli.Command {

	return &cli.Command{
		Name:  "parse",
		Usage: "merge prisma.schema files into .artifacts/data/schemas/launchboard.schema, then run phi generate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "out",
				Usage: "target path for the result",
			},
			&cli.StringFlag{
				Name:  "out-format",
				Usage: "format to output. If not provided then the out-format will be resolved via the out file's file suffix",
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			sources := c.Args().Slice()
			if len(sources) == 0 {
				return fmt.Errorf("parse requires at least one source file")
			}
			target := c.String("out")
			opts := []OptionerFn{}
			if c.String("out-format") != "" {
				opts = append(opts, internal.WithTargetFormat(c.String("out-format")))
			}
			result, err := internal.ParseSchemas(sources, target, opts...)
			if err != nil {
				fmt.Printf("\nerror: %s\n", err)
				return err
			}
			fmt.Printf("\n%s\n%s\n", strings.Join(sources, ", "), *result)
			return nil
		},
	}
}
