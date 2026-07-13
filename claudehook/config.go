package claudehook

import "encoding/json"

// Settings is the Claude Code settings.json hooks section shape.
type Settings struct {
	// Hooks maps event names to matcher groups.
	Hooks map[string][]MatcherGroup `json:"hooks"`
	// DisableAllHooks disables every hook when true.
	DisableAllHooks bool `json:"disableAllHooks,omitempty"`
}

// MatcherGroup is a Claude Code hook matcher group.
type MatcherGroup struct {
	// Matcher is the tool or event matcher string.
	Matcher string `json:"matcher,omitempty"`
	// If is a Claude-only group-level permission rule when present.
	If json.RawMessage `json:"if,omitempty"`
	// Hooks holds native handler JSON. Use ParseHandler and MarshalHandler for typed access.
	Hooks []json.RawMessage `json:"hooks"`
}

// Handler is a Claude Code hook handler definition.
type Handler struct {
	// Type is the handler kind: command, http, mcp_tool, prompt, or agent.
	Type string `json:"type"`
	// Command is the shell command for command-type handlers.
	Command string `json:"command,omitempty"`
	// Args is the exec-form argument list.
	Args []string `json:"args,omitempty"`
	// Shell is the shell interpreter (bash, powershell).
	Shell string `json:"shell,omitempty"`
	// Async runs the handler asynchronously when true.
	Async bool `json:"async,omitempty"`
	// AsyncRewake wakes the agent after async completion when true.
	AsyncRewake bool `json:"asyncRewake,omitempty"`
	// URL is the HTTP endpoint for http-type handlers.
	URL string `json:"url,omitempty"`
	// Headers are HTTP headers for http-type handlers.
	Headers map[string]string `json:"headers,omitempty"`
	// AllowedEnvVars lists env vars forwarded to http handlers.
	AllowedEnvVars []string `json:"allowedEnvVars,omitempty"`
	// Server is the MCP server name for mcp_tool handlers.
	Server string `json:"server,omitempty"`
	// Tool is the MCP tool name for mcp_tool handlers.
	Tool string `json:"tool,omitempty"`
	// Input is the MCP tool input for mcp_tool handlers.
	Input map[string]any `json:"input,omitempty"`
	// Prompt is the prompt text for prompt and agent handlers.
	Prompt string `json:"prompt,omitempty"`
	// Model is the model name for agent handlers.
	Model string `json:"model,omitempty"`
	// Timeout is the handler timeout in seconds.
	Timeout int `json:"timeout,omitempty"`
	// StatusMessage is the status message shown while the handler runs.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Once runs the handler only once per session when true.
	Once bool `json:"once,omitempty"`
}

// ParseHandler decodes native handler JSON into a Handler.
func ParseHandler(raw json.RawMessage) (Handler, error) {
	var h Handler
	if len(raw) == 0 {
		return h, nil
	}
	err := json.Unmarshal(raw, &h)
	return h, err
}

// MarshalHandler encodes a Handler as native handler JSON.
func MarshalHandler(h Handler) (json.RawMessage, error) {
	return json.Marshal(h)
}

// Handlers encodes typed handlers as native handler JSON blobs.
func Handlers(h ...Handler) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, 0, len(h))
	for _, handler := range h {
		raw, err := MarshalHandler(handler)
		if err != nil {
			return nil, err
		}
		out = append(out, raw)
	}
	return out, nil
}

// TimeoutSeconds returns the configured handler timeout in seconds.
func (h Handler) TimeoutSeconds() int {
	return h.Timeout
}
