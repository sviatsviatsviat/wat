package agenthooks

import "github.com/sviatsviatsviat/wat/claudehook"

// ClaudeEventForKind maps unified kinds to Claude Code hook event names.
var ClaudeEventForKind = map[Kind]string{
	KindSessionStart:      claudehook.EventSessionStart,
	KindSessionEnd:        claudehook.EventSessionEnd,
	KindUserPrompt:        claudehook.EventUserPromptSubmit,
	KindPreTool:           claudehook.EventPreToolUse,
	KindPostTool:          claudehook.EventPostToolUse,
	KindPostToolFailure:   claudehook.EventPostToolUseFailure,
	KindPermissionRequest: claudehook.EventPermissionRequest,
	KindSubagentStart:     claudehook.EventSubagentStart,
	KindSubagentStop:      claudehook.EventSubagentStop,
	KindStop:              claudehook.EventStop,
	KindPreCompact:        claudehook.EventPreCompact,
	KindNotification:      claudehook.EventNotification,
	KindAgentError:        claudehook.EventStopFailure,
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
