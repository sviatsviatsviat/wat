package tools

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

// Cursor builtin tool name constants.
const (
	// ToolShell is the shell execution tool.
	ToolShell = "Shell"
	// ToolRead is the file read tool.
	ToolRead = "Read"
	// ToolWrite is the file write tool.
	ToolWrite = "Write"
	// ToolGrep is the grep tool.
	ToolGrep = "Grep"
	// ToolTask is the agent task tool.
	ToolTask = "Task"
	// ToolDelete is the delete tool.
	ToolDelete = "Delete"
)

// ShellInput is the input schema for the Shell tool.
type ShellInput struct {
	Command string `json:"command"`
}

// ReadInput is the input schema for the Read tool.
type ReadInput struct {
	Path string `json:"path"`
}

// EditInput is the input schema for file edit tools.
type EditInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// ToolInputAs decodes raw tool input JSON into T.
func ToolInputAs[T any](raw json.RawMessage) (T, error) {
	return hookkit.ToolInputAs[T](raw)
}
