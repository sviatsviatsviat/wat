package agnostic

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Dialect identifies the coding agent that emitted a hook event.
type Dialect = model.Dialect

const (
	// Unknown is returned when the originating agent cannot be determined.
	Unknown = model.Unknown
	// Claude is the Claude Code agent dialect.
	Claude = model.Claude
	// Copilot is the GitHub Copilot CLI or cloud agent dialect.
	Copilot = model.Copilot
	// Cursor is the Cursor agent dialect.
	Cursor = model.Cursor
)

// ParseDialect parses a dialect name from a CLI flag or config value.
// It accepts "claude", "claude-code", "claudecode", "copilot", "github-copilot",
// "gh", and "cursor". Unknown or empty strings return Unknown.
func ParseDialect(s string) Dialect { return model.ParseDialect(s) }

// Detect infers the originating agent from a hook payload and environment hints.
// Payload shape is checked before environment variables because Cursor exports
// CLAUDE_PROJECT_DIR for Claude Code compatibility. When the dialect is already
// known (for example via ParseDialect from wat run --agent), callers should
// skip Detect and use the explicit value instead.
func Detect(payload []byte, getenv func(string) string) Dialect {
	return model.Detect(payload, getenv)
}
