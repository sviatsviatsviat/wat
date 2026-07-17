package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	// Tool holds tool invocation details.
	Tool *ToolCall
}

// PreToolEventFrom maps a decoded Event to PreToolEvent.
func PreToolEventFrom(ev *Event) (PreToolEvent, error) {
	if ev == nil {
		return PreToolEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPreTool {
		return PreToolEvent{}, fmt.Errorf("agnostic: expected PreTool kind, got %s", ev.Kind)
	}
	return PreToolEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool}, nil
}

// PreToolHook is the handler context for portable PreTool events.
type PreToolHook struct {
	PreToolEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreToolHook) Raw() json.RawMessage { return h.PreToolEvent.Raw }

func preToolHook(ctx run.Invocation, ev PreToolEvent) PreToolHook {
	return PreToolHook{PreToolEvent: ev, inv: ctx}
}

// PreToolResult is the portable hook response for PreTool events.
// Construct only via PreToolResults (Allow/Deny/Ask), then With*.
// A nil value is a no-op.
type PreToolResult interface {
	isPreToolResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// WithUpdatedInput replaces tool arguments when set.
	// On Cursor, updated_input is emitted only for preToolUse (not beforeShellExecution,
	// beforeMCPExecution, or beforeReadFile).
	WithUpdatedInput(input map[string]any) PreToolResult
}

// PreToolResults is the hook-scoped response builder supplied to PreToolHandler by registration.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolResult
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolResult
	// Ask returns an escalate-to-user verdict with an agent-facing reason.
	Ask(reason string) PreToolResult
	isPreToolResults()
}

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, hook PreToolHook, results PreToolResults) (PreToolResult, error)

// OnPreTool registers a handler for PreTool events across all agents.
func OnPreTool(fn PreToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().PreToolUse(adaptClaudePreTool(fn))
	sdkcopilot.Adapter().PreToolUse(adaptCopilotPreTool(fn))
	sdkcursor.Adapter().
		PreToolUse(adaptCursorPreTool(fn)).
		BeforeShellExecution(adaptCursorBeforeShell(fn)).
		BeforeMCPExecution(adaptCursorBeforeMCP(fn)).
		BeforeReadFile(adaptCursorBeforeRead(fn))
	return &Chain{}
}

// OnPreTool registers another PreTool handler on the chain.
func (c *Chain) OnPreTool(fn PreToolHandler) *Chain {
	return OnPreTool(fn)
}

func adaptClaudePreTool(fn PreToolHandler) func(context.Context, sdkclaude.Hook[sdkclaude.PreToolUse], sdkclaude.PreToolUseResults) (sdkclaude.PreToolUseOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PreToolUse], native sdkclaude.PreToolUseResults) (sdkclaude.PreToolUseOutput, error) {
		typed, err := PreToolEventFrom(mapClaudePreToolUse(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, preToolHook(hook.Invocation(), typed), claudePreToolResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudePreToolResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PreTool result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type claudePreToolResults struct {
	native sdkclaude.PreToolUseResults
}

func (claudePreToolResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (w claudePreToolResults) Allow() PreToolResult {
	return claudePreToolResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w claudePreToolResults) Deny(reason string) PreToolResult {
	return claudePreToolResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w claudePreToolResults) Ask(reason string) PreToolResult {
	return claudePreToolResult{native: w.native.Ask(reason)}
}

type claudePreToolResult struct {
	native sdkclaude.PreToolUseOutput
}

func (claudePreToolResult) isPreToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r claudePreToolResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

// WithUpdatedInput replaces tool arguments when set.
func (r claudePreToolResult) WithUpdatedInput(input map[string]any) PreToolResult {
	r.native = r.native.WithUpdatedInput(input)
	return r
}

func adaptCopilotPreTool(fn PreToolHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.PreToolUse], sdkcopilot.PreToolResults) (sdkcopilot.PreToolOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PreToolUse], native sdkcopilot.PreToolResults) (sdkcopilot.PreToolOutput, error) {
		typed, err := PreToolEventFrom(mapCopilotPreToolUse(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, preToolHook(hook.Invocation(), typed), copilotPreToolResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotPreToolResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PreTool result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type copilotPreToolResults struct {
	native sdkcopilot.PreToolResults
}

func (copilotPreToolResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (w copilotPreToolResults) Allow() PreToolResult {
	return copilotPreToolResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w copilotPreToolResults) Deny(reason string) PreToolResult {
	return copilotPreToolResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w copilotPreToolResults) Ask(reason string) PreToolResult {
	return copilotPreToolResult{native: w.native.Ask(reason)}
}

type copilotPreToolResult struct {
	native sdkcopilot.PreToolOutput
}

func (copilotPreToolResult) isPreToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r copilotPreToolResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

// WithUpdatedInput replaces tool arguments when set.
func (r copilotPreToolResult) WithUpdatedInput(input map[string]any) PreToolResult {
	r.native = r.native.WithModifiedArgs(input)
	return r
}

func adaptCursorPreTool(fn PreToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.PreToolUse], sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PreToolUse], native sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
		return callCursorPreTool(ctx, hook.Invocation(), mapCursorPreToolUse(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorBeforeShell(fn PreToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.BeforeShellExecution], sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.BeforeShellExecution], native sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
		return callCursorPreTool(ctx, hook.Invocation(), mapCursorBeforeShellExecution(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorBeforeMCP(fn PreToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.BeforeMCPExecution], sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.BeforeMCPExecution], native sdkcursor.PermissionResults) (sdkcursor.PermissionOutput, error) {
		return callCursorPreTool(ctx, hook.Invocation(), mapCursorBeforeMCPExecution(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorBeforeRead(fn PreToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.BeforeReadFile], sdkcursor.BeforeReadFileResults) (sdkcursor.PermissionOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.BeforeReadFile], native sdkcursor.BeforeReadFileResults) (sdkcursor.PermissionOutput, error) {
		typed, err := PreToolEventFrom(mapCursorBeforeReadFile(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, preToolHook(hook.Invocation(), typed), cursorBeforeReadResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(cursorPermissionResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PreTool result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

func callCursorPreTool(ctx context.Context, inv run.Invocation, ev *Event, native sdkcursor.PermissionResults, fn PreToolHandler) (sdkcursor.PermissionOutput, error) {
	typed, err := PreToolEventFrom(ev)
	if err != nil {
		return nil, err
	}
	out, err := fn(ctx, preToolHook(inv, typed), cursorPermissionResults{native: native})
	if err != nil || out == nil {
		return nil, err
	}
	r, ok := out.(cursorPermissionResult)
	if !ok {
		return nil, fmt.Errorf("agnostic: PreTool result must come from the injected Results builder")
	}
	return r.native, nil
}

type cursorPermissionResults struct {
	native sdkcursor.PermissionResults
}

func (cursorPermissionResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (w cursorPermissionResults) Allow() PreToolResult {
	return cursorPermissionResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w cursorPermissionResults) Deny(reason string) PreToolResult {
	return cursorPermissionResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w cursorPermissionResults) Ask(reason string) PreToolResult {
	return cursorPermissionResult{native: w.native.Ask(reason)}
}

type cursorBeforeReadResults struct {
	native sdkcursor.BeforeReadFileResults
}

func (cursorBeforeReadResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (w cursorBeforeReadResults) Allow() PreToolResult {
	return cursorPermissionResult{native: w.native.Allow()}
}

// Deny returns a deny verdict with an agent-facing reason.
func (w cursorBeforeReadResults) Deny(reason string) PreToolResult {
	return cursorPermissionResult{native: w.native.Deny(reason)}
}

// Ask returns an ask verdict with an agent-facing reason.
func (w cursorBeforeReadResults) Ask(reason string) PreToolResult {
	return cursorPermissionResult{native: w.native.Ask(reason)}
}

type cursorPermissionResult struct {
	native sdkcursor.PermissionOutput
}

func (cursorPermissionResult) isPreToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r cursorPermissionResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }

// WithUpdatedInput replaces tool arguments when set.
// On Cursor, updated_input is emitted only for preToolUse.
func (r cursorPermissionResult) WithUpdatedInput(input map[string]any) PreToolResult {
	r.native = r.native.WithUpdatedInput(input)
	return r
}

func mapClaudePreToolUse(e sdkclaude.PreToolUse, raw []byte) *Event {
	ev := claudeEvent(e, raw, KindPreTool)
	ev.Tool = newToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	return ev
}

func mapCopilotPreToolUse(e sdkcopilot.PreToolUse, raw []byte) *Event {
	ev := copilotEvent(e, raw, KindPreTool)
	ev.Tool = newToolCall(e.NativeToolName(), e.Input().Raw(), "")
	return ev
}

func mapCursorPreToolUse(e sdkcursor.PreToolUse, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindPreTool)
	ev.Tool = newToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
	return ev
}

func mapCursorBeforeShellExecution(e sdkcursor.BeforeShellExecution, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindPreTool)
	ev.Tool = &ToolCall{Name: tools.ToolBash, Native: cursorReceivedName(e), Shell: e.Command}
	return ev
}

func mapCursorBeforeMCPExecution(e sdkcursor.BeforeMCPExecution, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindPreTool)
	name := cursorReceivedName(e)
	nameNorm, _ := hookkit.NormalizeToolName(e.ToolName)
	toolInput := tools.NewInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	ev.Tool = &ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  toolInput,
		MCP:    true,
	}
	return ev
}

func mapCursorBeforeReadFile(e sdkcursor.BeforeReadFile, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindPreTool)
	name := cursorReceivedName(e)
	ev.Tool = &ToolCall{Name: tools.ToolRead, Native: name}
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
