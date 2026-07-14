package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hook is the handler context for a typed Copilot hook event.
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
	// PreToolUseHook is the handler context for preToolUse events.
	PreToolUseHook = Hook[PreToolUse]
	// PostToolUseHook is the handler context for postToolUse events.
	PostToolUseHook = Hook[PostToolUse]
	// PostToolUseFailureHook is the handler context for postToolUseFailure events.
	PostToolUseFailureHook = Hook[PostToolUseFailure]
	// AgentStopHook is the handler context for agentStop events.
	AgentStopHook = Hook[AgentStop]
	// SubagentStopHook is the handler context for subagentStop events.
	SubagentStopHook = Hook[SubagentStop]
	// PermissionRequestHook is the handler context for permissionRequest events.
	PermissionRequestHook = Hook[PermissionRequest]
	// SessionStartHook is the handler context for sessionStart events.
	SessionStartHook = Hook[SessionStart]
	// SessionEndHook is the handler context for sessionEnd events.
	SessionEndHook = Hook[SessionEnd]
	// UserPromptSubmittedHook is the handler context for userPromptSubmitted events.
	UserPromptSubmittedHook = Hook[UserPromptSubmitted]
	// SubagentStartHook is the handler context for subagentStart events.
	SubagentStartHook = Hook[SubagentStart]
	// PreCompactHook is the handler context for preCompact events.
	PreCompactHook = Hook[PreCompact]
	// NotificationHook is the handler context for notification events.
	NotificationHook = Hook[Notification]
	// ErrorOccurredHook is the handler context for errorOccurred events.
	ErrorOccurredHook = Hook[ErrorOccurred]
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
