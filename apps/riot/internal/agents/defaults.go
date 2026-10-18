package agents

import (
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var defChatModelConfig = ark.ChatModelConfig{
	Model:  "doubao-seed-1-8-251228",
	APIKey: "your_api_key",
}

var defAgentConfig = adk.ChatModelAgentConfig{
	Name:        "unnamed",
	Description: "riot agent",
	Instruction: "Instruct the user this agent is misconfigured",
	ToolsConfig: adk.ToolsConfig{
		ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				// Register your tools, e.g., webSearchTool
			},
		},
	},
	// Handlers: []adk.ChatModelAgentMiddleware{...}, // Register Middleware
}

var defRunnnerConfig = adk.RunnerConfig{
	EnableStreaming: true,
}

func getDefaultAgentConfig(model ark.ChatModel) (*adk.ChatModelAgentConfig, error) {
	config := DeepCopy[adk.ChatModelAgentConfig](defAgentConfig)
	config.Model = model
	return &config, nil
}

func getDefaultModelConfig() (ark.ChatModelConfig, error) {
	config := DeepCopy[ark.ChatModelConfig](defChatModelConfig)
	return config, nil
}

func getDefaultRunnerConfig(a *adk.ChatModelAgent) (adk.RunnerConfig, error) {
	config := DeepCopy[ark.RunnerConfig](defRunnnerConfig)
	config.Agent = a
	return config, nil
}
