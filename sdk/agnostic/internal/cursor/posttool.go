package cursor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterPostTool registers fn on Cursor PostToolUse, AfterShellExecution,
// AfterMCPExecution, and AfterFileEdit chains.
func RegisterPostTool(fn model.PostToolHandler) {
	if fn == nil {
		return
	}
	sdkcursor.OnPostToolUse(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PostToolUse], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		return callPostTool(ctx, hook.Invocation(), mapPostToolUse(hook.Event, hook.Raw()), native, fn)
	}).
		AfterShellExecution(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterShellExecution], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
			return callPostTool(ctx, hook.Invocation(), mapAfterShellExecution(hook.Event, hook.Raw()), native, fn)
		}).
		AfterMCPExecution(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterMCPExecution], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
			return callPostTool(ctx, hook.Invocation(), mapAfterMCPExecution(hook.Event, hook.Raw()), native, fn)
		}).
		AfterFileEdit(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterFileEdit], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
			return callPostTool(ctx, hook.Invocation(), mapAfterFileEdit(hook.Event, hook.Raw()), native, fn)
		})
}

func callPostTool(ctx context.Context, inv run.Invocation, ev *model.PostToolEvent, native sdkcursor.PostToolResults, fn model.PostToolHandler) (sdkcursor.PostToolOutput, error) {
	out, err := fn(ctx, model.NewPostToolHook(inv, ev), newPostToolResults(native))
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapPostTool(out)
	if !ok {
		return nil, fmt.Errorf("cursor: PostTool result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapPostToolUse(e sdkcursor.PostToolUse, raw []byte) *model.PostToolEvent {
	return &model.PostToolEvent{
		Envelope: envelope(e, raw),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
		Result:   &model.ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()},
	}
}

func mapAfterShellExecution(e sdkcursor.AfterShellExecution, raw []byte) *model.PostToolEvent {
	return &model.PostToolEvent{
		Envelope: envelope(e, raw),
		Tool:     &model.ToolCall{Name: tools.ToolBash, Native: receivedName(e), Shell: e.Command},
		Result:   &model.ToolResult{Text: e.Output, DurationMs: e.DurationMillis()},
	}
}

func mapAfterMCPExecution(e sdkcursor.AfterMCPExecution, raw []byte) *model.PostToolEvent {
	nameNorm, _ := hookkit.NormalizeToolName(e.ToolName)
	toolInput := tools.NewInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	name := receivedName(e)
	return &model.PostToolEvent{
		Envelope: envelope(e, raw),
		Tool: &model.ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  toolInput,
			MCP:    true,
		},
		Result: &model.ToolResult{Raw: json.RawMessage(e.ResultJSON), DurationMs: e.DurationMillis()},
	}
}

func mapAfterFileEdit(e sdkcursor.AfterFileEdit, raw []byte) *model.PostToolEvent {
	name := receivedName(e)
	ev := &model.PostToolEvent{
		Envelope: envelope(e, raw),
		Tool:     &model.ToolCall{Name: tools.ToolEdit, Native: name},
	}
	input, err := json.Marshal(map[string]any{
		"file_path": e.FilePath,
		"edits":     e.Edits,
	})
	if err != nil {
		return ev
	}
	editsRaw, err := json.Marshal(e.Edits)
	if err != nil {
		return ev
	}
	ev.Tool.Input = tools.NewInput(tools.ToolEdit, name, input)
	ev.Result = &model.ToolResult{Raw: editsRaw}
	return ev
}

func newPostToolResults(native sdkcursor.PostToolResults) model.PostToolResults {
	return postToolResults{native: native}
}

func unwrapPostTool(r model.PostToolResult) (sdkcursor.PostToolOutput, bool) {
	out, ok := r.(postToolResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type postToolResults struct {
	native sdkcursor.PostToolResults
}

// Context returns a context-injection result.
func (w postToolResults) Context(text string) model.PostToolResult {
	return postToolResult{native: w.native.Context(text)}
}

type postToolResult struct {
	native sdkcursor.PostToolOutput
}

// IsZero reports whether the result carries no instruction.
func (r postToolResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
func (r postToolResult) WithUpdatedOutput(output string) model.PostToolResult {
	r.native = r.native.WithUpdatedMCPOutput(output)
	return r
}
