package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolRead is the file read tool.
const ToolRead = "Read"

// ReadInput is the input schema for the Read tool.
type ReadInput struct {
	Path string `json:"path"`
}

// AsRead returns the Read tool input when this payload is for Read.
func (in Input) AsRead() (ReadInput, bool) {
	return hookkit.As[ReadInput](in.Input, ToolRead)
}
