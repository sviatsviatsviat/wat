package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterStop registers fn on the Copilot AgentStop chain.
func RegisterStop(fn model.StopHandler) {
	if fn == nil {
		return
	}
	new(sdkcopilot.Chain).AgentStop(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.AgentStop], native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapAgentStop(hook.Event, hook.Raw()), native, fn)
	})
}

// RegisterSubagentStop registers fn on the Copilot SubagentStop chain.
func RegisterSubagentStop(fn model.StopHandler) {
	if fn == nil {
		return
	}
	new(sdkcopilot.Chain).SubagentStop(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SubagentStop], native sdkcopilot.StopResults) (sdkcopilot.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapSubagentStop(hook.Event, hook.Raw()), native, fn)
	})
}

func callStop(ctx context.Context, inv run.Invocation, ev *model.StopEvent, native sdkcopilot.StopResults, fn model.StopHandler) (sdkcopilot.StopOutput, error) {
	out, err := fn(ctx, model.NewStopHook(inv, ev), newStopResults(native))
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapStop(out)
	if !ok {
		return nil, fmt.Errorf("copilot: Stop result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapAgentStop(e sdkcopilot.AgentStop, raw []byte) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e, raw),
		Turn:     &model.TurnEnd{Status: e.Reason()},
	}
}

func mapSubagentStop(e sdkcopilot.SubagentStop, raw []byte) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e, raw),
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
func (r stopResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }
