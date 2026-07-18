package claude

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterPostToolFailure registers fn on the Claude PostToolUseFailure chain.
func RegisterPostToolFailure(fn model.PostToolFailureHandler) {
	if fn == nil {
		return
	}
	sdkclaude.OnPostToolUseFailure(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PostToolUseFailure], native sdkclaude.PostToolUseFailureResults) (sdkclaude.PostToolUseOutput, error) {
		out, err := fn(ctx, model.NewPostToolFailureHook(hook.Invocation(), mapPostToolUseFailure(hook.Event)), newPostToolFailureResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPostToolFailure(out)
		if !ok {
			return nil, fmt.Errorf("claude: PostToolFailure result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPostToolUseFailure(e sdkclaude.PostToolUseFailure) *model.PostToolFailureEvent {
	return &model.PostToolFailureEvent{
		Envelope: envelope(e),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
		Result:   &model.ToolResult{Error: e.Error},
	}
}

func newPostToolFailureResults(native sdkclaude.PostToolUseFailureResults) model.PostToolFailureResults {
	return postToolFailureResults{native: native}
}

func unwrapPostToolFailure(r model.PostToolFailureResult) (sdkclaude.PostToolUseOutput, bool) {
	out, ok := r.(postToolFailureResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolFailureResults struct {
	native sdkclaude.PostToolUseFailureResults
}

// Context returns a context-injection result.
func (w postToolFailureResults) Context(text string) model.PostToolFailureResult {
	return postToolFailureResult{native: w.native.Context(text)}
}

type postToolFailureResult struct {
	native sdkclaude.PostToolUseOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolFailureResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }
