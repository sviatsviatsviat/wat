package model

// Kind is a normalized hook event category used by portconfig when mapping
// native config event names across Claude, Copilot, and Cursor.
type Kind string

const (
	// KindSessionStart is the normalized category for session start events.
	KindSessionStart Kind = "SessionStart"
	// KindSessionEnd is the normalized category for session end events.
	KindSessionEnd Kind = "SessionEnd"
	// KindUserPrompt is the normalized category for user prompt submission events.
	KindUserPrompt Kind = "UserPrompt"
	// KindPreTool is the normalized category for pre-tool hook events.
	KindPreTool Kind = "PreTool"
	// KindPostTool is the normalized category for successful post-tool events.
	KindPostTool Kind = "PostTool"
	// KindPostToolFailure is the normalized category for failed post-tool events.
	KindPostToolFailure Kind = "PostToolFailure"
	// KindPermissionRequest is the normalized category for permission request events.
	KindPermissionRequest Kind = "PermissionRequest"
	// KindSubagentStart is the normalized category for subagent start events.
	KindSubagentStart Kind = "SubagentStart"
	// KindSubagentStop is the normalized category for subagent stop events.
	KindSubagentStop Kind = "SubagentStop"
	// KindStop is the normalized category for agent stop events.
	KindStop Kind = "Stop"
	// KindPreCompact is the normalized category for pre-compaction events.
	KindPreCompact Kind = "PreCompact"
	// KindNotification is the normalized category for notification events.
	KindNotification Kind = "Notification"
	// KindAgentError is the normalized category for agent runtime error events.
	KindAgentError Kind = "AgentError"
	// KindOther is the normalized category for events with no dedicated mapping.
	KindOther Kind = "Other"
)
