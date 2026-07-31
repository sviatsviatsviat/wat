package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical web_fetch tool names accepted by [Input.AsWebFetch].
const (
	ToolWebFetch       = "web_fetch"
	ToolWebFetchClaude = "WebFetch"
)

// WebFetchInput is the input schema for the web_fetch tool.
type WebFetchInput struct {
	URL string `json:"url"`
}

// AsWebFetch returns the web_fetch tool input when this payload is for web_fetch or WebFetch.
func (in Input) AsWebFetch() (WebFetchInput, bool) {
	return hookkit.AsFold[WebFetchInput](in.Input, ToolWebFetch, ToolWebFetchClaude)
}
