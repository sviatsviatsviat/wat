package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolHook is the handler context for portable PostTool events.
type PostToolHook struct {
	PostToolEvent
	inv run.Invocation
}

// NewPostToolHook wraps ev with serve-time invocation settings.
func NewPostToolHook(inv run.Invocation, ev *PostToolEvent) PostToolHook {
	h := PostToolHook{inv: inv}
	if ev != nil {
		h.PostToolEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolHook) Raw() json.RawMessage { return h.PostToolEvent.Raw }

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, hook PostToolHook, results PostToolResults) (PostToolResult, error)

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook struct {
	PostToolFailureEvent
	inv run.Invocation
}

// NewPostToolFailureHook wraps ev with serve-time invocation settings.
func NewPostToolFailureHook(inv run.Invocation, ev *PostToolFailureEvent) PostToolFailureHook {
	h := PostToolFailureHook{inv: inv}
	if ev != nil {
		h.PostToolFailureEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolFailureHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolFailureHook) Raw() json.RawMessage { return h.PostToolFailureEvent.Raw }

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, hook PostToolFailureHook, results PostToolFailureResults) (PostToolFailureResult, error)

// PreToolHook is the handler context for portable PreTool events.
type PreToolHook struct {
	PreToolEvent
	inv run.Invocation
}

// NewPreToolHook wraps ev with serve-time invocation settings.
func NewPreToolHook(inv run.Invocation, ev *PreToolEvent) PreToolHook {
	h := PreToolHook{inv: inv}
	if ev != nil {
		h.PreToolEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreToolHook) Raw() json.RawMessage { return h.PreToolEvent.Raw }

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, hook PreToolHook, results PreToolResults) (PreToolResult, error)

// StopHook is the handler context for portable Stop and SubagentStop events.
type StopHook struct {
	StopEvent
	inv run.Invocation
}

// NewStopHook wraps ev with serve-time invocation settings.
func NewStopHook(inv run.Invocation, ev *StopEvent) StopHook {
	h := StopHook{inv: inv}
	if ev != nil {
		h.StopEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h StopHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h StopHook) Raw() json.RawMessage { return h.StopEvent.Raw }

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, hook StopHook, results StopResults) (StopResult, error)

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook struct {
	SessionStartEvent
	inv run.Invocation
}

// NewSessionStartHook wraps ev with serve-time invocation settings.
func NewSessionStartHook(inv run.Invocation, ev *SessionStartEvent) SessionStartHook {
	h := SessionStartHook{inv: inv}
	if ev != nil {
		h.SessionStartEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionStartHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SessionStartHook) Raw() json.RawMessage { return h.SessionStartEvent.Raw }

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, hook SessionStartHook, results SessionStartResults) (SessionStartResult, error)

// SessionEndHook is the handler context for portable SessionEnd events.
type SessionEndHook struct {
	SessionEndEvent
	inv run.Invocation
}

// NewSessionEndHook wraps ev with serve-time invocation settings.
func NewSessionEndHook(inv run.Invocation, ev *SessionEndEvent) SessionEndHook {
	h := SessionEndHook{inv: inv}
	if ev != nil {
		h.SessionEndEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionEndHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SessionEndHook) Raw() json.RawMessage { return h.SessionEndEvent.Raw }

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler func(ctx context.Context, hook SessionEndHook) error

// UserPromptHook is the handler context for portable UserPrompt events.
type UserPromptHook struct {
	UserPromptEvent
	inv run.Invocation
}

// NewUserPromptHook wraps ev with serve-time invocation settings.
func NewUserPromptHook(inv run.Invocation, ev *UserPromptEvent) UserPromptHook {
	h := UserPromptHook{inv: inv}
	if ev != nil {
		h.UserPromptEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h UserPromptHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h UserPromptHook) Raw() json.RawMessage { return h.UserPromptEvent.Raw }

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler func(ctx context.Context, hook UserPromptHook) error

// PreCompactHook is the handler context for portable PreCompact events.
type PreCompactHook struct {
	PreCompactEvent
	inv run.Invocation
}

// NewPreCompactHook wraps ev with serve-time invocation settings.
func NewPreCompactHook(inv run.Invocation, ev *PreCompactEvent) PreCompactHook {
	h := PreCompactHook{inv: inv}
	if ev != nil {
		h.PreCompactEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreCompactHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreCompactHook) Raw() json.RawMessage { return h.PreCompactEvent.Raw }

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler func(ctx context.Context, hook PreCompactHook) error

// SubagentStartHook is the handler context for portable SubagentStart events.
type SubagentStartHook struct {
	SubagentStartEvent
	inv run.Invocation
}

// NewSubagentStartHook wraps ev with serve-time invocation settings.
func NewSubagentStartHook(inv run.Invocation, ev *SubagentStartEvent) SubagentStartHook {
	h := SubagentStartHook{inv: inv}
	if ev != nil {
		h.SubagentStartEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h SubagentStartHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SubagentStartHook) Raw() json.RawMessage { return h.SubagentStartEvent.Raw }

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler func(ctx context.Context, hook SubagentStartHook) error
