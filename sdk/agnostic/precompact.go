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

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	// Compact holds compaction details for this pre-compact event.
	Compact *compactInfo
}

// PreCompactEventFrom maps a decoded Event to PreCompactEvent.
func PreCompactEventFrom(ev *Event) (PreCompactEvent, error) {
	if ev == nil {
		return PreCompactEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPreCompact {
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
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().PreCompact(adaptClaudePreCompact(fn))
	sdkcopilot.Adapter().PreCompact(adaptCopilotPreCompact(fn))
	sdkcursor.Adapter().PreCompact(adaptCursorPreCompact(fn))
	return &Chain{}
}

// OnPreCompact registers another observe-only PreCompact handler on the chain.
func (c *Chain) OnPreCompact(fn PreCompactHandler) *Chain {
	return OnPreCompact(fn)
}

func adaptClaudePreCompact(fn PreCompactHandler) func(context.Context, sdkclaude.Hook[sdkclaude.PreCompact], sdkclaude.PreCompactResults) (sdkclaude.CommonOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PreCompact], _ sdkclaude.PreCompactResults) (sdkclaude.CommonOutput, error) {
		typed, err := PreCompactEventFrom(mapClaudePreCompact(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, preCompactHook(hook.Invocation(), typed))
	}
}

func adaptCopilotPreCompact(fn PreCompactHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.PreCompact]) error {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PreCompact]) error {
		typed, err := PreCompactEventFrom(mapCopilotPreCompact(hook.Event, hook.Raw()))
		if err != nil {
			return err
		}
		return fn(ctx, preCompactHook(hook.Invocation(), typed))
	}
}

func adaptCursorPreCompact(fn PreCompactHandler) func(context.Context, sdkcursor.Hook[sdkcursor.PreCompact], sdkcursor.PreCompactResults) (sdkcursor.PreCompactOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PreCompact], _ sdkcursor.PreCompactResults) (sdkcursor.PreCompactOutput, error) {
		typed, err := PreCompactEventFrom(mapCursorPreCompact(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, preCompactHook(hook.Invocation(), typed))
	}
}

func mapClaudePreCompact(e sdkclaude.PreCompact, raw []byte) *Event {
	ev := claudeEvent(e, raw, KindPreCompact)
	ev.Compact = &compactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
	return ev
}

func mapCopilotPreCompact(e sdkcopilot.PreCompact, raw []byte) *Event {
	ev := copilotEvent(e, raw, KindPreCompact)
	ev.Compact = &compactInfo{
		Trigger:            e.Trigger,
		CustomInstructions: e.Instructions(),
	}
	return ev
}

func mapCursorPreCompact(e sdkcursor.PreCompact, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindPreCompact)
	ev.Compact = &compactInfo{Trigger: e.Trigger}
	return ev
}
