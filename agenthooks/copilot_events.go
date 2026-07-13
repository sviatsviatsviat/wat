package agenthooks

import (
	"github.com/sviatsviatsviat/wat/copilothook"
)

// CopilotEventForKind maps unified kinds to GitHub Copilot hook event names.
var CopilotEventForKind = map[Kind]string{
	KindSessionStart:      copilothook.EventSessionStart,
	KindSessionEnd:        copilothook.EventSessionEnd,
	KindUserPrompt:        copilothook.EventUserPromptSubmitted,
	KindPreTool:           copilothook.EventPreToolUse,
	KindPostTool:          copilothook.EventPostToolUse,
	KindPostToolFailure:   copilothook.EventPostToolUseFailure,
	KindPermissionRequest: copilothook.EventPermissionRequest,
	KindSubagentStart:     copilothook.EventSubagentStart,
	KindSubagentStop:      copilothook.EventSubagentStop,
	KindStop:              copilothook.EventAgentStop,
	KindPreCompact:        copilothook.EventPreCompact,
	KindNotification:      copilothook.EventNotification,
	KindAgentError:        copilothook.EventErrorOccurred,
}

// CopilotKindForEventMap maps Copilot hook event names to unified kinds.
var CopilotKindForEventMap = buildCopilotKindForEvent()

// CopilotKindForEvent returns the unified kind for a Copilot hook event name.
func CopilotKindForEvent(name string) (Kind, bool) {
	canonical, ok := copilothook.CanonicalEventName(name)
	if !ok {
		return KindOther, false
	}
	kind, ok := CopilotKindForEventMap[canonical]
	return kind, ok
}

func buildCopilotKindForEvent() map[string]Kind {
	out := make(map[string]Kind, len(CopilotEventForKind))
	for kind, event := range CopilotEventForKind {
		out[event] = kind
	}
	for alias, canonical := range copilothook.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
