package model

import (
	"context"
)

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent struct {
	Envelope
	Turn     *TurnEnd
	Subagent *Subagent
}

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, event StopEvent, results StopResults) (StopResult, error)

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
