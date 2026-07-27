package event

// StopActiveFields holds Claude stop-hook continuation fields shared by Stop
// and SubagentStop events.
type StopActiveFields struct {
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
}
