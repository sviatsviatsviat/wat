package cursor

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterStop registers fn on the Cursor Stop chain.
func RegisterStop(fn model.StopHandler) {
	if fn == nil {
		return
	}
	sdkcursor.OnStop(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.Stop], native sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapStop(hook.Event, hook.Raw()), native, fn)
	})
}

// RegisterSubagentStop registers fn on the Cursor SubagentStop chain.
func RegisterSubagentStop(fn model.StopHandler) {
	if fn == nil {
		return
	}
	sdkcursor.OnSubagentStop(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SubagentStop], native sdkcursor.StopResults) (sdkcursor.StopOutput, error) {
		return callStop(ctx, hook.Invocation(), mapSubagentStop(hook.Event, hook.Raw()), native, fn)
	})
}

func callStop(ctx context.Context, inv run.Invocation, ev *model.StopEvent, native sdkcursor.StopResults, fn model.StopHandler) (sdkcursor.StopOutput, error) {
	out, err := fn(ctx, model.NewStopHook(inv, ev), newStopResults(native))
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapStop(out)
	if !ok {
		return nil, fmt.Errorf("cursor: Stop result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapStop(e sdkcursor.Stop, raw []byte) *model.StopEvent {
	return &model.StopEvent{
		Envelope: envelope(e, raw),
		Turn:     &model.TurnEnd{Status: e.Status, LoopCount: e.LoopCount},
	}
}

func mapSubagentStop(e sdkcursor.SubagentStop, raw []byte) *model.StopEvent {
	tp := ""
	if e.AgentTranscriptPath != nil {
		tp = *e.AgentTranscriptPath
	}
	return &model.StopEvent{
		Envelope: envelope(e, raw),
		Subagent: &model.Subagent{
			ID:             e.SubagentID,
			Type:           e.SubagentType,
			Task:           e.Task,
			Summary:        e.Summary,
			Status:         e.Status,
			TranscriptPath: tp,
			LoopCount:      e.LoopCount,
		},
	}
}

func newStopResults(native sdkcursor.StopResults) model.StopResults {
	return stopResults{native: native}
}

func unwrapStop(r model.StopResult) (sdkcursor.StopOutput, bool) {
	out, ok := r.(stopResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type stopResults struct {
	native sdkcursor.StopResults
}

// FollowUp returns a stop follow-up result.
func (w stopResults) FollowUp(text string) model.StopResult {
	return stopResult{native: w.native.FollowUp(text)}
}

type stopResult struct {
	native sdkcursor.StopOutput
}

// IsZero reports whether the result carries no instruction.
func (r stopResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }
