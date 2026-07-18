package model

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	Compact *CompactInfo
}

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

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler func(ctx context.Context, hook PreCompactHook) error
