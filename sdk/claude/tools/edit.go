package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolEdit is the file edit tool.
const ToolEdit = "Edit"

// EditInput is the input schema for the Edit tool.
type EditInput struct {
	FilePath   string `json:"file_path"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all,omitempty"`
}

// AsEdit returns the Edit tool input when this payload is for Edit.
func (in Input) AsEdit() (EditInput, bool) {
	return hookkit.As[EditInput](in.Input, ToolEdit)
}
