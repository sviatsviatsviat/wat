package agnostic

import "github.com/sviatsviatsviat/wat/sdk/cursor"

// CursorEventForKind maps unified kinds to Cursor hook event names.
var CursorEventForKind = map[Kind]string{
	KindSessionStart:    cursor.EventSessionStart,
	KindSessionEnd:      cursor.EventSessionEnd,
	KindUserPrompt:      cursor.EventBeforeSubmitPrompt,
	KindPreTool:         cursor.EventPreToolUse,
	KindPostTool:        cursor.EventPostToolUse,
	KindPostToolFailure: cursor.EventPostToolUseFailure,
	KindSubagentStart:   cursor.EventSubagentStart,
	KindSubagentStop:    cursor.EventSubagentStop,
	KindStop:            cursor.EventStop,
	KindPreCompact:      cursor.EventPreCompact,
}

// CursorKindForEventMap maps Cursor hook event names to unified kinds.
var CursorKindForEventMap = buildCursorKindForEvent()

// CursorDedicatedEvents lists Cursor hook events folded into unified kinds but
// not portable as dedicated events on other agents.
var CursorDedicatedEvents = map[string]bool{
	cursor.EventBeforeShellExecution: true,
	cursor.EventAfterShellExecution:  true,
	cursor.EventBeforeMCPExecution:   true,
	cursor.EventAfterMCPExecution:    true,
	cursor.EventBeforeReadFile:       true,
	cursor.EventAfterFileEdit:        true,
}

// CursorKindForEvent returns the unified kind for a Cursor hook event name.
func CursorKindForEvent(name string) (Kind, bool) {
	kind, ok := CursorKindForEventMap[name]
	return kind, ok
}

// IsCursorDedicatedEvent reports whether name is a Cursor dedicated surface event.
func IsCursorDedicatedEvent(name string) bool {
	return CursorDedicatedEvents[name]
}

func buildCursorKindForEvent() map[string]Kind {
	out := make(map[string]Kind, len(CursorEventForKind)+12)
	for kind, event := range CursorEventForKind {
		out[event] = kind
	}
	for event, kind := range map[string]Kind{
		cursor.EventBeforeShellExecution: KindPreTool,
		cursor.EventAfterShellExecution:  KindPostTool,
		cursor.EventBeforeMCPExecution:   KindPreTool,
		cursor.EventAfterMCPExecution:    KindPostTool,
		cursor.EventBeforeReadFile:       KindPreTool,
		cursor.EventAfterFileEdit:        KindPostTool,
		cursor.EventAfterAgentResponse:   KindOther,
		cursor.EventAfterAgentThought:    KindOther,
		cursor.EventBeforeTabFileRead:    KindOther,
		cursor.EventAfterTabFileEdit:     KindOther,
		cursor.EventWorkspaceOpen:        KindOther,
	} {
		out[event] = kind
	}
	return out
}
