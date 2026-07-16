package agnostic

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Dialect identifies the coding agent that emitted a hook event.
type Dialect = model.Dialect

// Dialect constants.
const (
	Unknown = model.Unknown
	Claude  = model.Claude
	Copilot = model.Copilot
	Cursor  = model.Cursor
)

// Kind is a normalized hook event category shared by all supported agents.
type Kind = model.Kind

// Normalized event kind constants.
const (
	KindSessionStart      = model.KindSessionStart
	KindSessionEnd        = model.KindSessionEnd
	KindUserPrompt        = model.KindUserPrompt
	KindPreTool           = model.KindPreTool
	KindPostTool          = model.KindPostTool
	KindPostToolFailure   = model.KindPostToolFailure
	KindPermissionRequest = model.KindPermissionRequest
	KindSubagentStart     = model.KindSubagentStart
	KindSubagentStop      = model.KindSubagentStop
	KindStop              = model.KindStop
	KindPreCompact        = model.KindPreCompact
	KindNotification      = model.KindNotification
	KindAgentError        = model.KindAgentError
	KindOther             = model.KindOther
)

// Event is the unified, agent-independent view of a hook invocation.
type Event = model.Event

// Decision is the unified gate verdict for pre-events.
type Decision = model.Decision

// Decision constants.
const (
	DecisionUnset = model.DecisionUnset
	DecisionAllow = model.DecisionAllow
	DecisionAsk   = model.DecisionAsk
	DecisionDeny  = model.DecisionDeny
)

// Result is the wire hook response used by codecs.
type Result = model.Result

// PreToolResult is the portable hook response for PreTool events.
type PreToolResult = model.PreToolResult

// PostToolResult is the portable hook response for PostTool events.
type PostToolResult = model.PostToolResult

// PostToolFailureResult is the portable hook response for PostToolFailure events.
type PostToolFailureResult = model.PostToolFailureResult

// StopResult is the portable hook response for Stop and SubagentStop events.
type StopResult = model.StopResult

// SessionStartResult is the portable hook response for SessionStart events.
type SessionStartResult = model.SessionStartResult

// ParseDialect parses a dialect name from a CLI flag or config value.
func ParseDialect(s string) Dialect { return model.ParseDialect(s) }

// Detect infers the originating agent from a hook payload and environment hints.
func Detect(payload []byte, getenv func(string) string) Dialect {
	return model.Detect(payload, getenv)
}

// NormalizeToolName maps a native tool name onto the canonical vocabulary.
func NormalizeToolName(name string) (canonical string, mcp bool) {
	return model.NormalizeToolName(name)
}
