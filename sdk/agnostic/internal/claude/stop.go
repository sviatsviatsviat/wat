package claude

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterStop registers fn on the Claude Stop chain.
func RegisterStop(r *run.Registry, fn model.StopHandler) {
	if fn == nil {
		return
	}
	sdkclaude.UseHooks(r).Stop(func(ctx context.Context, hook run.Hook[sdkclaude.Stop], native sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapStop(hook.Event), native, fn)
	})
}

// RegisterSubagentStop registers fn on the Claude SubagentStop chain.
func RegisterSubagentStop(r *run.Registry, fn model.StopHandler) {
	if fn == nil {
		return
	}
	sdkclaude.UseHooks(r).SubagentStop(func(ctx context.Context, hook run.Hook[sdkclaude.SubagentStop], native sdkclaude.StopResults) (sdkclaude.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapSubagentStop(hook.Event), native, fn)
	})
}

func callStop(ctx context.Context, inv run.Invocation, ev *model.StopEvent, native sdkclaude.StopResults, fn model.StopHandler) (sdkclaude.StopOutput, error) {
	out, err := fn(ctx, model.NewStopHook(inv, ev), newStopResults(native))
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapStop(out)
	if !ok {
		return nil, fmt.Errorf("claude: Stop result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapStop(e sdkclaude.Stop) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Turn:     &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage},
	}
}

func mapSubagentStop(e sdkclaude.SubagentStop) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage},
		Turn:     &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage},
	}
}

func newStopResults(native sdkclaude.StopResults) model.StopResults {
	return stopResults{native: native}
}

func unwrapStop(r model.StopResult) (sdkclaude.StopOutput, bool) {
	out, ok := r.(stopResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type stopResults struct {
	native sdkclaude.StopResults
}

// FollowUp returns a stop follow-up result.
func (w stopResults) FollowUp(text string) model.StopResult {
	return stopResult{native: w.native.FollowUp(text)}
}

type stopResult struct {
	native sdkclaude.StopOutput
}

// IsZero reports whether the result carries no instruction.
func (r stopResult) IsZero() bool { return r.native.IsZero() }
