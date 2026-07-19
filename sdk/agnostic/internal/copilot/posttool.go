package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPostTool registers fn on the Copilot PostToolUse chain.
func RegisterPostTool(r *run.Registry, fn model.PostToolHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks(r).PostToolUse(func(ctx context.Context, hook run.Hook[sdkcopilot.PostToolUse], native sdkcopilot.PostToolResults) (sdkcopilot.PostToolOutput, error) {
		out, err := fn(ctx, model.NewPostToolHook(hook.Invocation(), mapPostToolUse(hook.Event)), newPostToolResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPostTool(out)
		if !ok {
			return nil, fmt.Errorf("copilot: PostTool result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPostToolUse(e sdkcopilot.PostToolUse) *model.PostToolEvent {
	return &model.PostToolEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     model.NewToolCall(e.NativeToolName(), e.Input().Raw(), ""),
		Result:   &model.ToolResult{Raw: hookkit.CloneRaw(e.ResultRaw()), Text: e.ResultText()},
	}
}

func newPostToolResults(native sdkcopilot.PostToolResults) model.PostToolResults {
	return postToolResults{native: native}
}

func unwrapPostTool(r model.PostToolResult) (sdkcopilot.PostToolOutput, bool) {
	out, ok := r.(postToolResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolResults struct {
	native sdkcopilot.PostToolResults
}

// Context returns a context-injection result.
func (w postToolResults) Context(text string) model.PostToolResult {
	return postToolResult{native: w.native.Context(text)}
}

type postToolResult struct {
	native sdkcopilot.PostToolOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
func (r postToolResult) WithUpdatedOutput(output string) model.PostToolResult {
	r.native = r.native.WithModifiedResult(output)
	return r
}
