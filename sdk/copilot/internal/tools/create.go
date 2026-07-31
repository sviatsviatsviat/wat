package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical create/write tool names accepted by [Input.AsCreate].
const (
	ToolCreate = "create"
	ToolWrite  = "Write"
)

// CreateInput is the input schema for the create tool.
type CreateInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// AsCreate returns the create tool input when this payload is for create or Write.
func (in Input) AsCreate() (CreateInput, bool) {
	return hookkit.AsFold[CreateInput](in.Input, ToolCreate, ToolWrite)
}
