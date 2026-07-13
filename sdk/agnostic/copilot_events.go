package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

// CopilotEventForKind maps unified kinds to GitHub Copilot hook event names.
var CopilotEventForKind = map[Kind]string{
	KindSessionStart:      copilot.EventSessionStart,
	KindSessionEnd:        copilot.EventSessionEnd,
	KindUserPrompt:        copilot.EventUserPromptSubmitted,
	KindPreTool:           copilot.EventPreToolUse,
	KindPostTool:          copilot.EventPostToolUse,
	KindPostToolFailure:   copilot.EventPostToolUseFailure,
	KindPermissionRequest: copilot.EventPermissionRequest,
	KindSubagentStart:     copilot.EventSubagentStart,
	KindSubagentStop:      copilot.EventSubagentStop,
	KindStop:              copilot.EventAgentStop,
	KindPreCompact:        copilot.EventPreCompact,
	KindNotification:      copilot.EventNotification,
	KindAgentError:        copilot.EventErrorOccurred,
}

// CopilotKindForEventMap maps Copilot hook event names to unified kinds.
var CopilotKindForEventMap = buildCopilotKindForEvent()

// CopilotKindForEvent returns the unified kind for a Copilot hook event name.
func CopilotKindForEvent(name string) (Kind, bool) {
	canonical, ok := copilot.CanonicalEventName(name)
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
	for alias, canonical := range copilot.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
