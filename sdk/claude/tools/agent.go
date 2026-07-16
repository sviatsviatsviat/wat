package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolAgent is the agent tool.
const ToolAgent = "Agent"

// AgentInput is the input schema for the Agent tool.
type AgentInput struct {
	Prompt       string `json:"prompt"`
	Description  string `json:"description,omitempty"`
	SubagentType string `json:"subagent_type,omitempty"`
	Model        string `json:"model,omitempty"`
}

// AsAgent returns the Agent tool input when this payload is for Agent.
func (in Input) AsAgent() (AgentInput, bool) {
	return hookkit.As[AgentInput](in.Input, ToolAgent)
}
