package model

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent struct {
	Envelope
	Subagent *Subagent
}

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

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler func(ctx context.Context, hook SubagentStartHook) error
