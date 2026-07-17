package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolGrep is the normalized name for grep search tools.
const ToolGrep = hookkit.ToolGrep

// GrepInput is the canonical grep tool input.
type GrepInput struct {
	Pattern         string `json:"pattern"`
	Path            string `json:"path,omitempty"`
	Glob            string `json:"glob,omitempty"`
	OutputMode      string `json:"output_mode,omitempty"`
	CaseInsensitive bool   `json:"-i,omitempty"`
	Multiline       bool   `json:"multiline,omitempty"`
}

// AsGrep returns the grep tool input when this payload is for a grep tool.
func (in Input) AsGrep() (GrepInput, bool) {
	return as[GrepInput](in, ToolGrep)
}
