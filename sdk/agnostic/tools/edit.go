package tools

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// ToolEdit is the normalized name for file edit tools.
const ToolEdit = hookkit.ToolEdit

// EditInput is the canonical edit tool input.
type EditInput struct {
	// Path is the file path (from path or file_path).
	Path string `json:"path"`
	// Content is replacement content when provided by the agent dialect.
	Content string `json:"content,omitempty"`
	// OldString is the text to replace (Claude Edit).
	OldString string `json:"old_string,omitempty"`
	// NewString is the replacement text (Claude Edit).
	NewString string `json:"new_string,omitempty"`
}

// AsEdit returns the edit tool input when this payload is for an edit tool.
func (in Input) AsEdit() (EditInput, bool) {
	if in.name != ToolEdit {
		return EditInput{}, false
	}
	var wire struct {
		Path      string `json:"path"`
		FilePath  string `json:"file_path"`
		Content   string `json:"content"`
		OldString string `json:"old_string"`
		NewString string `json:"new_string"`
	}
	if len(in.raw) == 0 {
		return EditInput{}, true
	}
	if json.Unmarshal(in.raw, &wire) == nil {
		path := wire.Path
		if path == "" {
			path = wire.FilePath
		}
		return EditInput{
			Path:      path,
			Content:   wire.Content,
			OldString: wire.OldString,
			NewString: wire.NewString,
		}, true
	}
	return EditInput{}, false
}
