package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *model.Lifecycle
}

// SessionEndEventFrom maps a decoded Event to SessionEndEvent.
func SessionEndEventFrom(ev *model.Event) (SessionEndEvent, error) {
	if ev == nil {
		return SessionEndEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindSessionEnd {
		return SessionEndEvent{}, fmt.Errorf("agnostic: expected SessionEnd kind, got %s", ev.Kind)
	}
	return SessionEndEvent{Envelope: envelopeFrom(ev), Life: ev.Life}, nil
}

// SessionEndHook is the handler context for portable SessionEnd events.
type SessionEndHook struct {
	SessionEndEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionEndHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SessionEndHook) Raw() json.RawMessage { return h.SessionEndEvent.Raw }

func sessionEndHook(ctx run.Invocation, ev SessionEndEvent) SessionEndHook {
	return SessionEndHook{SessionEndEvent: ev, inv: ctx}
}

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler func(ctx context.Context, hook SessionEndHook) error

// OnSessionEnd registers an observe-only handler for SessionEnd events.
func OnSessionEnd(fn SessionEndHandler) *Chain {
	registerObserveHandler(model.KindSessionEnd, func(ctx context.Context, ev *model.Event) error {
		typed, err := SessionEndEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, sessionEndHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnSessionEnd registers another observe-only SessionEnd handler on the chain.
func (c *Chain) OnSessionEnd(fn SessionEndHandler) *Chain {
	return OnSessionEnd(fn)
}
