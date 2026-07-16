package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolWebSearch is the web search tool.
const ToolWebSearch = "WebSearch"

// WebSearchInput is the input schema for the WebSearch tool.
type WebSearchInput struct {
	Query          string   `json:"query"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	BlockedDomains []string `json:"blocked_domains,omitempty"`
}

// AsWebSearch returns the WebSearch tool input when this payload is for WebSearch.
func (in Input) AsWebSearch() (WebSearchInput, bool) {
	return hookkit.As[WebSearchInput](in.Input, ToolWebSearch)
}
