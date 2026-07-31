package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical web_search tool names accepted by [Input.AsWebSearch].
const (
	ToolWebSearch       = "web_search"
	ToolWebSearchClaude = "WebSearch"
)

// WebSearchInput is the input schema for the web_search tool.
type WebSearchInput struct {
	Query          string   `json:"query"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	BlockedDomains []string `json:"blocked_domains,omitempty"`
}

// AsWebSearch returns the web_search tool input when this payload is for web_search or WebSearch.
func (in Input) AsWebSearch() (WebSearchInput, bool) {
	return hookkit.AsFold[WebSearchInput](in.Input, ToolWebSearch, ToolWebSearchClaude)
}
