package tools

import "encoding/json"

// Canonical view/read tool names accepted by [Input.AsView].
const (
	ToolView = "view"
	ToolRead = "Read"
)

// ViewInput is the input schema for the view tool.
type ViewInput struct {
	Path   string `json:"path"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// AsView returns the view tool input when this payload is for view or Read.
// Claude-format Read payloads may use file_path instead of path.
func (in Input) AsView() (ViewInput, bool) {
	if !toolNameFold(in.Name(), ToolView, ToolRead) {
		return ViewInput{}, false
	}
	if !in.HasRaw() {
		return ViewInput{}, true
	}
	var wire struct {
		Path     string `json:"path"`
		FilePath string `json:"file_path"`
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
	}
	if json.Unmarshal(in.Raw(), &wire) != nil {
		return ViewInput{}, false
	}
	return ViewInput{
		Path:   firstNonEmpty(wire.Path, wire.FilePath),
		Offset: wire.Offset,
		Limit:  wire.Limit,
	}, true
}
