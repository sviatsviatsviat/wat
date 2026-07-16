package claude

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hook is the handler context for a typed Claude hook event.
type Hook[E Event] struct {
	// Event is the decoded native hook payload.
	Event E
	inv   run.Invocation
}

// NewHook wraps ev with serve-time invocation settings.
func NewHook[E Event](inv run.Invocation, ev E) Hook[E] {
	return Hook[E]{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h Hook[E]) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload when available.
func (h Hook[E]) Raw() json.RawMessage { return RawBytes(h.Event) }

type (
	// PreToolUseHook is the handler context for PreToolUse events.
	PreToolUseHook = Hook[PreToolUse]
	// PostToolUseHook is the handler context for PostToolUse events.
	PostToolUseHook = Hook[PostToolUse]
	// PostToolUseFailureHook is the handler context for PostToolUseFailure events.
	PostToolUseFailureHook = Hook[PostToolUseFailure]
	// PermissionRequestHook is the handler context for PermissionRequest events.
	PermissionRequestHook = Hook[PermissionRequest]
	// PermissionDeniedHook is the handler context for PermissionDenied events.
	PermissionDeniedHook = Hook[PermissionDenied]
	// UserPromptSubmitHook is the handler context for UserPromptSubmit events.
	UserPromptSubmitHook = Hook[UserPromptSubmit]
	// UserPromptExpansionHook is the handler context for UserPromptExpansion events.
	UserPromptExpansionHook = Hook[UserPromptExpansion]
	// StopHook is the handler context for Stop events.
	StopHook = Hook[Stop]
	// SubagentStopHook is the handler context for SubagentStop events.
	SubagentStopHook = Hook[SubagentStop]
	// SessionStartHook is the handler context for SessionStart events.
	SessionStartHook = Hook[SessionStart]
	// SubagentStartHook is the handler context for SubagentStart events.
	SubagentStartHook = Hook[SubagentStart]
	// NotificationHook is the handler context for Notification events.
	NotificationHook = Hook[Notification]
	// PreCompactHook is the handler context for PreCompact events.
	PreCompactHook = Hook[PreCompact]
	// SessionEndHook is the handler context for SessionEnd events.
	SessionEndHook = Hook[SessionEnd]
	// MessageDisplayHook is the handler context for MessageDisplay events.
	MessageDisplayHook = Hook[MessageDisplay]
	// TaskCreatedHook is the handler context for TaskCreated events.
	TaskCreatedHook = Hook[TaskCreated]
	// TaskCompletedHook is the handler context for TaskCompleted events.
	TaskCompletedHook = Hook[TaskCompleted]
	// WorktreeCreateHook is the handler context for WorktreeCreate events.
	WorktreeCreateHook = Hook[WorktreeCreate]
)

// AnyHook is the handler context for catch-all OnAny handlers.
type AnyHook struct {
	Event
	inv run.Invocation
}

// NewAnyHook wraps ev with serve-time invocation settings.
func NewAnyHook(inv run.Invocation, ev Event) AnyHook {
	return AnyHook{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h AnyHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload when available.
func (h AnyHook) Raw() json.RawMessage { return RawBytes(h.Event) }
