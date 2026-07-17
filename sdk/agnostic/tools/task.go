package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolTask is the normalized name for subagent or task tools.
const ToolTask = hookkit.ToolTask

// TaskInput is the canonical task/agent tool input.
type TaskInput struct {
	Prompt       string `json:"prompt"`
	Description  string `json:"description,omitempty"`
	SubagentType string `json:"subagent_type,omitempty"`
	Model        string `json:"model,omitempty"`
}

// AsTask returns the task tool input when this payload is for a task/agent tool.
func (in Input) AsTask() (TaskInput, bool) {
	return as[TaskInput](in, ToolTask)
}
