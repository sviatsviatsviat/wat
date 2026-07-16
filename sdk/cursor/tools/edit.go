package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWrite is the file write tool.
const ToolWrite = "Write"

// EditInput is the input schema for file edit tools.
type EditInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// AsEdit returns the edit tool input when this payload is for Write (edit).
func (in Input) AsEdit() (EditInput, bool) {
	return hookkit.As[EditInput](in.Input, ToolWrite)
}
