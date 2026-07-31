package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical edit tool names accepted by [Input.AsEdit].
const (
	ToolEdit             = "edit"
	ToolEditClaude       = "Edit"
	ToolStrReplaceEditor = "str_replace_editor"
	ToolApplyPatch       = "apply_patch"
)

// EditInput is the input schema for the edit tool.
type EditInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// AsEdit returns the edit tool input when this payload is for edit, Edit,
// str_replace_editor, or apply_patch.
func (in Input) AsEdit() (EditInput, bool) {
	return hookkit.AsFold[EditInput](in.Input, ToolEdit, ToolEditClaude, ToolStrReplaceEditor, ToolApplyPatch)
}
