package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWebFetch is the normalized name for web fetch tools.
const ToolWebFetch = hookkit.ToolWebFetch

// WebFetchInput is the canonical web_fetch tool input.
type WebFetchInput struct {
	URL    string `json:"url"`
	Prompt string `json:"prompt,omitempty"`
}

// AsWebFetch returns the web_fetch tool input when this payload is for web_fetch.
func (in Input) AsWebFetch() (WebFetchInput, bool) {
	return as[WebFetchInput](in, ToolWebFetch)
}
