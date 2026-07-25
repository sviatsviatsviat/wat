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

var kindForEventMap = model.InvertEventForKind(EventForKind)

func kindForEvent(name string) (model.Kind, bool) {
	kind, ok := kindForEventMap[name]
	return kind, ok
}
