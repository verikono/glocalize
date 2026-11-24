package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/riot/internal/mcp"
)

// getToolsCmd returns the `riot tools` command, which lists the MCP tool
// surface to the console.
func getToolsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tools",
		Usage: "list the mcp tools riot serves",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "include each tool's description",
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			return writeTools(os.Stdout, c.Bool("list"))
		},
	}
}

// writeTools renders the tool surface. Names are padded to the width of the
// longest name so the separator column lines up on every row.
func writeTools(w io.Writer, withDescriptions bool) error {
	tools := mcp.Tools()

	width := 0
	for _, tool := range tools {
		if len(tool.Name) > width {
			width = len(tool.Name)
		}
	}

	for _, tool := range tools {
		if !withDescriptions {
			fmt.Fprintln(w, tool.Name)
			continue
		}
		fmt.Fprintf(w, "%-*s : %s\n", width, tool.Name, oneLine(tool.Description))
	}
	return nil
}

// oneLine collapses a description onto a single line so every entry stays on
// one row.
func oneLine(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
