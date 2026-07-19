package claude

// PermissionMode is the Claude Code permission mode on hook events.
type PermissionMode string

const (
	// PermissionDefault is the default permission mode.
	PermissionDefault PermissionMode = "default"
	// PermissionPlan is the plan permission mode.
	PermissionPlan PermissionMode = "plan"
	// PermissionAcceptEdits auto-accepts edit tools.
	PermissionAcceptEdits PermissionMode = "acceptEdits"
	// PermissionAuto is the auto permission mode.
	PermissionAuto PermissionMode = "auto"
	// PermissionDontAsk suppresses permission prompts.
	PermissionDontAsk PermissionMode = "dontAsk"
	// PermissionBypassPermissions bypasses permission checks.
	PermissionBypassPermissions PermissionMode = "bypassPermissions"
)

// EffortLevel is the effort level on hook events.
type EffortLevel string

const (
	// EffortLow is the low effort level.
	EffortLow EffortLevel = "low"
	// EffortMedium is the medium effort level.
	EffortMedium EffortLevel = "medium"
	// EffortHigh is the high effort level.
	EffortHigh EffortLevel = "high"
	// EffortXHigh is the extra-high effort level.
	EffortXHigh EffortLevel = "xhigh"
	// EffortMax is the maximum effort level.
	EffortMax EffortLevel = "max"
)

// Effort carries effort metadata on hook events.
type Effort struct {
	// Level is the effort level (low, medium, high, xhigh, max).
	Level EffortLevel `json:"level"`
}

// Envelope holds fields shared by every Claude Code hook event payload.
type Envelope struct {
	// SessionID is the Claude session identifier.
	SessionID string `json:"session_id"`
	// PromptID is the prompt identifier when provided.
	PromptID string `json:"prompt_id,omitempty"`
	// TranscriptPath is the conversation transcript path.
	TranscriptPath string `json:"transcript_path"`
	// Cwd is the working directory.
	Cwd string `json:"cwd"`
	// PermissionMode is the active permission mode.
	PermissionMode PermissionMode `json:"permission_mode,omitempty"`
	// Effort is the effort metadata when provided.
	Effort *Effort `json:"effort,omitempty"`
	// HookEventName is the native hook event name.
	HookEventName string `json:"hook_event_name"`
	// AgentID is the subagent identifier inside subagent events.
	AgentID string `json:"agent_id,omitempty"`
	// AgentType is the subagent type inside subagent events.
	AgentType string `json:"agent_type,omitempty"`
}
