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

// RegisterPreTool registers fn on Cursor PreToolUse, BeforeShellExecution, and BeforeMCPExecution chains.
func RegisterPreTool(fn model.PreToolHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcursor.UseHooks().PreToolUse(func(ctx context.Context, hook sdkcursor.PreToolUse, native sdkcursor.PreToolUseResults) (sdkcursor.PermissionOutput, error) {
		return callPreTool(ctx, mapPreToolUse(hook), newPreToolUseResults(native), fn)
	}).
		BeforeShellExecution(func(ctx context.Context, hook sdkcursor.BeforeShellExecution, native sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
			return callPreTool(ctx, mapBeforeShellExecution(hook), newPermissionResults(native), fn)
		}).
		BeforeMCPExecution(func(ctx context.Context, hook sdkcursor.BeforeMCPExecution, native sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
			return callPreTool(ctx, mapBeforeMCPExecution(hook), newPermissionResults(native), fn)
		})
}

// RegisterBeforeReadFile registers fn on the Cursor BeforeReadFile chain.
func RegisterBeforeReadFile(fn model.PreToolHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcursor.UseHooks().BeforeReadFile(func(ctx context.Context, hook sdkcursor.BeforeReadFile, native sdkcursor.BeforeReadFileResults) (sdkcursor.PermissionOutput, error) {
		return callPreTool(ctx, mapBeforeReadFile(hook), newBeforeReadResults(native), fn)
	})
}

func callPreTool(ctx context.Context, ev *model.PreToolEvent, results model.PreToolResults, fn model.PreToolHandler) (sdkcursor.PermissionOutput, error) {
	out, err := fn(ctx, *ev, results)
	if err != nil || out == nil {
		return nil, err
	}
	nativeOut, ok := unwrapPreTool(out)
	if !ok {
		return nil, fmt.Errorf("cursor: PreTool result must come from the injected Results builder")
	}
	return nativeOut, nil
}

func mapPreToolUse(e sdkcursor.PreToolUse) *model.PreToolEvent {
	ev := &model.PreToolEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
	}
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
	return ev
}

func mapBeforeShellExecution(e sdkcursor.BeforeShellExecution) *model.PreToolEvent {
	return &model.PreToolEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     &model.ToolCall{Name: tools.ToolBash, Native: e.EventName(), Shell: e.Command},
	}
}

func mapBeforeMCPExecution(e sdkcursor.BeforeMCPExecution) *model.PreToolEvent {
	nameNorm, _ := hookkit.NormalizeToolName(e.ToolName)
	toolInput := tools.NewInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	return &model.PreToolEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool: &model.ToolCall{
			Name:   nameNorm,
			Native: e.ToolName,
			Input:  toolInput,
			MCP:    true,
		},
	}
}

func mapBeforeReadFile(e sdkcursor.BeforeReadFile) *model.PreToolEvent {
	name := e.EventName()
	ev := &model.PreToolEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     &model.ToolCall{Name: tools.ToolRead, Native: name},
	}
	input, err := json.Marshal(map[string]any{
		"file_path":   e.FilePath,
		"content":     e.Content,
		"attachments": e.Attachments,
	})
	if err != nil {
		return ev
	}
	ev.Tool.Input = tools.NewInput(tools.ToolRead, name, input)
	return ev
}

func newPermissionResults(native sdkcursor.PermissionResults) model.PreToolResults {
	return permissionResults{native: native}
}

func newPreToolUseResults(native sdkcursor.PreToolUseResults) model.PreToolResults {
	return preToolUseResults{native: native}
}

func newBeforeReadResults(native sdkcursor.BeforeReadFileResults) model.PreToolResults {
	return beforeReadResults{native: native}
}

func unwrapPreTool(r model.PreToolResult) (sdkcursor.PermissionOutput, bool) {
	out, ok := r.(permissionResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type permissionResults struct {
	native sdkcursor.PermissionResults
}

// Allow returns an allow verdict.
func (w permissionResults) Allow() model.PreToolResult {
	return permissionResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w permissionResults) Deny(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w permissionResults) Ask(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Ask(reason)}
}

type preToolUseResults struct {
	native sdkcursor.PreToolUseResults
}

// Allow returns an allow verdict.
func (w preToolUseResults) Allow() model.PreToolResult {
	return permissionResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w preToolUseResults) Deny(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Deny(reason)}
}

// Ask encodes Cursor preToolUse permission "ask". Cursor accepts the value but
// does not enforce escalation for preToolUse today.
func (w preToolUseResults) Ask(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Ask(reason)}
}

type beforeReadResults struct {
	native sdkcursor.BeforeReadFileResults
}

// Allow returns an allow verdict.
func (w beforeReadResults) Allow() model.PreToolResult {
	return permissionResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w beforeReadResults) Deny(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w beforeReadResults) Ask(reason string) model.PreToolResult {
	return permissionResult{native: w.native.Ask(reason)}
}

type permissionResult struct {
	native sdkcursor.PermissionOutput
}

// IsZero reports whether the result carries no instruction.
func (r permissionResult) IsZero() bool { return r.native.IsZero() }

// WithUpdatedInput replaces tool arguments when set.
func (r permissionResult) WithUpdatedInput(input map[string]any) model.PreToolResult {
	r.native = r.native.WithUpdatedInput(input)
	return r
}
