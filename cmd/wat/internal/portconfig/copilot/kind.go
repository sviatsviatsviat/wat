package copilot

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
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

var kindForEventMap = buildKindForEvent()

func kindForEvent(name string) (model.Kind, bool) {
	canonical, ok := sdkcopilot.CanonicalEventName(name)
	if !ok {
		return model.KindOther, false
	}
	kind, ok := kindForEventMap[canonical]
	return kind, ok
}

func buildKindForEvent() map[string]model.Kind {
	out := model.InvertEventForKind(EventForKind)
	for alias, canonical := range sdkcopilot.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
