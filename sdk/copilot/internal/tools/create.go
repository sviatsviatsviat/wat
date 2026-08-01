package tools

import "encoding/json"

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
// Claude-format Write payloads may use file_path instead of path.
func (in Input) AsCreate() (CreateInput, bool) {
	if !toolNameFold(in.Name(), ToolCreate, ToolWrite) {
		return CreateInput{}, false
	}
	if !in.HasRaw() {
		return CreateInput{}, true
	}
	var wire struct {
		Path     string `json:"path"`
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}
	if json.Unmarshal(in.Raw(), &wire) != nil {
		return CreateInput{}, false
	}
	return CreateInput{
		Path:    firstNonEmpty(wire.Path, wire.FilePath),
		Content: wire.Content,
	}, true
}
