package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *model.ToolCall
	Result *model.ToolResult
}

// PostToolFailureEventFrom maps a decoded Event to PostToolFailureEvent.
func PostToolFailureEventFrom(ev *model.Event) (PostToolFailureEvent, error) {
	if ev == nil {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindPostToolFailure {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: expected PostToolFailure kind, got %s", ev.Kind)
	}
	return PostToolFailureEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool, Result: ev.Result}, nil
}

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook struct {
	PostToolFailureEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolFailureHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolFailureHook) Raw() json.RawMessage { return h.PostToolFailureEvent.Raw }

func postToolFailureHook(ctx run.Invocation, ev PostToolFailureEvent) PostToolFailureHook {
	return PostToolFailureHook{PostToolFailureEvent: ev, inv: ctx}
}

// PostToolFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolFailureResults interface {
	// Context returns recovery guidance for PostToolFailure events.
	Context(text string) model.PostToolFailureResult
	isPostToolFailureResults()
}

type postToolFailureResults struct{}

func (postToolFailureResults) isPostToolFailureResults() {}

// Context returns recovery guidance for PostToolFailure events.
func (postToolFailureResults) Context(text string) model.PostToolFailureResult {
	return model.PostToolFailureContext(text)
}

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, hook PostToolFailureHook, results PostToolFailureResults) (model.PostToolFailureResult, error)

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindPostToolFailure, func(ctx context.Context, ev *model.Event) (model.PostToolFailureResult, error) {
		typed, err := PostToolFailureEventFrom(ev)
		if err != nil {
			return nil, err
		}
		return fn(ctx, postToolFailureHook(run.InvocationFrom(ctx), typed), postToolFailureResults{})
	})
	return &Chain{}
}

// OnPostToolFailure registers another PostToolFailure handler on the chain.
func (c *Chain) OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	return OnPostToolFailure(fn)
}
