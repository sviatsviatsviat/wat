package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical task/agent tool names accepted by [Input.AsTask].
const (
	ToolTask        = "task"
	ToolTaskClaude  = "Task"
	ToolAgentClaude = "Agent"
)

// TaskInput is the input schema for the task/agent tool.
type TaskInput struct {
	Prompt       string `json:"prompt"`
	Description  string `json:"description,omitempty"`
	SubagentType string `json:"subagent_type,omitempty"`
	Model        string `json:"model,omitempty"`
}

// AsTask returns the task tool input when this payload is for task, Task, or Agent.
func (in Input) AsTask() (TaskInput, bool) {
	return hookkit.AsFold[TaskInput](in.Input, ToolTask, ToolTaskClaude, ToolAgentClaude)
}
