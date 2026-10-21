package agents

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	T "github.com/verikono/glocalize/types/go"
)

type (
	// OptionerFn is the shared functional-option signature used across the monorepo. Options receive the value they configure as an argument.
	OptionerFn = T.OptionerFn
)

// AgentIdentifier names an agent inside an AgenticBrothel.
type AgentIdentifier = string

// ExposedProsti is the read-only view of a Prosti.
type ExposedProsti interface {
	GetContext() context.Context
	GetModel() *ark.ChatModel
	GetAgent() *adk.ChatModelAgent
	GetRunner() *adk.Runner
	GetName() string
	Active() bool
	Alive() bool
}

// AgentProvider looks agents up by identifier.
type AgentProvider interface {
	GetAgent(id AgentIdentifier) (*Prosti, error)
	ListAgents() []AgentIdentifier
}

// AgentManager owns the lifecycle of the agents it provides.
type AgentManager interface {
	AgentProvider
	NewAgent(id AgentIdentifier, options ...OptionerFn) (*Prosti, error)
	RegisterAgent(id AgentIdentifier, agent *Prosti) error
	RemoveAgent(id AgentIdentifier) error
}

type AgentKind string

const (
	SimpleAgent      AgentKind = "simple"
	SupervisingAgent AgentKind = "supervisor"
)
