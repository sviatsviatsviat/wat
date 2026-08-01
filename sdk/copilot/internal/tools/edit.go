package tools

import (
	"encoding/json"
	"strings"
)

// Canonical edit tool names accepted by [Input.AsEdit].
const (
	ToolEdit             = "edit"
	ToolEditClaude       = "Edit"
	ToolStrReplaceEditor = "str_replace_editor"
	ToolApplyPatch       = "apply_patch"
)

// EditInput is the input schema for the edit tool.
type EditInput struct {
	Path      string `json:"path"`
	Content   string `json:"content,omitempty"`
	OldString string `json:"old_string,omitempty"`
	NewString string `json:"new_string,omitempty"`
	Patch     string `json:"patch,omitempty"`
}

// AsEdit returns the edit tool input when this payload is for edit, Edit,
// str_replace_editor, or apply_patch. Claude-format and text-editor field
// aliases are normalized onto this struct.
func (in Input) AsEdit() (EditInput, bool) {
	if !toolNameFold(in.Name(), ToolEdit, ToolEditClaude, ToolStrReplaceEditor, ToolApplyPatch) {
		return EditInput{}, false
	}
	if !in.HasRaw() {
		return EditInput{}, true
	}
	var wire struct {
		Path      string `json:"path"`
		FilePath  string `json:"file_path"`
		Content   string `json:"content"`
		FileText  string `json:"file_text"`
		OldString string `json:"old_string"`
		NewString string `json:"new_string"`
		OldStr    string `json:"old_str"`
		NewStr    string `json:"new_str"`
		Patch     string `json:"patch"`
	}
	if json.Unmarshal(in.Raw(), &wire) != nil {
		return EditInput{}, false
	}
	return EditInput{
		Path:      firstNonEmpty(wire.Path, wire.FilePath),
		Content:   firstNonEmpty(wire.Content, wire.FileText),
		OldString: firstNonEmpty(wire.OldString, wire.OldStr),
		NewString: firstNonEmpty(wire.NewString, wire.NewStr),
		Patch:     wire.Patch,
	}, true
}

func toolNameFold(got string, want ...string) bool {
	for _, name := range want {
		if strings.EqualFold(got, name) {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
