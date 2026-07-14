package agnostic

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolHook is the handler context for portable PreTool events.
type PreToolHook struct {
	PreToolEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreToolHook) Invocation() run.Invocation { return h.inv }

// PostToolHook is the handler context for portable PostTool events.
type PostToolHook struct {
	PostToolEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolHook) Invocation() run.Invocation { return h.inv }

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook struct {
	PostToolFailureEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolFailureHook) Invocation() run.Invocation { return h.inv }

// StopHook is the handler context for portable Stop and SubagentStop events.
type StopHook struct {
	StopEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h StopHook) Invocation() run.Invocation { return h.inv }

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook struct {
	SessionStartEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionStartHook) Invocation() run.Invocation { return h.inv }

// SessionEndHook is the handler context for portable SessionEnd events.
type SessionEndHook struct {
	SessionEndEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionEndHook) Invocation() run.Invocation { return h.inv }

// UserPromptHook is the handler context for portable UserPrompt events.
type UserPromptHook struct {
	UserPromptEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h UserPromptHook) Invocation() run.Invocation { return h.inv }

// PreCompactHook is the handler context for portable PreCompact events.
type PreCompactHook struct {
	PreCompactEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreCompactHook) Invocation() run.Invocation { return h.inv }

// SubagentStartHook is the handler context for portable SubagentStart events.
type SubagentStartHook struct {
	SubagentStartEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SubagentStartHook) Invocation() run.Invocation { return h.inv }

// AnyHook is the handler context for catch-all OnAny handlers.
type AnyHook struct {
	AnyEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h AnyHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreToolHook) Raw() json.RawMessage { return h.PreToolEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h PostToolHook) Raw() json.RawMessage { return h.PostToolEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h PostToolFailureHook) Raw() json.RawMessage { return h.PostToolFailureEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h StopHook) Raw() json.RawMessage { return h.StopEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h SessionStartHook) Raw() json.RawMessage { return h.SessionStartEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h SessionEndHook) Raw() json.RawMessage { return h.SessionEndEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h UserPromptHook) Raw() json.RawMessage { return h.UserPromptEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h PreCompactHook) Raw() json.RawMessage { return h.PreCompactEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h SubagentStartHook) Raw() json.RawMessage { return h.SubagentStartEvent.Raw }

// Raw returns the untouched native JSON payload.
func (h AnyHook) Raw() json.RawMessage { return h.AnyEvent.Raw }

func preToolHook(ctx run.Invocation, ev PreToolEvent) PreToolHook {
	return PreToolHook{PreToolEvent: ev, inv: ctx}
}
func postToolHook(ctx run.Invocation, ev PostToolEvent) PostToolHook {
	return PostToolHook{PostToolEvent: ev, inv: ctx}
}
func postToolFailureHook(ctx run.Invocation, ev PostToolFailureEvent) PostToolFailureHook {
	return PostToolFailureHook{PostToolFailureEvent: ev, inv: ctx}
}
func stopHook(ctx run.Invocation, ev StopEvent) StopHook {
	return StopHook{StopEvent: ev, inv: ctx}
}
func sessionStartHook(ctx run.Invocation, ev SessionStartEvent) SessionStartHook {
	return SessionStartHook{SessionStartEvent: ev, inv: ctx}
}
func sessionEndHook(ctx run.Invocation, ev SessionEndEvent) SessionEndHook {
	return SessionEndHook{SessionEndEvent: ev, inv: ctx}
}
func userPromptHook(ctx run.Invocation, ev UserPromptEvent) UserPromptHook {
	return UserPromptHook{UserPromptEvent: ev, inv: ctx}
}
func preCompactHook(ctx run.Invocation, ev PreCompactEvent) PreCompactHook {
	return PreCompactHook{PreCompactEvent: ev, inv: ctx}
}
func subagentStartHook(ctx run.Invocation, ev SubagentStartEvent) SubagentStartHook {
	return SubagentStartHook{SubagentStartEvent: ev, inv: ctx}
}
func anyHook(ctx run.Invocation, ev AnyEvent) AnyHook {
	return AnyHook{AnyEvent: ev, inv: ctx}
}
