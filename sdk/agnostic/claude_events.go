package agnostic

import "github.com/sviatsviatsviat/wat/sdk/claude"

// ClaudeEventForKind maps unified kinds to Claude Code hook event names.
var ClaudeEventForKind = map[Kind]string{
	KindSessionStart:      claude.EventSessionStart,
	KindSessionEnd:        claude.EventSessionEnd,
	KindUserPrompt:        claude.EventUserPromptSubmit,
	KindPreTool:           claude.EventPreToolUse,
	KindPostTool:          claude.EventPostToolUse,
	KindPostToolFailure:   claude.EventPostToolUseFailure,
	KindPermissionRequest: claude.EventPermissionRequest,
	KindSubagentStart:     claude.EventSubagentStart,
	KindSubagentStop:      claude.EventSubagentStop,
	KindStop:              claude.EventStop,
	KindPreCompact:        claude.EventPreCompact,
	KindNotification:      claude.EventNotification,
	KindAgentError:        claude.EventStopFailure,
}

// ClaudeKindForEventMap maps Claude hook event names to unified kinds.
var ClaudeKindForEventMap = invertKindEvent(ClaudeEventForKind)

// ClaudeKindForEvent returns the unified kind for a Claude hook event name.
func ClaudeKindForEvent(name string) (Kind, bool) {
	kind, ok := ClaudeKindForEventMap[name]
	return kind, ok
}

func invertKindEvent(m map[Kind]string) map[string]Kind {
	out := make(map[string]Kind, len(m))
	for kind, event := range m {
		out[event] = kind
	}
	return out
}
