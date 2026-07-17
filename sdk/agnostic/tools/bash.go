package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolBash is the normalized name for shell execution tools.
const ToolBash = hookkit.ToolBash

// BashInput is the canonical shell tool input.
type BashInput struct {
	Command string `json:"command"`
}

// AsBash returns the bash tool input when this payload is for a shell tool.
func (in Input) AsBash() (BashInput, bool) {
	return as[BashInput](in, ToolBash)
}
