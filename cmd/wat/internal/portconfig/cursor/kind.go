package cursor

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// EventForKind maps unified kinds to Cursor hook event names.
var EventForKind = map[agnostic.Kind]string{
	agnostic.KindSessionStart:    sdkcursor.EventSessionStart,
	agnostic.KindSessionEnd:      sdkcursor.EventSessionEnd,
	agnostic.KindUserPrompt:      sdkcursor.EventBeforeSubmitPrompt,
	agnostic.KindPreTool:         sdkcursor.EventPreToolUse,
	agnostic.KindPostTool:        sdkcursor.EventPostToolUse,
	agnostic.KindPostToolFailure: sdkcursor.EventPostToolUseFailure,
	agnostic.KindSubagentStart:   sdkcursor.EventSubagentStart,
	agnostic.KindSubagentStop:    sdkcursor.EventSubagentStop,
	agnostic.KindStop:            sdkcursor.EventStop,
	agnostic.KindPreCompact:      sdkcursor.EventPreCompact,
}

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

// IsDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsDedicatedEvent(name string) bool {
	return DedicatedEvents[name]
}

var kindForEventMap = buildKindForEvent()

func kindForEvent(name string) (agnostic.Kind, bool) {
	kind, ok := kindForEventMap[name]
	return kind, ok
}

func buildKindForEvent() map[string]agnostic.Kind {
	out := model.InvertEventForKind(EventForKind)
	for event, kind := range map[string]agnostic.Kind{
		sdkcursor.EventBeforeShellExecution: agnostic.KindPreTool,
		sdkcursor.EventAfterShellExecution:  agnostic.KindPostTool,
		sdkcursor.EventBeforeMCPExecution:   agnostic.KindPreTool,
		sdkcursor.EventAfterMCPExecution:    agnostic.KindPostTool,
		sdkcursor.EventBeforeReadFile:       agnostic.KindPreTool,
		sdkcursor.EventAfterFileEdit:        agnostic.KindPostTool,
		sdkcursor.EventAfterAgentResponse:   agnostic.KindOther,
		sdkcursor.EventAfterAgentThought:    agnostic.KindOther,
		sdkcursor.EventBeforeTabFileRead:    agnostic.KindOther,
		sdkcursor.EventAfterTabFileEdit:     agnostic.KindOther,
		sdkcursor.EventWorkspaceOpen:        agnostic.KindOther,
	} {
		out[event] = kind
	}
	return out
}
