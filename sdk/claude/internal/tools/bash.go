package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolBash is the shell execution tool.
const ToolBash = "Bash"

// BashInput is the input schema for the Bash tool.
type BashInput struct {
	Command         string `json:"command"`
	Description     string `json:"description,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	RunInBackground bool   `json:"run_in_background,omitempty"`
}

// AsBash returns the Bash tool input when this payload is for Bash.
func (in Input) AsBash() (BashInput, bool) {
	return hookkit.As[BashInput](in.Input, ToolBash)
}
