package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	Compact *model.CompactInfo
}

// PreCompactEventFrom maps a decoded Event to PreCompactEvent.
func PreCompactEventFrom(ev *model.Event) (PreCompactEvent, error) {
	if ev == nil {
		return PreCompactEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindPreCompact {
		return PreCompactEvent{}, fmt.Errorf("agnostic: expected PreCompact kind, got %s", ev.Kind)
	}
	return PreCompactEvent{Envelope: envelopeFrom(ev), Compact: ev.Compact}, nil
}

// PreCompactHook is the handler context for portable PreCompact events.
type PreCompactHook struct {
	PreCompactEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreCompactHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreCompactHook) Raw() json.RawMessage { return h.PreCompactEvent.Raw }

func preCompactHook(ctx run.Invocation, ev PreCompactEvent) PreCompactHook {
	return PreCompactHook{PreCompactEvent: ev, inv: ctx}
}

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler func(ctx context.Context, hook PreCompactHook) error

// OnPreCompact registers an observe-only handler for PreCompact events.
func OnPreCompact(fn PreCompactHandler) *Chain {
	registerObserveHandler(model.KindPreCompact, func(ctx context.Context, ev *model.Event) error {
		typed, err := PreCompactEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, preCompactHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnPreCompact registers another observe-only PreCompact handler on the chain.
func (c *Chain) OnPreCompact(fn PreCompactHandler) *Chain {
	return OnPreCompact(fn)
}
