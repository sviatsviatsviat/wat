package runtime

import "github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"

// knownEvents lists canonical decoder/handler keys.
var knownEvents = map[string]struct{}{
	event.SessionStart:         {},
	event.SessionEnd:           {},
	event.BeforeSubmitPrompt:   {},
	event.PreToolUse:           {},
	event.PostToolUse:          {},
	event.PostToolUseFailure:   {},
	event.BeforeShellExecution: {},
	event.AfterShellExecution:  {},
	event.BeforeMCPExecution:   {},
	event.AfterMCPExecution:    {},
	event.BeforeReadFile:       {},
	event.AfterFileEdit:        {},
	event.SubagentStart:        {},
	event.SubagentStop:         {},
	event.Stop:                 {},
	event.PreCompact:           {},
	event.AfterAgentResponse:   {},
	event.AfterAgentThought:    {},
	event.BeforeTabFileRead:    {},
	event.AfterTabFileEdit:     {},
	event.WorkspaceOpen:        {},
}

// CanonicalEventName reports whether name is a known Cursor hook event name.
// Known names are returned unchanged; unknown non-empty names return false.
func CanonicalEventName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	if _, ok := knownEvents[name]; ok {
		return name, true
	}
	return name, false
}

// EventAliasMap returns a copy of known event name to itself (identity map).
func EventAliasMap() map[string]string {
	out := make(map[string]string, len(knownEvents))
	for k := range knownEvents {
		out[k] = k
	}
	return out
}
