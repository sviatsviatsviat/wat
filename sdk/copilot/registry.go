package copilot

// Canonical PascalCase GitHub Copilot hook event names for config keys and mux dispatch.
const (
	EventSessionStart        = "SessionStart"
	EventSessionEnd          = "SessionEnd"
	EventUserPromptSubmitted = "UserPromptSubmit"
	EventPreToolUse          = "PreToolUse"
	EventPostToolUse         = "PostToolUse"
	EventPostToolUseFailure  = "PostToolUseFailure"
	EventPermissionRequest   = "PermissionRequest"
	EventSubagentStart       = "SubagentStart"
	EventSubagentStop        = "SubagentStop"
	EventAgentStop           = "Stop"
	EventPreCompact          = "PreCompact"
	EventNotification        = "Notification"
	EventErrorOccurred       = "ErrorOccurred"
)

// knownEvents lists canonical decoder/handler keys.
var knownEvents = map[string]struct{}{
	EventSessionStart:        {},
	EventSessionEnd:          {},
	EventUserPromptSubmitted: {},
	EventPreToolUse:          {},
	EventPostToolUse:         {},
	EventPostToolUseFailure:  {},
	EventPermissionRequest:   {},
	EventSubagentStart:       {},
	EventSubagentStop:        {},
	EventAgentStop:           {},
	EventPreCompact:          {},
	EventNotification:        {},
	EventErrorOccurred:       {},
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
