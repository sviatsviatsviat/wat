package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolShell is the shell execution tool.
const ToolShell = "Shell"

// ShellInput is the input schema for the Shell tool.
type ShellInput struct {
	Command string `json:"command"`
}

// AsShell returns the Shell tool input when this payload is for Shell.
func (in Input) AsShell() (ShellInput, bool) {
	return hookkit.AsFold[ShellInput](in.Input, ToolShell)
}
