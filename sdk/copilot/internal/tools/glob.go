package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical glob tool names accepted by [Input.AsGlob].
const (
	ToolGlob       = "glob"
	ToolGlobClaude = "Glob"
)

// GlobInput is the input schema for the glob tool.
type GlobInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// AsGlob returns the glob tool input when this payload is for glob or Glob.
func (in Input) AsGlob() (GlobInput, bool) {
	return hookkit.AsFold[GlobInput](in.Input, ToolGlob, ToolGlobClaude)
}
