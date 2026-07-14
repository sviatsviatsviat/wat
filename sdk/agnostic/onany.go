package agnostic

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AnyEvent is the catch-all normalized view for OnAny handlers.
type AnyEvent struct {
	Envelope
	Kind     model.Kind
	Prompt   string
	Tool     *model.ToolCall
	Result   *model.ToolResult
	Subagent *model.Subagent
	Turn     *model.TurnEnd
	Compact  *model.CompactInfo
	Note     *model.Note
	Life     *model.Lifecycle
}

// AnyEventFrom maps a decoded Event to AnyEvent.
func AnyEventFrom(ev *model.Event) AnyEvent {
	if ev == nil {
		return AnyEvent{}
	}
	return AnyEvent{
		Envelope: envelopeFrom(ev),
		Kind:     ev.Kind,
		Prompt:   ev.Prompt,
		Tool:     ev.Tool,
		Result:   ev.Result,
		Subagent: ev.Subagent,
		Turn:     ev.Turn,
		Compact:  ev.Compact,
		Note:     ev.Note,
		Life:     ev.Life,
	}
}

// AnyHook is the handler context for catch-all OnAny handlers.
type AnyHook struct {
	AnyEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h AnyHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h AnyHook) Raw() json.RawMessage { return h.AnyEvent.Raw }

func anyHook(ctx run.Invocation, ev AnyEvent) AnyHook {
	return AnyHook{AnyEvent: ev, inv: ctx}
}

// AnyHandler handles every event before kind-specific handlers with no hook response.
type AnyHandler func(ctx context.Context, hook AnyHook) error

// OnAny registers an observe-only handler invoked for every event before kind-specific handlers.
func OnAny(fn AnyHandler) *Chain {
	registerAny(fn)
	return &Chain{}
}

// OnAny registers another observe-only catch-all handler on the chain.
func (c *Chain) OnAny(fn AnyHandler) *Chain {
	return OnAny(fn)
}
