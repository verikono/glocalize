package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	cli "github.com/urfave/cli/v3"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

// getStatusCmd returns the `riot status` command, reporting the background
// server's URLs, database and uptime.
func getStatusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "report the background server's mcp/gql urls, database and uptime",
		Action: func(_ context.Context, _ *cli.Command) error {
			return writeStatus(os.Stdout)
		},
	}
}

func writeStatus(w io.Writer) error {
	base := fmt.Sprintf("http://localhost:%d", utils.ResolvePort())

	fmt.Fprintf(w, "%-12s : %s\n", "name", "riot")
	fmt.Fprintf(w, "%-12s : %s\n", "version", Version)
	fmt.Fprintf(w, "%-12s : %s\n", "gql", base+"/graph")
	fmt.Fprintf(w, "%-12s : %s\n", "mcp", base+"/mcp")
	fmt.Fprintf(w, "%-12s : %s\n", "database", utils.ResolveDataURL())
	fmt.Fprintf(w, "%-12s : %s\n", "uptime", daemonUptime())
	return nil
}

func daemonUptime() string {
	pid, err := RunningPid()
	if err != nil {
		return "not running"
	}
	up, err := procUptime(pid)
	if err != nil {
		return "not running"
	}
	return formatUptime(up)
}

func formatUptime(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d.Truncate(time.Second) / time.Second)
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60
	seconds := total % 60

	switch {
	case days > 0:
		return fmt.Sprintf("up %dd%02dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("up %dh%02dm", hours, minutes)
	case minutes > 0:
		return fmt.Sprintf("up %dm%02ds", minutes, seconds)
	default:
		return fmt.Sprintf("up %ds", seconds)
	}
}
