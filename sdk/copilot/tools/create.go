package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolCreate is the file creation tool.
const ToolCreate = "create"

// CreateInput is the input schema for the create tool.
type CreateInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// AsCreate returns the create tool input when this payload is for create.
func (in Input) AsCreate() (CreateInput, bool) {
	return hookkit.As[CreateInput](in.Input, ToolCreate)
}
