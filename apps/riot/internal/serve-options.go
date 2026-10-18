package internal

import (
	"os"

	T "github.com/verikono/glocalize/types/go"
)

// WithGqlPath overrides the GraphQL mount point.
func WithGqlPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.GqlURL = value
			}
		}
		return nil
	}
}

// WithMcpPath overrides the MCP mount point.
func WithMcpPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.McpURL = value
			}
		}
		return nil
	}
}

// WithPort overrides the serve port.
func WithPort(port int) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.Port = port
			}
		}
		return nil
	}
}

// WithPlaygroundPath mounts the GraphQL playground at the given path when
// non-empty.
func WithPlaygroundPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.Playground = value
			}
		}
		return nil
	}
}

func (c *RiotConfig) applyEnvVars() error {

	if os.Getenv("RIOT_GRAPHL_PLAYGROUND") != "" {
		c.Playground = os.Getenv("RIOT_GRAPHL_PLAYGROUND")
	}

	return nil
}
