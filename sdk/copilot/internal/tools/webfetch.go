package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWebFetch is the web fetch tool.
const ToolWebFetch = "web_fetch"

// WebFetchInput is the input schema for the web_fetch tool.
type WebFetchInput struct {
	URL string `json:"url"`
}

// AsWebFetch returns the web_fetch tool input when this payload is for web_fetch.
func (in Input) AsWebFetch() (WebFetchInput, bool) {
	return hookkit.As[WebFetchInput](in.Input, ToolWebFetch)
}
