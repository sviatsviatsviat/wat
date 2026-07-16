package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolGlob is the glob tool.
const ToolGlob = "Glob"

// GlobInput is the input schema for the Glob tool.
type GlobInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// AsGlob returns the Glob tool input when this payload is for Glob.
func (in Input) AsGlob() (GlobInput, bool) {
	return hookkit.As[GlobInput](in.Input, ToolGlob)
}
