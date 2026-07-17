package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent struct {
	Envelope
	Turn     *TurnEnd
	Subagent *Subagent
}

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

// StopResult is the portable hook response for Stop and SubagentStop events.
// Construct only via StopResults (FollowUp).
// A nil value is a no-op.
type StopResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// StopResults is the hook-scoped response builder for Stop handlers.
type StopResults interface {
	// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
	FollowUp(text string) StopResult
}
