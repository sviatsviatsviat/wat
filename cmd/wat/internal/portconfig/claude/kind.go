package claude

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// EventForKind maps unified kinds to Claude Code hook event names.
var EventForKind = map[agnostic.Kind]string{
	agnostic.KindSessionStart:      sdkclaude.EventSessionStart,
	agnostic.KindSessionEnd:        sdkclaude.EventSessionEnd,
	agnostic.KindUserPrompt:        sdkclaude.EventUserPromptSubmit,
	agnostic.KindPreTool:           sdkclaude.EventPreToolUse,
	agnostic.KindPostTool:          sdkclaude.EventPostToolUse,
	agnostic.KindPostToolFailure:   sdkclaude.EventPostToolUseFailure,
	agnostic.KindPermissionRequest: sdkclaude.EventPermissionRequest,
	agnostic.KindSubagentStart:     sdkclaude.EventSubagentStart,
	agnostic.KindSubagentStop:      sdkclaude.EventSubagentStop,
	agnostic.KindStop:              sdkclaude.EventStop,
	agnostic.KindPreCompact:        sdkclaude.EventPreCompact,
	agnostic.KindNotification:      sdkclaude.EventNotification,
	agnostic.KindAgentError:        sdkclaude.EventStopFailure,
}

var kindForEventMap = model.InvertEventForKind(EventForKind)

func kindForEvent(name string) (agnostic.Kind, bool) {
	kind, ok := kindForEventMap[name]
	return kind, ok
}
