package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// EventForKind maps unified kinds to Cursor hook event names.
var EventForKind = map[model.Kind]string{
	model.KindSessionStart:    sdkcursor.EventSessionStart,
	model.KindSessionEnd:      sdkcursor.EventSessionEnd,
	model.KindUserPrompt:      sdkcursor.EventBeforeSubmitPrompt,
	model.KindPreTool:         sdkcursor.EventPreToolUse,
	model.KindPostTool:        sdkcursor.EventPostToolUse,
	model.KindPostToolFailure: sdkcursor.EventPostToolUseFailure,
	model.KindSubagentStart:   sdkcursor.EventSubagentStart,
	model.KindSubagentStop:    sdkcursor.EventSubagentStop,
	model.KindStop:            sdkcursor.EventStop,
	model.KindPreCompact:      sdkcursor.EventPreCompact,
}

// KindForEventMap maps Cursor hook event names to unified kinds.
var KindForEventMap = buildKindForEvent()

// DedicatedEvents lists Cursor hook events folded into unified kinds but not
// portable as dedicated events on other agents.
var DedicatedEvents = map[string]bool{
	sdkcursor.EventBeforeShellExecution: true,
	sdkcursor.EventAfterShellExecution:  true,
	sdkcursor.EventBeforeMCPExecution:   true,
	sdkcursor.EventAfterMCPExecution:    true,
	sdkcursor.EventBeforeReadFile:       true,
	sdkcursor.EventAfterFileEdit:        true,
}

// KindForEvent returns the unified kind for a Cursor hook event name.
func KindForEvent(name string) (model.Kind, bool) {
	kind, ok := KindForEventMap[name]
	return kind, ok
}

// IsDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsDedicatedEvent(name string) bool {
	return DedicatedEvents[name]
}

func buildKindForEvent() map[string]model.Kind {
	out := make(map[string]model.Kind, len(EventForKind)+12)
	for kind, event := range EventForKind {
		out[event] = kind
	}
	for event, kind := range map[string]model.Kind{
		sdkcursor.EventBeforeShellExecution: model.KindPreTool,
		sdkcursor.EventAfterShellExecution:  model.KindPostTool,
		sdkcursor.EventBeforeMCPExecution:   model.KindPreTool,
		sdkcursor.EventAfterMCPExecution:    model.KindPostTool,
		sdkcursor.EventBeforeReadFile:       model.KindPreTool,
		sdkcursor.EventAfterFileEdit:        model.KindPostTool,
		sdkcursor.EventAfterAgentResponse:   model.KindOther,
		sdkcursor.EventAfterAgentThought:    model.KindOther,
		sdkcursor.EventBeforeTabFileRead:    model.KindOther,
		sdkcursor.EventAfterTabFileEdit:     model.KindOther,
		sdkcursor.EventWorkspaceOpen:        model.KindOther,
	} {
		out[event] = kind
	}
	return out
}
