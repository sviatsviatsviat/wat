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
