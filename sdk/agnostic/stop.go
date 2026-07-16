package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent struct {
	Envelope
	Turn     *model.TurnEnd
	Subagent *model.Subagent
}

// StopEventFrom maps a decoded Event to StopEvent.
func StopEventFrom(ev *model.Event) (StopEvent, error) {
	if ev == nil {
		return StopEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindStop && ev.Kind != model.KindSubagentStop {
		return StopEvent{}, fmt.Errorf("agnostic: expected Stop or SubagentStop kind, got %s", ev.Kind)
	}
	return StopEvent{Envelope: envelopeFrom(ev), Turn: ev.Turn, Subagent: ev.Subagent}, nil
}

// StopHook is the handler context for portable Stop and SubagentStop events.
type StopHook struct {
	StopEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h StopHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h StopHook) Raw() json.RawMessage { return h.StopEvent.Raw }

func stopHook(ctx run.Invocation, ev StopEvent) StopHook {
	return StopHook{StopEvent: ev, inv: ctx}
}

// StopResults is the hook-scoped response builder supplied to OnStop and OnSubagentStop handlers by registration.
type StopResults interface {
	// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
	FollowUp(text string) model.StopResult
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
func (stopResults) FollowUp(text string) model.StopResult { return model.StopFollowUp(text) }

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, hook StopHook, results StopResults) (model.StopResult, error)

// OnStop registers a handler for Stop events across all agents.
func OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindStop, func(ctx context.Context, ev *model.Event) (model.StopResult, error) {
		typed, err := StopEventFrom(ev)
		if err != nil {
			return nil, err
		}
		return fn(ctx, stopHook(run.InvocationFrom(ctx), typed), stopResults{})
	})
	return &Chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindSubagentStop, func(ctx context.Context, ev *model.Event) (model.StopResult, error) {
		typed, err := StopEventFrom(ev)
		if err != nil {
			return nil, err
		}
		return fn(ctx, stopHook(run.InvocationFrom(ctx), typed), stopResults{})
	})
	return &Chain{}
}

// OnStop registers another Stop handler on the chain.
func (c *Chain) OnStop(fn StopHandler) *Chain {
	return OnStop(fn)
}

// OnSubagentStop registers another SubagentStop handler on the chain.
func (c *Chain) OnSubagentStop(fn StopHandler) *Chain {
	return OnSubagentStop(fn)
}
