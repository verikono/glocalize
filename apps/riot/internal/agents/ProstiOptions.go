package agents

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
)

// WithContext overrides the context the agent is built with.
func WithContext(ctx context.Context) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.ctx = ctx
			}
		}
		return nil
	}
}

// WithName overrides the agent's name.
func WithName(name string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.name = name
			}
		}
		return nil
	}
}

// With Agent Kind
func WithKind(kind AgentKind) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(Prosti); ok {
				inst.kind = kind
			}
		}
		return nil
	}
}

// WithDescription overrides the agent's description.
func WithDescription(description string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.description = description
			}
		}
		return nil
	}
}

// WithInstruction overrides the agent's system prompt.
func WithInstruction(instruction string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.instruction = instruction
			}
		}
		return nil
	}
}

// WithModelConfig overrides the ark chat model configuration.
func WithModelConfig(config ark.ChatModelConfig) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.model = config
			}
		}
		return nil
	}
}

// WithModel overrides only the ark model id, keeping the other model settings.
func WithModel(model string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.model.Model = model
			}
		}
		return nil
	}
}

// WithAPIKey overrides the ark api key.
func WithAPIKey(apiKey string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.model.APIKey = apiKey
			}
		}
		return nil
	}
}

// WithTools replaces the agent's tool set.
func WithTools(tools ...tool.BaseTool) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.tools = tools
			}
		}
		return nil
	}
}

// WithStreaming toggles streaming on the runner.
func WithStreaming(enabled bool) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if cfg, ok := arg.(*prostiConfig); ok {
				cfg.streaming = enabled
			}
		}
		return nil
	}
}
