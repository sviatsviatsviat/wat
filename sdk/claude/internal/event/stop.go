package event

// StopActiveFields holds Claude stop-hook continuation fields shared by Stop
// and SubagentStop events.
type StopActiveFields struct {
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
	// BackgroundTasks lists in-flight background tasks when the task registry
	// is reachable (Claude Code v2.1.145+). Empty when nothing is in flight.
	// On SubagentStop, entries are scoped to the parent session.
	BackgroundTasks []BackgroundTask `json:"background_tasks,omitempty"`
	// SessionCrons lists session-scoped scheduled wakeups when the task
	// registry is reachable (Claude Code v2.1.145+). Empty when none are
	// scheduled. On SubagentStop, entries are scoped to the parent session.
	SessionCrons []SessionCron `json:"session_crons,omitempty"`
}

// BackgroundTask describes one in-flight background task on Stop and
// SubagentStop payloads (Claude Code v2.1.145+).
type BackgroundTask struct {
	// ID is the task identifier.
	ID string `json:"id"`
	// Type is a friendly task-type label such as shell, subagent, monitor,
	// workflow, teammate, "cloud session", or "MCP task".
	Type string `json:"type"`
	// Status is the current task status.
	Status string `json:"status"`
	// Description is a free-text description (may be truncated by the host).
	Description string `json:"description,omitempty"`
	// Command is the shell command line for shell tasks.
	Command string `json:"command,omitempty"`
	// AgentType is the subagent type name for subagent tasks.
	AgentType string `json:"agent_type,omitempty"`
	// Server is the MCP server name for monitor and MCP task entries.
	Server string `json:"server,omitempty"`
	// Tool is the MCP tool name for monitor and MCP task entries.
	Tool string `json:"tool,omitempty"`
	// Name is the workflow name for workflow tasks.
	Name string `json:"name,omitempty"`
}

// SessionCron describes one session-scoped scheduled wakeup on Stop and
// SubagentStop payloads (Claude Code v2.1.145+).
type SessionCron struct {
	// ID is the cron task identifier.
	ID string `json:"id"`
	// Schedule is the cron expression (for example "0 9 * * 1-5").
	Schedule string `json:"schedule"`
	// Recurring is false for one-shot wakeups and true for recurring matches.
	Recurring bool `json:"recurring"`
	// Prompt is the prompt submitted when the cron fires (may be truncated).
	Prompt string `json:"prompt,omitempty"`
}
