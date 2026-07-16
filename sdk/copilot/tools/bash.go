package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical shell tool names accepted by [Input.AsBash].
const (
	ToolBash       = "bash"
	ToolPowerShell = "powershell"
	ToolShell      = "shell"
)

// BashInput is the input schema for the bash tool.
type BashInput struct {
	Command string `json:"command"`
}

// AsBash returns the bash tool input when this payload is for a shell tool
// (bash, powershell, or shell; case-insensitive).
func (in Input) AsBash() (BashInput, bool) {
	return hookkit.AsFold[BashInput](in.Input, ToolBash, ToolPowerShell, ToolShell)
}
