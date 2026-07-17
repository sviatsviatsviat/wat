package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
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
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().SubagentStart(adaptClaudeSubagentStart(fn))
	sdkcopilot.Adapter().SubagentStart(adaptCopilotSubagentStart(fn))
	sdkcursor.Adapter().SubagentStart(adaptCursorSubagentStart(fn))
	return &Chain{}
}

// OnSubagentStart registers another observe-only SubagentStart handler on the chain.
func (c *Chain) OnSubagentStart(fn SubagentStartHandler) *Chain {
	return OnSubagentStart(fn)
}

func adaptClaudeSubagentStart(fn SubagentStartHandler) func(context.Context, sdkclaude.Hook[sdkclaude.SubagentStart], sdkclaude.SubagentStartResults) (sdkclaude.CommonOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SubagentStart], _ sdkclaude.SubagentStartResults) (sdkclaude.CommonOutput, error) {
		typed, err := SubagentStartEventFrom(agclaude.MapSubagentStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, subagentStartHook(hook.Invocation(), typed))
	}
}

func adaptCopilotSubagentStart(fn SubagentStartHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.SubagentStart], sdkcopilot.SubagentStartResults) (sdkcopilot.SubagentStartOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SubagentStart], _ sdkcopilot.SubagentStartResults) (sdkcopilot.SubagentStartOutput, error) {
		typed, err := SubagentStartEventFrom(agcopilot.MapSubagentStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, subagentStartHook(hook.Invocation(), typed))
	}
}

func adaptCursorSubagentStart(fn SubagentStartHandler) func(context.Context, sdkcursor.Hook[sdkcursor.SubagentStart], sdkcursor.SubagentStartResults) (sdkcursor.PermissionOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SubagentStart], _ sdkcursor.SubagentStartResults) (sdkcursor.PermissionOutput, error) {
		typed, err := SubagentStartEventFrom(agcursor.MapSubagentStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, subagentStartHook(hook.Invocation(), typed))
	}
}
