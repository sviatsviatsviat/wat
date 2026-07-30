package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterStop registers fn on the Copilot AgentStop chain for agent-scoped
// Stop payloads only (not when agent_name / agent_display_name is set).
func RegisterStop(fn model.StopHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcopilot.UseHooks().AgentStop(func(ctx context.Context, hook sdkcopilot.AgentStop, native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		if hook.IsSubagent() {
			return nil, nil
		}
		return callStop(ctx, mapAgentStop(hook), native, fn)
	})
}

// RegisterSubagentStop registers fn for explicit SubagentStop wire events and for
// Stop payloads scoped to a subagent (agent_name / agent_display_name set).
func RegisterSubagentStop(fn model.StopHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcopilot.UseHooks().SubagentStop(func(ctx context.Context, hook sdkcopilot.SubagentStop, native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		return callStop(ctx, mapSubagentStop(hook), native, fn)
	}).AgentStop(func(ctx context.Context, hook sdkcopilot.AgentStop, native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		if !hook.IsSubagent() {
			return nil, nil
		}
		return callStop(ctx, mapAgentStopAsSubagent(hook), native, fn)
	})
}

func callStop(ctx context.Context, ev *model.StopEvent, native sdkcopilot.StopResults, fn model.StopHandler) (sdkcopilot.StopOutput, error) {
	out, err := fn(ctx, *ev, newStopResults(native))
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapStop(out)
	if !ok {
		return nil, fmt.Errorf("copilot: Stop result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapAgentStop(e sdkcopilot.AgentStop) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Turn: &model.TurnEnd{
			Status:         e.Reason(),
			StopHookActive: e.StopHookActive,
		},
	}
}

func mapAgentStopAsSubagent(e sdkcopilot.AgentStop) *model.StopEvent {
	typ := e.Name()
	if typ == "" {
		typ = e.DisplayName()
	}
	return &model.StopEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{
			Type:   typ,
			Status: e.Reason(),
		},
		Turn: &model.TurnEnd{
			Status:         e.Reason(),
			StopHookActive: e.StopHookActive,
		},
	}
}

func mapSubagentStop(e sdkcopilot.SubagentStop) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{
			Type:   e.Name(),
			Status: e.Reason(),
		},
		Turn: &model.TurnEnd{Status: e.Reason()},
	}
}

func newStopResults(native sdkcopilot.StopResults) model.StopResults {
	return stopResults{native: native}
}

func unwrapStop(r model.StopResult) (sdkcopilot.StopOutput, bool) {
	out, ok := r.(stopResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type stopResults struct {
	native sdkcopilot.StopResults
}

// FollowUp returns a stop follow-up result.
func (w stopResults) FollowUp(text string) model.StopResult {
	return stopResult{native: w.native.FollowUp(text)}
}

type stopResult struct {
	native sdkcopilot.StopOutput
}

// IsZero reports whether the result carries no instruction.
func (r stopResult) IsZero() bool { return r.native.IsZero() }
