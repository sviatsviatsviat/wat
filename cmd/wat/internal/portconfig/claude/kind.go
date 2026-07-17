package claude

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// EventForKind maps unified kinds to Claude Code hook event names.
var EventForKind = map[model.Kind]string{
	model.KindSessionStart:      sdkclaude.EventSessionStart,
	model.KindSessionEnd:        sdkclaude.EventSessionEnd,
	model.KindUserPrompt:        sdkclaude.EventUserPromptSubmit,
	model.KindPreTool:           sdkclaude.EventPreToolUse,
	model.KindPostTool:          sdkclaude.EventPostToolUse,
	model.KindPostToolFailure:   sdkclaude.EventPostToolUseFailure,
	model.KindPermissionRequest: sdkclaude.EventPermissionRequest,
	model.KindSubagentStart:     sdkclaude.EventSubagentStart,
	model.KindSubagentStop:      sdkclaude.EventSubagentStop,
	model.KindStop:              sdkclaude.EventStop,
	model.KindPreCompact:        sdkclaude.EventPreCompact,
	model.KindNotification:      sdkclaude.EventNotification,
	model.KindAgentError:        sdkclaude.EventStopFailure,
}

var kindForEventMap = model.InvertEventForKind(EventForKind)

func kindForEvent(name string) (model.Kind, bool) {
	kind, ok := kindForEventMap[name]
	return kind, ok
}
