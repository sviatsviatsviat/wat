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
// The map is a snapshot; mutating it does not affect the canonical registry.
var ClaudeEventForKind = cloneEventForKind(agclaude.EventForKind)

// ClaudeKindForEventMap maps Claude hook event names to unified kinds.
// The map is a snapshot; mutating it does not affect the canonical registry.
var ClaudeKindForEventMap = cloneKindForEvent(agclaude.KindForEventMap)

// ClaudeKindForEvent returns the unified kind for a Claude hook event name.
func ClaudeKindForEvent(name string) (Kind, bool) {
	return agclaude.KindForEvent(name)
}

// CopilotEventForKind maps unified kinds to Copilot hook event names.
// The map is a snapshot; mutating it does not affect the canonical registry.
var CopilotEventForKind = cloneEventForKind(agcopilot.EventForKind)

// CopilotKindForEventMap maps Copilot hook event names to unified kinds.
// The map is a snapshot; mutating it does not affect the canonical registry.
var CopilotKindForEventMap = cloneKindForEvent(agcopilot.KindForEventMap)

// CopilotKindForEvent returns the unified kind for a Copilot hook event name.
func CopilotKindForEvent(name string) (Kind, bool) {
	return agcopilot.KindForEvent(name)
}

// CursorEventForKind maps unified kinds to Cursor hook event names.
// The map is a snapshot; mutating it does not affect the canonical registry.
var CursorEventForKind = cloneEventForKind(agcursor.EventForKind)

// CursorKindForEventMap maps Cursor hook event names to unified kinds.
// The map is a snapshot; mutating it does not affect the canonical registry.
var CursorKindForEventMap = cloneKindForEvent(agcursor.KindForEventMap)

// CursorDedicatedEvents lists Cursor-only surface events.
// The map is a snapshot; mutating it does not affect the canonical registry.
var CursorDedicatedEvents = cloneDedicatedEvents(agcursor.DedicatedEvents)

// CursorKindForEvent returns the unified kind for a Cursor hook event name.
func CursorKindForEvent(name string) (Kind, bool) {
	return agcursor.KindForEvent(name)
}

// IsCursorDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsCursorDedicatedEvent(name string) bool {
	return agcursor.IsDedicatedEvent(name)
}

func cloneEventForKind(src map[Kind]string) map[Kind]string {
	out := make(map[Kind]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func cloneKindForEvent(src map[string]Kind) map[string]Kind {
	out := make(map[string]Kind, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func cloneDedicatedEvents(src map[string]bool) map[string]bool {
	out := make(map[string]bool, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
