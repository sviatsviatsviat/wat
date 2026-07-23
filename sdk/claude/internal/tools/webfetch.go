package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWebFetch is the web fetch tool.
const ToolWebFetch = "WebFetch"

// WebFetchInput is the input schema for the WebFetch tool.
type WebFetchInput struct {
	URL    string `json:"url"`
	Prompt string `json:"prompt,omitempty"`
}

// AsWebFetch returns the WebFetch tool input when this payload is for WebFetch.
func (in Input) AsWebFetch() (WebFetchInput, bool) {
	return hookkit.As[WebFetchInput](in.Input, ToolWebFetch)
}
