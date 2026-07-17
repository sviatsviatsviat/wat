package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolGlob is the normalized name for glob search tools.
const ToolGlob = hookkit.ToolGlob

// GlobInput is the canonical glob tool input.
type GlobInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// AsGlob returns the glob tool input when this payload is for a glob tool.
func (in Input) AsGlob() (GlobInput, bool) {
	return as[GlobInput](in, ToolGlob)
}
