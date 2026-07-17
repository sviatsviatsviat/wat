package model

import "encoding/json"

// Envelope carries shared metadata present on every normalized hook event.
type Envelope struct {
	// Agent is the dialect that emitted this hook event (e.g. "claude").
	Agent string
	// Name is the native event name as received.
	Name string
	// Session holds session_id, sessionId, or conversation_id from the native payload.
	Session string
	// Cwd is the working directory from the native payload.
	Cwd string
	// TranscriptPath is the conversation transcript path when provided.
	TranscriptPath string
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	Tool *ToolCall
}

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent struct {
	Envelope
	Turn     *TurnEnd
	Subagent *Subagent
}

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *Lifecycle
}

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent struct {
	Envelope
	Prompt string
}

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	Compact *CompactInfo
}

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent struct {
	Envelope
	Subagent *Subagent
}
