package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWrite is the file write tool.
const ToolWrite = "Write"

// WriteInput is the input schema for the Write tool.
type WriteInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// AsWrite returns the Write tool input when this payload is for Write.
func (in Input) AsWrite() (WriteInput, bool) {
	return hookkit.As[WriteInput](in.Input, ToolWrite)
}
