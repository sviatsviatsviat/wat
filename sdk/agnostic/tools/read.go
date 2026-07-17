package tools

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// ToolRead is the normalized name for file read tools.
const ToolRead = hookkit.ToolRead

// ReadInput is the canonical read/view tool input.
type ReadInput struct {
	// Path is the file path (from path or file_path).
	Path string `json:"path"`
	// Offset is the optional start line.
	Offset int `json:"offset,omitempty"`
	// Limit is the optional line limit.
	Limit int `json:"limit,omitempty"`
}

// AsRead returns the read tool input when this payload is for a read tool.
func (in Input) AsRead() (ReadInput, bool) {
	if in.name != ToolRead {
		return ReadInput{}, false
	}
	var wire struct {
		Path     string `json:"path"`
		FilePath string `json:"file_path"`
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
	}
	if len(in.raw) == 0 {
		return ReadInput{}, true
	}
	if json.Unmarshal(in.raw, &wire) == nil {
		path := wire.Path
		if path == "" {
			path = wire.FilePath
		}
		return ReadInput{Path: path, Offset: wire.Offset, Limit: wire.Limit}, true
	}
	return ReadInput{}, false
}
