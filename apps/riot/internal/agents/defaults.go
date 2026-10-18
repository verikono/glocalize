package agents

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
)

const (
	defaultAgentName        = "prosti"
	defaultAgentDescription = "riot agent"
	defaultAgentInstruction = "You are a riot agent. Answer the user's request directly and concisely."
	defaultModel            = "doubao-seed-1-8-251228"
	defaultMaxTokens        = 4096
)

// prostiConfig is the mutable configuration an OptionerFn writes into. It is built fresh for every Prosti, so options never leak between agents.
type prostiConfig struct {
	ctx         context.Context
	name        string
	description string
	instruction string
	model       ark.ChatModelConfig
	tools       []tool.BaseTool
	streaming   bool
}

// defaultProstiConfig returns a fresh configuration seeded from the environment. Every call builds a new value.
func defaultProstiConfig() prostiConfig {
	maxTokens := defaultMaxTokens
	return prostiConfig{
		ctx:         context.Background(),
		name:        defaultAgentName,
		description: defaultAgentDescription,
		instruction: defaultAgentInstruction,
		model: ark.ChatModelConfig{
			Model:     envOr("RIOT_DEFAULT_MODEL", defaultModel),
			APIKey:    strings.TrimSpace(os.Getenv("RIOT_MODEL_KEY")),
			BaseURL:   strings.TrimSpace(os.Getenv("RIOT_MODEL_BASEURL")),
			MaxTokens: &maxTokens,
		},
		streaming: true,
	}
}

// envOr returns the trimmed value of key, or fallback when it is unset.
func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// applyOptions applies every non-nil option to cfg in order.
func applyOptions(cfg *prostiConfig, options ...OptionerFn) error {
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(cfg); err != nil {
			return err
		}
	}
	return nil
}

// toolsOf returns the configured tools, never nil, so the agent config is
// always well formed.
func toolsOf(cfg prostiConfig) []tool.BaseTool {
	if cfg.tools == nil {
		return []tool.BaseTool{}
	}
	return cfg.tools
}

// validateModelConfig rejects a model configuration the ark client cannot use,
// so the failure names the missing setting instead of surfacing deep inside the
// provider.
func validateModelConfig(cfg ark.ChatModelConfig) error {
	if strings.TrimSpace(cfg.Model) == "" {
		return fmt.Errorf("agents: a model id is required (WithModel)")
	}
	hasAPIKey := strings.TrimSpace(cfg.APIKey) != ""
	hasKeyPair := strings.TrimSpace(cfg.AccessKey) != "" && strings.TrimSpace(cfg.SecretKey) != ""
	if !hasAPIKey && !hasKeyPair {
		return fmt.Errorf("agents: ark credentials are required: set RIOT_ARK_API_KEY or WithAPIKey")
	}
	return nil
}
