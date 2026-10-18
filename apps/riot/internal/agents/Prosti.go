package agents

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	_ "log"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	_ "github.com/cloudwego/eino/components/tool"
	_ "github.com/cloudwego/eino/compose"
)

func NewProsti(options ...OptionerFn) (*Prosti, error) {

	prosti := Prosti{
		ctx: context.Background(),
	}

	var err error

	err = prosti.newModel(options...)
	if err != nil {
		return nil, err
	}

	err = prosti.newAgent(options...)
	if err != nil {
		return nil, err
	}

	return &prosti, nil
}

type Prosti struct {
	ctx    context.Context
	agent  *adk.ChatModelAgent
	model  *ark.ChatModel
	runner *adk.Runner
	ExposedProsti
}

func (p *Prosti) GetModel() *ark.ChatModel {
	return p.model
}

func (p *Prosti) GetContext() context.Context {
	return p.ctx
}

func (p *Prosti) GetAgent() *adk.ChatModelAgent {
	return p.agent
}

func (p *Prosti) Start(options ...OptionerFn) error {
	if p.model != nil && p.agent != nil {
		config, err := getDefaultRunnerConfig(p.agent)
		if err != nil {
			return err
		}
		for _, option := range options {
			err := option(config)
			if err != nil {
				return err
			}
		}
		p.runner = adk.NewRunner(p.ctx, config)
		return nil
	}
	return fmt.Errorf("failed starting agent")
}

func (p *Prosti) newAgent(options ...OptionerFn) error {
	if p.model != nil {
		config, err := getDefaultAgentConfig(p.model)
		if err != nil {
			return err
		}
		agent, err := adk.NewChatModelAgent(p.ctx, config)
		if err != nil {
			return err
		}
		p.agent = agent
		return nil
	}
	return fmt.Errorf("failed instantiating new agent")
}

func (p *Prosti) newModel(options ...OptionerFn) error {
	if p.ctx != nil {
		config, err := getDefaultModelConfig()
		if err != nil {
			return err
		}
		model, err := ark.NewChatModel(p.ctx, &config)
		if err != nil {
			log.Fatal(err)
		}

		p.model = model
		return nil
	}
	return fmt.Errorf("failed instantiating new model reference")
}
