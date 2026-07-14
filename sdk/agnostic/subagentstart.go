package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent struct {
	Envelope
	Subagent *model.Subagent
}

// SubagentStartEventFrom maps a decoded Event to SubagentStartEvent.
func SubagentStartEventFrom(ev *model.Event) (SubagentStartEvent, error) {
	if ev == nil {
		return SubagentStartEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindSubagentStart {
		return SubagentStartEvent{}, fmt.Errorf("agnostic: expected SubagentStart kind, got %s", ev.Kind)
	}
	return SubagentStartEvent{Envelope: envelopeFrom(ev), Subagent: ev.Subagent}, nil
}

// SubagentStartHook is the handler context for portable SubagentStart events.
type SubagentStartHook struct {
	SubagentStartEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SubagentStartHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SubagentStartHook) Raw() json.RawMessage { return h.SubagentStartEvent.Raw }

func subagentStartHook(ctx run.Invocation, ev SubagentStartEvent) SubagentStartHook {
	return SubagentStartHook{SubagentStartEvent: ev, inv: ctx}
}

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler func(ctx context.Context, hook SubagentStartHook) error

// OnSubagentStart registers an observe-only handler for SubagentStart events.
func OnSubagentStart(fn SubagentStartHandler) *Chain {
	registerObserveHandler(model.KindSubagentStart, func(ctx context.Context, ev *model.Event) error {
		typed, err := SubagentStartEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, subagentStartHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnSubagentStart registers another observe-only SubagentStart handler on the chain.
func (c *Chain) OnSubagentStart(fn SubagentStartHandler) *Chain {
	return OnSubagentStart(fn)
}
