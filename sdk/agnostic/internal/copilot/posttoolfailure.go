package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPostToolFailure registers fn on the Copilot PostToolUseFailure chain.
func RegisterPostToolFailure(fn model.PostToolFailureHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks().PostToolUseFailure(func(ctx context.Context, hook run.Hook[sdkcopilot.PostToolUseFailure], native sdkcopilot.PostToolFailureResults) (sdkcopilot.PostToolFailureOutput, error) {
		out, err := fn(ctx, model.NewPostToolFailureHook(hook.Invocation(), mapPostToolUseFailure(hook.Event)), newPostToolFailureResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPostToolFailure(out)
		if !ok {
			return nil, fmt.Errorf("copilot: PostToolFailure result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPostToolUseFailure(e sdkcopilot.PostToolUseFailure) *model.PostToolFailureEvent {
	ev := &model.PostToolFailureEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     model.NewToolCall(e.NativeToolName(), e.Input().Raw(), ""),
	}
	if msg := e.ErrorMessage(); msg != "" {
		ev.Result = &model.ToolResult{Error: msg}
	}
	return ev
}

func newPostToolFailureResults(native sdkcopilot.PostToolFailureResults) model.PostToolFailureResults {
	return postToolFailureResults{native: native}
}

func unwrapPostToolFailure(r model.PostToolFailureResult) (sdkcopilot.PostToolFailureOutput, bool) {
	out, ok := r.(postToolFailureResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolFailureResults struct {
	native sdkcopilot.PostToolFailureResults
}

// Context returns a context-injection result.
func (w postToolFailureResults) Context(text string) model.PostToolFailureResult {
	return postToolFailureResult{native: w.native.Context(text)}
}

type postToolFailureResult struct {
	native sdkcopilot.PostToolFailureOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolFailureResult) IsZero() bool { return r.native.IsZero() }
