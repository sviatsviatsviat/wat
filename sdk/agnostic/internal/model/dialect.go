package model

import (
	"encoding/json"
	"strings"
)

// Dialect identifies the coding agent that emitted a hook event.
type Dialect int

const (
	// Unknown is returned when the originating agent cannot be determined.
	Unknown Dialect = iota
	// Claude is the Claude Code agent dialect.
	Claude
	// Copilot is the GitHub Copilot CLI or cloud agent dialect.
	Copilot
	// Cursor is the Cursor agent dialect.
	Cursor
)

// String returns the lowercase dialect name ("claude", "copilot", "cursor", "unknown").
func (d Dialect) String() string {
	switch d {
	case Claude:
		return "claude"
	case Copilot:
		return "copilot"
	case Cursor:
		return "cursor"
	default:
		return "unknown"
	}
}

// ParseDialect parses a dialect name from a CLI flag or config value.
// It accepts "claude", "claude-code", "claudecode", "copilot", "github-copilot",
// "gh", and "cursor". Unknown or empty strings return Unknown.
func ParseDialect(s string) Dialect {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude", "claude-code", "claudecode":
		return Claude
	case "copilot", "github-copilot", "gh":
		return Copilot
	case "cursor":
		return Cursor
	default:
		return Unknown
	}
}

// Detect infers the originating agent from a hook payload and environment hints.
// Payload shape is checked before environment variables because Cursor exports
// CLAUDE_PROJECT_DIR for Claude Code compatibility. When the dialect is already
// known (for example via ParseDialect from wat run --agent), callers should
// skip Detect and use the explicit value instead.
func Detect(payload []byte, getenv func(string) string) Dialect {
	if getenv == nil {
		getenv = func(string) string { return "" }
	}
	var probe map[string]json.RawMessage
	_ = json.Unmarshal(payload, &probe)
	has := func(k string) bool { _, ok := probe[k]; return ok }

	switch {
	case has("cursor_version") || has("conversation_id"):
		return Cursor
	case has("sessionId"):
		return Copilot
	case has("hook_event_name") && has("timestamp"):
		return Copilot
	case has("session_id"):
		return Claude
	}

	switch {
	case getenv("CURSOR_VERSION") != "":
		return Cursor
	case getenv("CLAUDE_PROJECT_DIR") != "":
		return Claude
	case getenv("COPILOT_HOME") != "":
		return Copilot
	}
	return Unknown
}
