package copilot

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// EventForKind maps unified kinds to GitHub Copilot hook event names.
var EventForKind = map[agnostic.Kind]string{
	agnostic.KindSessionStart:      sdkcopilot.EventSessionStart,
	agnostic.KindSessionEnd:        sdkcopilot.EventSessionEnd,
	agnostic.KindUserPrompt:        sdkcopilot.EventUserPromptSubmitted,
	agnostic.KindPreTool:           sdkcopilot.EventPreToolUse,
	agnostic.KindPostTool:          sdkcopilot.EventPostToolUse,
	agnostic.KindPostToolFailure:   sdkcopilot.EventPostToolUseFailure,
	agnostic.KindPermissionRequest: sdkcopilot.EventPermissionRequest,
	agnostic.KindSubagentStart:     sdkcopilot.EventSubagentStart,
	agnostic.KindSubagentStop:      sdkcopilot.EventSubagentStop,
	agnostic.KindStop:              sdkcopilot.EventAgentStop,
	agnostic.KindPreCompact:        sdkcopilot.EventPreCompact,
	agnostic.KindNotification:      sdkcopilot.EventNotification,
	agnostic.KindAgentError:        sdkcopilot.EventErrorOccurred,
}

var kindForEventMap = buildKindForEvent()

func kindForEvent(name string) (agnostic.Kind, bool) {
	canonical, ok := sdkcopilot.CanonicalEventName(name)
	if !ok {
		return agnostic.KindOther, false
	}
	kind, ok := kindForEventMap[canonical]
	return kind, ok
}

func buildKindForEvent() map[string]agnostic.Kind {
	out := model.InvertEventForKind(EventForKind)
	for alias, canonical := range sdkcopilot.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
