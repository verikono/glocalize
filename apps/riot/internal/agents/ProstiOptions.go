package agents

import "github.com/cloudwego/eino-ext/components/model/ark"

func WithModelConfig(c ark.ChatModelAgentConfig) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(ark.ChatModelConfig); ok {
				// iterate through the keys in `c` and apply them to the inst[c.key] = c.value
				inst
			}
		}
		return nil
	}
}

func WithAgentConfig(c ark.ChatModelConfig) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(ark.ChatModelConfig); ok {
				// iterate through the keys in `c` and apply them to the inst[c.key] = c.value
				inst
			}
		}
		return nil
	}
}
