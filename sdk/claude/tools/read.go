package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolRead is the file read tool.
const ToolRead = "Read"

// ReadInput is the input schema for the Read tool.
type ReadInput struct {
	FilePath string `json:"file_path"`
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// AsRead returns the Read tool input when this payload is for Read.
func (in Input) AsRead() (ReadInput, bool) {
	return hookkit.As[ReadInput](in.Input, ToolRead)
}
