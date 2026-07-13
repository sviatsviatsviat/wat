package copilothook

import "encoding/json"

// Canonical camelCase GitHub Copilot hook event names for config keys and mux dispatch.
const (
	EventSessionStart        = "sessionStart"
	EventSessionEnd          = "sessionEnd"
	EventUserPromptSubmitted = "userPromptSubmitted"
	EventPreToolUse          = "preToolUse"
	EventPostToolUse         = "postToolUse"
	EventPostToolUseFailure  = "postToolUseFailure"
	EventPermissionRequest   = "permissionRequest"
	EventSubagentStart       = "subagentStart"
	EventSubagentStop        = "subagentStop"
	EventAgentStop           = "agentStop"
	EventPreCompact          = "preCompact"
	EventNotification        = "notification"
	EventErrorOccurred       = "errorOccurred"
)

// eventAliases maps wire event names (camelCase and VS Code PascalCase) to canonical names.
var eventAliases = map[string]string{
	EventSessionStart:        EventSessionStart,
	"SessionStart":           EventSessionStart,
	EventSessionEnd:          EventSessionEnd,
	"SessionEnd":             EventSessionEnd,
	EventUserPromptSubmitted: EventUserPromptSubmitted,
	"UserPromptSubmit":       EventUserPromptSubmitted,
	EventPreToolUse:          EventPreToolUse,
	"PreToolUse":             EventPreToolUse,
	EventPostToolUse:         EventPostToolUse,
	"PostToolUse":            EventPostToolUse,
	EventPostToolUseFailure:  EventPostToolUseFailure,
	"PostToolUseFailure":     EventPostToolUseFailure,
	EventPermissionRequest:   EventPermissionRequest,
	"PermissionRequest":      EventPermissionRequest,
	EventSubagentStart:       EventSubagentStart,
	"SubagentStart":          EventSubagentStart,
	EventSubagentStop:        EventSubagentStop,
	"SubagentStop":           EventSubagentStop,
	EventAgentStop:           EventAgentStop,
	"Stop":                   EventAgentStop,
	EventPreCompact:          EventPreCompact,
	"PreCompact":             EventPreCompact,
	EventNotification:        EventNotification,
	"Notification":           EventNotification,
	EventErrorOccurred:       EventErrorOccurred,
	"ErrorOccurred":          EventErrorOccurred,
}

// CanonicalEventName normalizes a wire event name to the canonical camelCase form.
func CanonicalEventName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	if canonical, ok := eventAliases[name]; ok {
		return canonical, true
	}
	return name, false
}

// ResolveCanonical maps a received event name to a canonical decoder key using
// payload scope when wire names are ambiguous.
func ResolveCanonical(raw []byte, received string) (string, bool) {
	if received == "Stop" {
		if payloadHasSubagentScope(raw) {
			return EventSubagentStop, true
		}
		return EventAgentStop, true
	}
	return CanonicalEventName(received)
}

func payloadHasSubagentScope(raw []byte) bool {
	var peek struct {
		AgentName             string `json:"agent_name"`
		AgentNameCamel        string `json:"agentName"`
		AgentDisplayName      string `json:"agent_display_name"`
		AgentDisplayNameCamel string `json:"agentDisplayName"`
	}
	if json.Unmarshal(raw, &peek) != nil {
		return false
	}
	return peek.AgentName != "" || peek.AgentNameCamel != "" ||
		peek.AgentDisplayName != "" || peek.AgentDisplayNameCamel != ""
}

// EventAliasMap returns a copy of wire-name to canonical-name aliases.
func EventAliasMap() map[string]string {
	out := make(map[string]string, len(eventAliases))
	for k, v := range eventAliases {
		out[k] = v
	}
	return out
}
