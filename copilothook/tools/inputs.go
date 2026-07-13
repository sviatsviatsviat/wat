package tools

import "encoding/json"

// Copilot builtin tool name constants (camelCase CLI wire names).
const (
	// ToolBash is the shell execution tool.
	ToolBash = "bash"
	// ToolCreate is the file creation tool.
	ToolCreate = "create"
	// ToolView is the file read tool.
	ToolView = "view"
	// ToolEdit is the file edit tool.
	ToolEdit = "edit"
	// ToolWebFetch is the web fetch tool.
	ToolWebFetch = "web_fetch"
	// ToolTask is the agent task tool.
	ToolTask = "task"
)

// BashInput is the input schema for the bash tool.
type BashInput struct {
	Command string `json:"command"`
}

// CreateInput is the input schema for the create tool.
type CreateInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// ViewInput is the input schema for the view tool.
type ViewInput struct {
	Path string `json:"path"`
}

// EditInput is the input schema for the edit tool.
type EditInput struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// WebFetchInput is the input schema for the web_fetch tool.
type WebFetchInput struct {
	URL string `json:"url"`
}

// ToolInputAs decodes raw tool input JSON into T.
func ToolInputAs[T any](raw json.RawMessage) (T, error) {
	var v T
	if len(raw) == 0 {
		return v, nil
	}
	err := json.Unmarshal(raw, &v)
	return v, err
}
