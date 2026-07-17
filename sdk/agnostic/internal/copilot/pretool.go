package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPreTool registers fn on the Copilot PreToolUse chain.
func RegisterPreTool(fn model.PreToolHandler) {
	if fn == nil {
		return
	}
	new(sdkcopilot.Chain).PreToolUse(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PreToolUse], native sdkcopilot.PreToolResults) (sdkcopilot.PreToolOutput, error) {
		out, err := fn(ctx, model.NewPreToolHook(hook.Invocation(), mapPreToolUse(hook.Event, hook.Raw())), newPreToolResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapPreTool(out)
		if !ok {
			return nil, fmt.Errorf("copilot: PreTool result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapPreToolUse(e sdkcopilot.PreToolUse, raw []byte) *model.PreToolEvent {
	return &model.PreToolEvent{
		Envelope: envelope(e, raw),
		Tool:     model.NewToolCall(e.NativeToolName(), e.Input().Raw(), ""),
	}
}

func newPreToolResults(native sdkcopilot.PreToolResults) model.PreToolResults {
	return preToolResults{native: native}
}

func unwrapPreTool(r model.PreToolResult) (sdkcopilot.PreToolOutput, bool) {
	out, ok := r.(preToolResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type preToolResults struct {
	native sdkcopilot.PreToolResults
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
	native sdkcopilot.PreToolOutput
}

// IsZero reports whether the result carries no instruction.
func (r preToolResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

// WithUpdatedInput replaces tool arguments when set.
func (r preToolResult) WithUpdatedInput(input map[string]any) model.PreToolResult {
	r.native = r.native.WithModifiedArgs(input)
	return r
}
