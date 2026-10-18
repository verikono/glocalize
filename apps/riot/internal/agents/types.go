package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	T "github.com/verikono/glocalize/types/go"
)

type (
	OptionerFn = T.OptionerFn
)

type ExposedProsti interface {
	GetModel() adk.ChatModelAgent
	GetContext() context.Context
	Active() bool
}

type AgentIdentifier = string

type AgentProvider interface {
	GetAgent()
	GetDeepAgent()
}

type AgentManager interface {
	AgentProvider
}
