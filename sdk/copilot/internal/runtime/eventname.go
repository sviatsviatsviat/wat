package runtime

import "github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"

// knownEvents lists canonical decoder/handler keys.
var knownEvents = map[string]struct{}{
	event.SessionStart:        {},
	event.SessionEnd:          {},
	event.UserPromptSubmitted: {},
	event.PreToolUse:          {},
	event.PostToolUse:         {},
	event.PostToolUseFailure:  {},
	event.PermissionRequest:   {},
	event.SubagentStart:       {},
	event.SubagentStop:        {},
	event.AgentStop:           {},
	event.PreCompact:          {},
	event.Notification:        {},
	event.ErrorOccurred:       {},
}

// CanonicalEventName reports whether name is a known Copilot hook event name.
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
