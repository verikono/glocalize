package agents

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// NewProsti builds a Prosti from the defaults, applying options in order. The model, the agent and the runner are all constructed here, so a Prosti returned without error is ready to run.
func NewProsti(options ...OptionerFn) (*Prosti, error) {
	cfg := defaultProstiConfig()
	if err := applyOptions(&cfg, options...); err != nil {
		return nil, fmt.Errorf("agents: applying options: %w", err)
	}
	if cfg.ctx == nil {
		cfg.ctx = context.Background()
	}
	if err := validateModelConfig(cfg.model); err != nil {
		return nil, err
	}

	model, err := ark.NewChatModel(cfg.ctx, &cfg.model)
	if err != nil {
		return nil, fmt.Errorf("agents: creating ark chat model: %w", err)
	}

	agent, err := adk.NewChatModelAgent(cfg.ctx, &adk.ChatModelAgentConfig{
		Name:        cfg.name,
		Description: cfg.description,
		Instruction: cfg.instruction,
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{Tools: toolsOf(cfg)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("agents: creating chat model agent: %w", err)
	}

	runner := adk.NewRunner(cfg.ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: cfg.streaming,
	})

	inst := Prosti{
		kind:   SimpleAgent,
		ctx:    cfg.ctx,
		name:   cfg.name,
		agent:  agent,
		model:  model,
		runner: runner,
	}

	for _, option := range options {
		err := option(inst)
		if err != nil {
			return nil, err
		}
	}

	return &inst, nil

}

// Prosti is a single riot agent: an ark chat model wrapped in an eino ChatModelAgent and driven by a Runner.
type Prosti struct {
	ctx    context.Context
	name   string
	kind   AgentKind
	agent  *adk.ChatModelAgent
	model  *ark.ChatModel
	runner *adk.Runner
	class  AgentKind
}

// GetContext returns the context the agent was built with.
func (p *Prosti) GetContext() context.Context {
	return p.ctx
}

// GetName returns the agent's configured name.
func (p *Prosti) GetName() string {
	return p.name
}

// GetModel returns the underlying ark chat model.
func (p *Prosti) GetModel() *ark.ChatModel {
	return p.model
}

func (p *Prosti) GetKind() AgentKind {
	return p.kind
}

// GetAgent returns the underlying eino agent.
func (p *Prosti) GetAgent() *adk.ChatModelAgent {
	return p.agent
}

// GetRunner returns the runner that executes the agent.
func (p *Prosti) GetRunner() *adk.Runner {
	return p.runner
}

// Active reports whether the agent is fully constructed and runnable.
func (p *Prosti) Active() bool {
	return p.ctx != nil && p.model != nil && p.agent != nil && p.runner != nil
}

// Returns true when the agent is capable of holding a conversation.
func (p *Prosti) Alive() bool {
	// work a way out to determine if the model will respond as expected
	return false
}

// Run executes the agent over the given messages and returns the event stream.
func (p *Prosti) Run(ctx context.Context, messages []adk.Message, options ...adk.AgentRunOption) (*adk.AsyncIterator[*adk.AgentEvent], error) {
	if !p.Active() {
		return nil, fmt.Errorf("agents: %q is not active", p.name)
	}
	if ctx == nil {
		ctx = p.ctx
	}
	return p.runner.Run(ctx, messages, options...), nil
}

// Query executes the agent with a single user message.
func (p *Prosti) Query(ctx context.Context, query string, options ...adk.AgentRunOption) (*adk.AsyncIterator[*adk.AgentEvent], error) {
	if !p.Active() {
		return nil, fmt.Errorf("agents: %q is not active", p.name)
	}
	if ctx == nil {
		ctx = p.ctx
	}
	return p.runner.Query(ctx, query, options...), nil
}

// Collect drains an event stream and returns the concatenated assistant text.
func Collect(events *adk.AsyncIterator[*adk.AgentEvent]) (string, error) {
	if events == nil {
		return "", fmt.Errorf("agents: nil event stream")
	}
	var out string
	for {
		event, ok := events.Next()
		if !ok {
			return out, nil
		}
		if event.Err != nil {
			return out, event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		variant := event.Output.MessageOutput
		if variant.Role != schema.Assistant {
			continue
		}
		message, err := variant.GetMessage()
		if err != nil {
			return out, err
		}
		if message != nil {
			out += message.Content
		}
	}
}
