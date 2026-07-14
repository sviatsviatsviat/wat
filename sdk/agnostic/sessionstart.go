package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	Life *model.Lifecycle
}

// SessionStartEventFrom maps a decoded Event to SessionStartEvent.
func SessionStartEventFrom(ev *model.Event) (SessionStartEvent, error) {
	if ev == nil {
		return SessionStartEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindSessionStart {
		return SessionStartEvent{}, fmt.Errorf("agnostic: expected SessionStart kind, got %s", ev.Kind)
	}
	return SessionStartEvent{Envelope: envelopeFrom(ev), Life: ev.Life}, nil
}

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook struct {
	SessionStartEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionStartHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SessionStartHook) Raw() json.RawMessage { return h.SessionStartEvent.Raw }

func sessionStartHook(ctx run.Invocation, ev SessionStartEvent) SessionStartHook {
	return SessionStartHook{SessionStartEvent: ev, inv: ctx}
}

// SessionStartResults is the hook-scoped response builder supplied to OnSessionStart handlers by registration.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) model.SessionStartResult
	isSessionStartResults()
}

type sessionStartResults struct{}

func (sessionStartResults) isSessionStartResults() {}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) model.SessionStartResult {
	return model.SessionStartContext(text)
}

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, hook SessionStartHook, results SessionStartResults) (model.SessionStartResult, error)

// OnSessionStart registers a handler for SessionStart events across all agents.
func OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindSessionStart, func(ctx context.Context, ev *model.Event) (model.SessionStartResult, error) {
		typed, err := SessionStartEventFrom(ev)
		if err != nil {
			return model.SessionStartResult{}, err
		}
		return fn(ctx, sessionStartHook(run.InvocationFrom(ctx), typed), sessionStartResults{})
	})
	return &Chain{}
}

// OnSessionStart registers another SessionStart handler on the chain.
func (c *Chain) OnSessionStart(fn SessionStartHandler) *Chain {
	return OnSessionStart(fn)
}
