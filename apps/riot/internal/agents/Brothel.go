package agents

import (
	"fmt"
	"sort"
	"sync"

	T "github.com/verikono/glocalize/types/go"
)

// AgenticBrothel is a concurrency-safe registry of named Prosti agents.
type AgenticBrothel struct {
	mu     sync.RWMutex
	agents map[AgentIdentifier]*Prosti
}

// NewAgenticBrothel returns an empty registry.
func NewAgenticBrothel(options ...T.OptionerFn) (AgentManager, error) {
	inst := &AgenticBrothel{
		agents: make(map[AgentIdentifier]*Prosti),
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(inst); err != nil {
			return nil, err
		}
	}
	return inst, nil
}

// NewAgent builds an agent and registers it under id.
func (a *AgenticBrothel) NewAgent(id AgentIdentifier, options ...OptionerFn) (*Prosti, error) {
	if id == "" {
		return nil, fmt.Errorf("agents: an agent identifier is required")
	}
	agent, err := NewProsti(options...)
	if err != nil {
		return nil, err
	}
	if err := a.RegisterAgent(id, agent); err != nil {
		return nil, err
	}
	return agent, nil
}

// RegisterAgent adds an already-built agent under id.
func (a *AgenticBrothel) RegisterAgent(id AgentIdentifier, agent *Prosti) error {
	if id == "" {
		return fmt.Errorf("agents: an agent identifier is required")
	}
	if agent == nil {
		return fmt.Errorf("agents: agent %q is nil", id)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, exists := a.agents[id]; exists {
		return fmt.Errorf("agents: agent %q is already registered", id)
	}
	a.agents[id] = agent
	return nil
}

// GetAgent returns the agent registered under id.
func (a *AgenticBrothel) GetAgent(id AgentIdentifier) (*Prosti, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	agent, ok := a.agents[id]
	if !ok {
		return nil, fmt.Errorf("agents: agent %q is not registered", id)
	}
	return agent, nil
}

// ListAgents returns the registered identifiers in sorted order.
func (a *AgenticBrothel) ListAgents() []AgentIdentifier {
	a.mu.RLock()
	defer a.mu.RUnlock()
	ids := make([]AgentIdentifier, 0, len(a.agents))
	for id := range a.agents {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// RemoveAgent unregisters the agent under id.
func (a *AgenticBrothel) RemoveAgent(id AgentIdentifier) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.agents[id]; !ok {
		return fmt.Errorf("agents: agent %q is not registered", id)
	}
	delete(a.agents, id)
	return nil
}
