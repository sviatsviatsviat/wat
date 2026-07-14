package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// EventForKind maps unified kinds to GitHub Copilot hook event names.
var EventForKind = map[model.Kind]string{
	model.KindSessionStart:      sdkcopilot.EventSessionStart,
	model.KindSessionEnd:        sdkcopilot.EventSessionEnd,
	model.KindUserPrompt:        sdkcopilot.EventUserPromptSubmitted,
	model.KindPreTool:           sdkcopilot.EventPreToolUse,
	model.KindPostTool:          sdkcopilot.EventPostToolUse,
	model.KindPostToolFailure:   sdkcopilot.EventPostToolUseFailure,
	model.KindPermissionRequest: sdkcopilot.EventPermissionRequest,
	model.KindSubagentStart:     sdkcopilot.EventSubagentStart,
	model.KindSubagentStop:      sdkcopilot.EventSubagentStop,
	model.KindStop:              sdkcopilot.EventAgentStop,
	model.KindPreCompact:        sdkcopilot.EventPreCompact,
	model.KindNotification:      sdkcopilot.EventNotification,
	model.KindAgentError:        sdkcopilot.EventErrorOccurred,
}

// KindForEventMap maps Copilot hook event names to unified kinds.
var KindForEventMap = buildKindForEvent()

// KindForEvent returns the unified kind for a Copilot hook event name.
func KindForEvent(name string) (model.Kind, bool) {
	canonical, ok := sdkcopilot.CanonicalEventName(name)
	if !ok {
		return model.KindOther, false
	}
	kind, ok := KindForEventMap[canonical]
	return kind, ok
}

func buildKindForEvent() map[string]model.Kind {
	out := make(map[string]model.Kind, len(EventForKind))
	for kind, event := range EventForKind {
		out[event] = kind
	}
	for alias, canonical := range sdkcopilot.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
