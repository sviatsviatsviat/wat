package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical grep tool names accepted by [Input.AsGrep].
const (
	ToolGrep       = "grep"
	ToolGrepClaude = "Grep"
	ToolRG         = "rg"
)

// GrepInput is the input schema for the grep tool.
type GrepInput struct {
	Pattern         string `json:"pattern"`
	Path            string `json:"path,omitempty"`
	Glob            string `json:"glob,omitempty"`
	OutputMode      string `json:"output_mode,omitempty"`
	CaseInsensitive bool   `json:"-i,omitempty"`
	Multiline       bool   `json:"multiline,omitempty"`
}

// AsGrep returns the grep tool input when this payload is for grep, rg, or Grep.
func (in Input) AsGrep() (GrepInput, bool) {
	return hookkit.AsFold[GrepInput](in.Input, ToolGrep, ToolGrepClaude, ToolRG)
}
