package agnostic

import (
	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
)

// ClaudeCodec adapts Claude Code hook payloads to the unified agnostic model.
type ClaudeCodec = agclaude.Codec

// CopilotCodec adapts GitHub Copilot hook payloads to the unified agnostic model.
type CopilotCodec = agcopilot.Codec

// CursorCodec adapts Cursor hook payloads to the unified agnostic model.
type CursorCodec = agcursor.Codec

// CopilotPreToolErrorExit is the exit code when a Copilot preToolUse handler returns an error.
const CopilotPreToolErrorExit = agcopilot.PreToolErrorExit

// CopilotWarnExit is Copilot exit code 2 for documented warn/deny paths.
const CopilotWarnExit = agcopilot.WarnExit

// CursorWarnExit is Cursor exit code 2 for block/deny permission paths.
const CursorWarnExit = agcursor.WarnExit

// CursorHandlerErrorExit is Cursor exit code 1 for handler errors under fail-open policy.
const CursorHandlerErrorExit = agcursor.HandlerErrorExit

// ClaudeEventForKind maps unified kinds to Claude hook event names.
var ClaudeEventForKind = agclaude.EventForKind

// ClaudeKindForEventMap maps Claude hook event names to unified kinds.
var ClaudeKindForEventMap = agclaude.KindForEventMap

// ClaudeKindForEvent returns the unified kind for a Claude hook event name.
func ClaudeKindForEvent(name string) (Kind, bool) {
	return agclaude.KindForEvent(name)
}

// CopilotEventForKind maps unified kinds to Copilot hook event names.
var CopilotEventForKind = agcopilot.EventForKind

// CopilotKindForEventMap maps Copilot hook event names to unified kinds.
var CopilotKindForEventMap = agcopilot.KindForEventMap

// CopilotKindForEvent returns the unified kind for a Copilot hook event name.
func CopilotKindForEvent(name string) (Kind, bool) {
	return agcopilot.KindForEvent(name)
}

// CursorEventForKind maps unified kinds to Cursor hook event names.
var CursorEventForKind = agcursor.EventForKind

// CursorKindForEventMap maps Cursor hook event names to unified kinds.
var CursorKindForEventMap = agcursor.KindForEventMap

// CursorDedicatedEvents lists Cursor-only surface events.
var CursorDedicatedEvents = agcursor.DedicatedEvents

// CursorKindForEvent returns the unified kind for a Cursor hook event name.
func CursorKindForEvent(name string) (Kind, bool) {
	return agcursor.KindForEvent(name)
}

// IsCursorDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsCursorDedicatedEvent(name string) bool {
	return agcursor.IsDedicatedEvent(name)
}
