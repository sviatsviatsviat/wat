package agenthooks

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
