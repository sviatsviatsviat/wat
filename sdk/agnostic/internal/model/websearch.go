package model

// WebSearchInput is the canonical web_search tool input.
type WebSearchInput struct {
	Query          string   `json:"query"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	BlockedDomains []string `json:"blocked_domains,omitempty"`
}

// AsWebSearch returns the web_search tool input when this payload is for web_search.
func (in ToolInput) AsWebSearch() (WebSearchInput, bool) {
	return as[WebSearchInput](in, ToolWebSearch)
}
