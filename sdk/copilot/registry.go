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

// CanonicalEventName normalizes a wire event name to the canonical form.
func CanonicalEventName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	if _, ok := knownEvents[name]; ok {
		return name, true
	}
	return name, false
}

// ResolveCanonical maps a received event name to a canonical decoder key using
// payload scope when wire names are ambiguous.
func ResolveCanonical(raw []byte, received string) (string, bool) {
	if received == "Stop" {
		return resolveCanonical(received, payloadHasSubagentScope(raw))
	}
	return CanonicalEventName(received)
}

func resolveCanonical(received string, hasSubagentScope bool) (string, bool) {
	if received == "Stop" {
		if hasSubagentScope {
			return EventSubagentStop, true
		}
		return EventAgentStop, true
	}
	return CanonicalEventName(received)
}

func payloadHasSubagentScope(raw []byte) bool {
	peek, err := peekPayload(raw)
	if err != nil {
		return false
	}
	return peek.hasSubagentScope()
}

// EventAliasMap returns a copy of known event name to itself (identity map).
func EventAliasMap() map[string]string {
	out := make(map[string]string, len(knownEvents))
	for k := range knownEvents {
		out[k] = k
	}
	return out
}
