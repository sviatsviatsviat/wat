package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *lifecycle
}

// SessionEndEventFrom maps a decoded Event to SessionEndEvent.
func SessionEndEventFrom(ev *Event) (SessionEndEvent, error) {
	if ev == nil {
		return SessionEndEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindSessionEnd {
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
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().SessionEnd(adaptClaudeSessionEnd(fn))
	sdkcopilot.Adapter().SessionEnd(adaptCopilotSessionEnd(fn))
	sdkcursor.Adapter().SessionEnd(adaptCursorSessionEnd(fn))
	return &Chain{}
}

// OnSessionEnd registers another observe-only SessionEnd handler on the chain.
func (c *Chain) OnSessionEnd(fn SessionEndHandler) *Chain {
	return OnSessionEnd(fn)
}

func adaptClaudeSessionEnd(fn SessionEndHandler) func(context.Context, sdkclaude.Hook[sdkclaude.SessionEnd]) error {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SessionEnd]) error {
		typed, err := SessionEndEventFrom(mapClaudeSessionEnd(hook.Event, hook.Raw()))
		if err != nil {
			return err
		}
		return fn(ctx, sessionEndHook(hook.Invocation(), typed))
	}
}

func adaptCopilotSessionEnd(fn SessionEndHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.SessionEnd]) error {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SessionEnd]) error {
		typed, err := SessionEndEventFrom(mapCopilotSessionEnd(hook.Event, hook.Raw()))
		if err != nil {
			return err
		}
		return fn(ctx, sessionEndHook(hook.Invocation(), typed))
	}
}

func adaptCursorSessionEnd(fn SessionEndHandler) func(context.Context, sdkcursor.Hook[sdkcursor.SessionEnd]) error {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SessionEnd]) error {
		typed, err := SessionEndEventFrom(mapCursorSessionEnd(hook.Event, hook.Raw()))
		if err != nil {
			return err
		}
		return fn(ctx, sessionEndHook(hook.Invocation(), typed))
	}
}

func mapClaudeSessionEnd(e sdkclaude.SessionEnd, raw []byte) *Event {
	ev := claudeEvent(e, raw, KindSessionEnd)
	ev.Life = &lifecycle{Reason: e.Reason}
	return ev
}

func mapCopilotSessionEnd(e sdkcopilot.SessionEnd, raw []byte) *Event {
	ev := copilotEvent(e, raw, KindSessionEnd)
	ev.Life = &lifecycle{Reason: e.Reason}
	return ev
}

func mapCursorSessionEnd(e sdkcursor.SessionEnd, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindSessionEnd)
	ev.Life = &lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
	return ev
}
