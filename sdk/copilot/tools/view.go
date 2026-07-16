package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolView is the file read tool.
const ToolView = "view"

// ViewInput is the input schema for the view tool.
type ViewInput struct {
	Path string `json:"path"`
}

// AsView returns the view tool input when this payload is for view.
func (in Input) AsView() (ViewInput, bool) {
	return hookkit.As[ViewInput](in.Input, ToolView)
}
