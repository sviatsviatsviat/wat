package agenthooks

import "github.com/sviatsviatsviat/wat/cursorhook"

// CursorEventForKind maps unified kinds to Cursor hook event names.
var CursorEventForKind = map[Kind]string{
	KindSessionStart:    cursorhook.EventSessionStart,
	KindSessionEnd:      cursorhook.EventSessionEnd,
	KindUserPrompt:      cursorhook.EventBeforeSubmitPrompt,
	KindPreTool:         cursorhook.EventPreToolUse,
	KindPostTool:        cursorhook.EventPostToolUse,
	KindPostToolFailure: cursorhook.EventPostToolUseFailure,
	KindSubagentStart:   cursorhook.EventSubagentStart,
	KindSubagentStop:    cursorhook.EventSubagentStop,
	KindStop:            cursorhook.EventStop,
	KindPreCompact:      cursorhook.EventPreCompact,
}

// CursorKindForEventMap maps Cursor hook event names to unified kinds.
var CursorKindForEventMap = buildCursorKindForEvent()

// CursorDedicatedEvents lists Cursor hook events folded into unified kinds but
// not portable as dedicated events on other agents.
var CursorDedicatedEvents = map[string]bool{
	cursorhook.EventBeforeShellExecution: true,
	cursorhook.EventAfterShellExecution:  true,
	cursorhook.EventBeforeMCPExecution:   true,
	cursorhook.EventAfterMCPExecution:    true,
	cursorhook.EventBeforeReadFile:       true,
	cursorhook.EventAfterFileEdit:        true,
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
		cursorhook.EventBeforeShellExecution: KindPreTool,
		cursorhook.EventAfterShellExecution:  KindPostTool,
		cursorhook.EventBeforeMCPExecution:   KindPreTool,
		cursorhook.EventAfterMCPExecution:    KindPostTool,
		cursorhook.EventBeforeReadFile:       KindPreTool,
		cursorhook.EventAfterFileEdit:        KindPostTool,
		cursorhook.EventAfterAgentResponse:   KindOther,
		cursorhook.EventAfterAgentThought:    KindOther,
		cursorhook.EventBeforeTabFileRead:    KindOther,
		cursorhook.EventAfterTabFileEdit:     KindOther,
		cursorhook.EventWorkspaceOpen:        KindOther,
	} {
		out[event] = kind
	}
	return out
}
