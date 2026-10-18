package agents

import (
	_ "fmt"
	_ "log"

	_ "github.com/cloudwego/eino/components/tool"
	_ "github.com/cloudwego/eino/compose"
	T "github.com/verikono/glocalize/types/go"
)

type AgenticBrothel struct {
	agents map[AgentIdentifier]Prosti
}

func NewAgenticBrothel(options ...T.OptionerFn) (AgentManager, error) {
	inst := AgenticBrothel{
		agents: make(map[AgentIdentifier]Prosti),
	}
	return &inst, nil
}

func (a *AgenticBrothel) GetAgent() {

}

func (a *AgenticBrothel) GetDeepAgent() {

}
