package claude

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterPostTool registers fn on the Claude PostToolUse chain.
func RegisterPostTool(fn model.PostToolHandler) {
	if fn == nil {
		return
	}
	sdkclaude.OnPostToolUse(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PostToolUse], native sdkclaude.PostToolUseResults) (sdkclaude.PostToolUseOutput, error) {
		out, err := fn(ctx, model.NewPostToolHook(hook.Invocation(), mapPostToolUse(hook.Event)), newPostToolResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPostTool(out)
		if !ok {
			return nil, fmt.Errorf("claude: PostTool result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPostToolUse(e sdkclaude.PostToolUse) *model.PostToolEvent {
	return &model.PostToolEvent{
		Envelope: envelope(e),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
		Result:   &model.ToolResult{Raw: hookkit.CloneRaw(e.ToolResponse), Text: hookkit.RawToText(e.ToolResponse)},
	}
}

func newPostToolResults(native sdkclaude.PostToolUseResults) model.PostToolResults {
	return postToolResults{native: native}
}

func unwrapPostTool(r model.PostToolResult) (sdkclaude.PostToolUseOutput, bool) {
	out, ok := r.(postToolResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolResults struct {
	native sdkclaude.PostToolUseResults
}

// Context returns a context-injection result.
func (w postToolResults) Context(text string) model.PostToolResult {
	return postToolResult{native: w.native.Context(text)}
}

type postToolResult struct {
	native sdkclaude.PostToolUseOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
func (r postToolResult) WithUpdatedOutput(output string) model.PostToolResult {
	r.native = r.native.WithUpdatedToolOutput(output)
	return r
}
