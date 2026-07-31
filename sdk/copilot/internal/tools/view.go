package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical view/read tool names accepted by [Input.AsView].
const (
	ToolView = "view"
	ToolRead = "Read"
)

// ViewInput is the input schema for the view tool.
type ViewInput struct {
	Path string `json:"path"`
}

// AsView returns the view tool input when this payload is for view or Read.
func (in Input) AsView() (ViewInput, bool) {
	return hookkit.AsFold[ViewInput](in.Input, ToolView, ToolRead)
}
