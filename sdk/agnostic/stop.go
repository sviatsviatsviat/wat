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

// StopResult is the portable hook response for Stop and SubagentStop events.
// Construct only via StopResults (FollowUp).
// A nil value is a no-op.
type StopResult interface {
	isStopResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// StopResults is the hook-scoped response builder supplied to OnStop and OnSubagentStop handlers by registration.
type StopResults interface {
	// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
	FollowUp(text string) StopResult
	isStopResults()
}

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, hook StopHook, results StopResults) (StopResult, error)

// OnStop registers a handler for Stop events across all agents.
func OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().Stop(adaptClaudeStop(fn))
	sdkcopilot.Adapter().AgentStop(adaptCopilotStop(fn))
	sdkcursor.Adapter().Stop(adaptCursorStop(fn))
	return &Chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().SubagentStop(adaptClaudeSubagentStop(fn))
	sdkcopilot.Adapter().SubagentStop(adaptCopilotSubagentStop(fn))
	sdkcursor.Adapter().SubagentStop(adaptCursorSubagentStop(fn))
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

func adaptClaudeStop(fn StopHandler) func(context.Context, sdkclaude.Hook[sdkclaude.Stop], sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.Stop], native sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
		typed, err := StopEventFrom(agclaude.MapStop(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, stopHook(hook.Invocation(), typed), claudeStopResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudeStopResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: Stop result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

func adaptClaudeSubagentStop(fn StopHandler) func(context.Context, sdkclaude.Hook[sdkclaude.SubagentStop], sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SubagentStop], native sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
		typed, err := StopEventFrom(agclaude.MapSubagentStop(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, stopHook(hook.Invocation(), typed), claudeStopResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudeStopResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: Stop result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type claudeStopResults struct {
	native sdkclaude.StopResults
}

func (claudeStopResults) isStopResults() {}

// FollowUp returns a stop follow-up result.
func (w claudeStopResults) FollowUp(text string) StopResult {
	return claudeStopResult{native: w.native.FollowUp(text)}
}

type claudeStopResult struct {
	native sdkclaude.StopOutput
}

func (claudeStopResult) isStopResult() {}

// IsZero reports whether the result carries no instruction.
func (r claudeStopResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

func adaptCopilotStop(fn StopHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.AgentStop], sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.AgentStop], native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		typed, err := StopEventFrom(agcopilot.MapAgentStop(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, stopHook(hook.Invocation(), typed), copilotStopResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotStopResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: Stop result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

func adaptCopilotSubagentStop(fn StopHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.SubagentStop], sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SubagentStop], native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		typed, err := StopEventFrom(agcopilot.MapSubagentStop(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, stopHook(hook.Invocation(), typed), copilotStopResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotStopResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: Stop result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type copilotStopResults struct {
	native sdkcopilot.StopResults
}

func (copilotStopResults) isStopResults() {}

// FollowUp returns a stop follow-up result.
func (w copilotStopResults) FollowUp(text string) StopResult {
	return copilotStopResult{native: w.native.FollowUp(text)}
}

type copilotStopResult struct {
	native sdkcopilot.StopOutput
}

func (copilotStopResult) isStopResult() {}

// IsZero reports whether the result carries no instruction.
func (r copilotStopResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

func adaptCursorStop(fn StopHandler) func(context.Context, sdkcursor.Hook[sdkcursor.Stop], sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.Stop], native sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
		return callCursorStop(ctx, hook.Invocation(), agcursor.MapStop(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorSubagentStop(fn StopHandler) func(context.Context, sdkcursor.Hook[sdkcursor.SubagentStop], sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SubagentStop], native sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
		return callCursorStop(ctx, hook.Invocation(), agcursor.MapSubagentStop(hook.Event, hook.Raw()), native, fn)
	}
}

func callCursorStop(ctx context.Context, inv run.Invocation, ev *model.Event, native sdkcursor.StopResults, fn StopHandler) (sdkcursor.StopOutput, error) {
	typed, err := StopEventFrom(ev)
	if err != nil {
		return nil, err
	}
	out, err := fn(ctx, stopHook(inv, typed), cursorStopResults{native: native})
	if err != nil || out == nil {
		return nil, err
	}
	r, ok := out.(cursorStopResult)
	if !ok {
		return nil, fmt.Errorf("agnostic: Stop result must come from the injected Results builder")
	}
	return r.native, nil
}

type cursorStopResults struct {
	native sdkcursor.StopResults
}

func (cursorStopResults) isStopResults() {}

// FollowUp returns a stop follow-up result.
func (w cursorStopResults) FollowUp(text string) StopResult {
	return cursorStopResult{native: w.native.FollowUp(text)}
}

type cursorStopResult struct {
	native sdkcursor.StopOutput
}

func (cursorStopResult) isStopResult() {}

// IsZero reports whether the result carries no instruction.
func (r cursorStopResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }
