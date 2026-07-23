package event

// GitHub Copilot hook_event_name values for config keys and stdin payloads.
const (
	SessionStart        = "SessionStart"
	SessionEnd          = "SessionEnd"
	UserPromptSubmitted = "UserPromptSubmit"
	PreToolUse          = "PreToolUse"
	PostToolUse         = "PostToolUse"
	PostToolUseFailure  = "PostToolUseFailure"
	PermissionRequest   = "PermissionRequest"
	SubagentStart       = "SubagentStart"
	SubagentStop        = "SubagentStop"
	AgentStop           = "Stop"
	PreCompact          = "PreCompact"
	Notification        = "Notification"
	ErrorOccurred       = "ErrorOccurred"
)
