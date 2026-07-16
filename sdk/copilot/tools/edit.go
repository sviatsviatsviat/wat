package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolEdit is the file edit tool.
const ToolEdit = "edit"

// EditInput is the input schema for the edit tool.
type EditInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// AsEdit returns the edit tool input when this payload is for edit.
func (in Input) AsEdit() (EditInput, bool) {
	return hookkit.As[EditInput](in.Input, ToolEdit)
}
