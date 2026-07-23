package cursor

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterPostToolFailure registers fn on the Cursor PostToolUseFailure chain.
func RegisterPostToolFailure(fn model.PostToolFailureHandler) {
	if fn == nil {
		return
	}
	sdkcursor.UseHooks().PostToolUseFailure(func(ctx context.Context, hook sdkcursor.PostToolUseFailure, native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		out, err := fn(ctx, *mapPostToolUseFailure(hook), newPostToolFailureResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPostToolFailure(out)
		if !ok {
			return nil, fmt.Errorf("cursor: PostToolFailure result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPostToolUseFailure(e sdkcursor.PostToolUseFailure) *model.PostToolFailureEvent {
	return &model.PostToolFailureEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
		Result: &model.ToolResult{
			Error:       e.ErrorMessage,
			FailureType: e.FailureType,
			DurationMs:  e.DurationMillis(),
		},
	}
}

func newPostToolFailureResults(native sdkcursor.PostToolResults) model.PostToolFailureResults {
	return postToolFailureResults{native: native}
}

func unwrapPostToolFailure(r model.PostToolFailureResult) (sdkcursor.PostToolOutput, bool) {
	out, ok := r.(postToolFailureResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolFailureResults struct {
	native sdkcursor.PostToolResults
}

// Context returns a context-injection result.
func (w postToolFailureResults) Context(text string) model.PostToolFailureResult {
	return postToolFailureResult{native: w.native.Context(text)}
}

type postToolFailureResult struct {
	native sdkcursor.PostToolOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolFailureResult) IsZero() bool { return r.native.IsZero() }
