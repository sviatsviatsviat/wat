package claude

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterPreTool registers fn on the Claude PreToolUse chain.
func RegisterPreTool(fn model.PreToolHandler) {
	if fn == nil {
		return
	}
	sdkclaude.Adapter().PreToolUse(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PreToolUse], native sdkclaude.PreToolUseResults) (sdkclaude.PreToolUseOutput, error) {
		out, err := fn(ctx, model.NewPreToolHook(hook.Invocation(), mapPreToolUse(hook.Event, hook.Raw())), newPreToolResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPreTool(out)
		if !ok {
			return nil, fmt.Errorf("claude: PreTool result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPreToolUse(e sdkclaude.PreToolUse, raw []byte) *model.PreToolEvent {
	return &model.PreToolEvent{
		Envelope: envelope(e, raw),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
	}
}

func newPreToolResults(native sdkclaude.PreToolUseResults) model.PreToolResults {
	return preToolResults{native: native}
}

func unwrapPreTool(r model.PreToolResult) (sdkclaude.PreToolUseOutput, bool) {
	out, ok := r.(preToolResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type preToolResults struct {
	native sdkclaude.PreToolUseResults
}

// Allow returns an allow verdict.
func (w preToolResults) Allow() model.PreToolResult {
	return preToolResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w preToolResults) Deny(reason string) model.PreToolResult {
	return preToolResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w preToolResults) Ask(reason string) model.PreToolResult {
	return preToolResult{native: w.native.Ask(reason)}
}

type preToolResult struct {
	native sdkclaude.PreToolUseOutput
}

// IsZero reports whether the result carries no instruction.
func (r preToolResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

// WithUpdatedInput replaces tool arguments when set.
func (r preToolResult) WithUpdatedInput(input map[string]any) model.PreToolResult {
	r.native = r.native.WithUpdatedInput(input)
	return r
}
