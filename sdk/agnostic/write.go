package agnostic

import "encoding/json"

// WriteInput is the canonical write/create tool input.
type WriteInput struct {
	// Path is the file path (from path or file_path).
	Path string `json:"path"`
	// Content is the file contents.
	Content string `json:"content"`
}

// AsWrite returns the write tool input when this payload is for a write tool.
func (in ToolInput) AsWrite() (WriteInput, bool) {
	if in.name != ToolWrite {
		return WriteInput{}, false
	}
	var wire struct {
		Path     string `json:"path"`
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}
	if len(in.raw) == 0 {
		return WriteInput{}, true
	}
	if json.Unmarshal(in.raw, &wire) == nil {
		path := wire.Path
		if path == "" {
			path = wire.FilePath
		}
		return WriteInput{Path: path, Content: wire.Content}, true
	}
	return WriteInput{}, false
}
