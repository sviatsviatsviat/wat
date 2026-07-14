package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

type (
	// Dialect identifies the coding agent that emitted a hook event.
	Dialect = model.Dialect
	// Kind is a normalized hook event category shared by all supported agents.
	Kind = model.Kind
	// Event is the unified, agent-independent view of a hook invocation.
	Event = model.Event
	// ToolCall describes the tool invocation a pre/post tool event refers to.
	ToolCall = model.ToolCall
	// ToolResult describes the outcome of a completed or failed tool call.
	ToolResult = model.ToolResult
	// Subagent describes subagent lifecycle events.
	Subagent = model.Subagent
	// TurnEnd describes agent stop events.
	TurnEnd = model.TurnEnd
	// CompactInfo describes context compaction events.
	CompactInfo = model.CompactInfo
	// Note describes notifications and runtime errors.
	Note = model.Note
	// Lifecycle describes session start and end events.
	Lifecycle = model.Lifecycle
	// Decision is the unified gate verdict for pre-events.
	Decision = model.Decision
	// PreToolResult is the portable hook response for PreTool events.
	PreToolResult = model.PreToolResult
	// PostToolResult is the portable hook response for PostTool events.
	PostToolResult = model.PostToolResult
	// PostToolFailureResult is the portable hook response for PostToolFailure events.
	PostToolFailureResult = model.PostToolFailureResult
	// StopResult is the portable hook response for Stop and SubagentStop events.
	StopResult = model.StopResult
	// SessionStartResult is the portable hook response for SessionStart events.
	SessionStartResult = model.SessionStartResult
	// Codec decodes native payloads and encodes wire results.
	Codec = model.Codec
)

const (
	// Unknown is returned when the originating agent cannot be determined.
	Unknown = model.Unknown
	// Claude is the Claude Code agent dialect.
	Claude = model.Claude
	// Copilot is the GitHub Copilot CLI or cloud agent dialect.
	Copilot = model.Copilot
	// Cursor is the Cursor agent dialect.
	Cursor = model.Cursor

	// KindSessionStart is the normalized category for session start events.
	KindSessionStart = model.KindSessionStart
	// KindSessionEnd is the normalized category for session end events.
	KindSessionEnd = model.KindSessionEnd
	// KindUserPrompt is the normalized category for user prompt submission events.
	KindUserPrompt = model.KindUserPrompt
	// KindPreTool is the normalized category for pre-tool hook events.
	KindPreTool = model.KindPreTool
	// KindPostTool is the normalized category for successful post-tool events.
	KindPostTool = model.KindPostTool
	// KindPostToolFailure is the normalized category for failed post-tool events.
	KindPostToolFailure = model.KindPostToolFailure
	// KindSubagentStart is the normalized category for subagent start events.
	KindSubagentStart = model.KindSubagentStart
	// KindSubagentStop is the normalized category for subagent stop events.
	KindSubagentStop = model.KindSubagentStop
	// KindStop is the normalized category for agent stop events.
	KindStop = model.KindStop
	// KindPreCompact is the normalized category for pre-compaction events.
	KindPreCompact = model.KindPreCompact
	// KindOther is the normalized category for events with no dedicated mapping.
	KindOther = model.KindOther

	// DecisionUnset means the handler expressed no gate verdict.
	DecisionUnset = model.DecisionUnset
	// DecisionAllow permits the gated action to proceed.
	DecisionAllow = model.DecisionAllow
	// DecisionAsk escalates the decision to the user.
	DecisionAsk = model.DecisionAsk
	// DecisionDeny blocks the gated action.
	DecisionDeny = model.DecisionDeny

	// ToolBash is the normalized name for shell execution tools.
	ToolBash = model.ToolBash
	// ToolEdit is the normalized name for file edit tools.
	ToolEdit = model.ToolEdit
	// ToolWrite is the normalized name for file write tools.
	ToolWrite = model.ToolWrite
	// ToolRead is the normalized name for file read tools.
	ToolRead = model.ToolRead
	// ToolGlob is the normalized name for glob search tools.
	ToolGlob = model.ToolGlob
	// ToolGrep is the normalized name for grep search tools.
	ToolGrep = model.ToolGrep
	// ToolTask is the normalized name for subagent or task tools.
	ToolTask = model.ToolTask
	// ToolWebFetch is the normalized name for web fetch tools.
	ToolWebFetch = model.ToolWebFetch
	// ToolWebSearch is the normalized name for web search tools.
	ToolWebSearch = model.ToolWebSearch
	// ToolDelete is the normalized name for file delete tools.
	ToolDelete = model.ToolDelete
)

// ClaudeCodec implements Codec for Claude Code hooks.
type ClaudeCodec = claude.Codec

// CopilotCodec implements Codec for GitHub Copilot hooks.
type CopilotCodec = copilot.Codec

// CursorCodec implements Codec for Cursor hooks.
type CursorCodec = cursor.Codec

// CopilotPreToolErrorExit is the exit code when a preToolUse handler returns an error.
const CopilotPreToolErrorExit = copilot.PreToolErrorExit

// CopilotWarnExit is Copilot exit code 2.
const CopilotWarnExit = copilot.WarnExit

// CursorWarnExit is Cursor exit code 2.
const CursorWarnExit = cursor.WarnExit

// CursorHandlerErrorExit is Cursor exit code 1 for handler errors.
const CursorHandlerErrorExit = cursor.HandlerErrorExit

// ClaudeEventForKind maps unified kinds to Claude Code hook event names.
var ClaudeEventForKind = claude.EventForKind

// ClaudeKindForEventMap maps Claude hook event names to unified kinds.
var ClaudeKindForEventMap = claude.KindForEventMap

// ClaudeKindForEvent returns the unified kind for a Claude hook event name.
func ClaudeKindForEvent(name string) (Kind, bool) {
	return claude.KindForEvent(name)
}

// CopilotEventForKind maps unified kinds to GitHub Copilot hook event names.
var CopilotEventForKind = copilot.EventForKind

// CopilotKindForEventMap maps Copilot hook event names to unified kinds.
var CopilotKindForEventMap = copilot.KindForEventMap

// CopilotKindForEvent returns the unified kind for a Copilot hook event name.
func CopilotKindForEvent(name string) (Kind, bool) {
	return copilot.KindForEvent(name)
}

// CursorEventForKind maps unified kinds to Cursor hook event names.
var CursorEventForKind = cursor.EventForKind

// CursorKindForEventMap maps Cursor hook event names to unified kinds.
var CursorKindForEventMap = cursor.KindForEventMap

// CursorDedicatedEvents lists Cursor hook events folded into unified kinds.
var CursorDedicatedEvents = cursor.DedicatedEvents

// CursorKindForEvent returns the unified kind for a Cursor hook event name.
func CursorKindForEvent(name string) (Kind, bool) {
	return cursor.KindForEvent(name)
}

// IsCursorDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsCursorDedicatedEvent(name string) bool {
	return cursor.IsDedicatedEvent(name)
}

// ParseDialect parses a dialect name from a CLI flag or config value.
func ParseDialect(s string) Dialect { return model.ParseDialect(s) }

// Detect infers the originating agent from a hook payload and environment hints.
func Detect(payload []byte, getenv func(string) string) Dialect {
	return model.Detect(payload, getenv)
}

// InputAs decodes the native tool input into T.
func InputAs[T any](t *ToolCall) (T, error) { return model.InputAs[T](t) }

// NormalizeToolName maps a native tool name onto the canonical vocabulary.
func NormalizeToolName(native string) (name string, mcp bool) {
	return model.NormalizeToolName(native)
}

// PreToolAllow returns an allow verdict for PreTool events.
func PreToolAllow() PreToolResult { return model.PreToolAllow() }

// PreToolDeny returns a deny verdict for PreTool events.
func PreToolDeny(reason string) PreToolResult { return model.PreToolDeny(reason) }

// PreToolAsk returns an ask verdict for PreTool events.
func PreToolAsk(reason string) PreToolResult { return model.PreToolAsk(reason) }

// PostToolContext returns a context-injection-only PostTool result.
func PostToolContext(text string) PostToolResult { return model.PostToolContext(text) }

// PostToolFailureContext returns recovery guidance for PostToolFailure events.
func PostToolFailureContext(text string) PostToolFailureResult {
	return model.PostToolFailureContext(text)
}

// StopFollowUp returns a stop-gate result with a follow-up instruction.
func StopFollowUp(text string) StopResult { return model.StopFollowUp(text) }

// SessionStartContext returns a context-injection-only SessionStart result.
func SessionStartContext(text string) SessionStartResult { return model.SessionStartContext(text) }
