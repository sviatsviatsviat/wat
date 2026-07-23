package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolGrep is the grep tool.
const ToolGrep = "Grep"

// GrepInput is the input schema for the Grep tool.
type GrepInput struct {
	Pattern         string `json:"pattern"`
	Path            string `json:"path,omitempty"`
	Glob            string `json:"glob,omitempty"`
	OutputMode      string `json:"output_mode,omitempty"`
	CaseInsensitive bool   `json:"-i,omitempty"`
	Multiline       bool   `json:"multiline,omitempty"`
}

// AsGrep returns the Grep tool input when this payload is for Grep.
func (in Input) AsGrep() (GrepInput, bool) {
	return hookkit.As[GrepInput](in.Input, ToolGrep)
}
