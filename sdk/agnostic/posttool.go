package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent struct {
	Envelope
	Tool   *model.ToolCall
	Result *model.ToolResult
}

// PostToolEventFrom maps a decoded Event to PostToolEvent.
func PostToolEventFrom(ev *model.Event) (PostToolEvent, error) {
	if ev == nil {
		return PostToolEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindPostTool {
		return PostToolEvent{}, fmt.Errorf("agnostic: expected PostTool kind, got %s", ev.Kind)
	}
	return PostToolEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool, Result: ev.Result}, nil
}

// PostToolHook is the handler context for portable PostTool events.
type PostToolHook struct {
	PostToolEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolHook) Raw() json.RawMessage { return h.PostToolEvent.Raw }

func postToolHook(ctx run.Invocation, ev PostToolEvent) PostToolHook {
	return PostToolHook{PostToolEvent: ev, inv: ctx}
}

// PostToolResult is the portable hook response for PostTool events.
// Construct only via PostToolResults (Context), then With*.
// A nil value is a no-op.
type PostToolResult interface {
	isPostToolResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// WithUpdatedOutput replaces tool result text when set.
	// On Cursor this maps to updated_mcp_tool_output (MCP AfterMCP / post-tool only).
	WithUpdatedOutput(output string) PostToolResult
}

// PostToolResults is the hook-scoped response builder supplied to PostToolHandler by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolResult
	isPostToolResults()
}

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, hook PostToolHook, results PostToolResults) (PostToolResult, error)

// OnPostTool registers a handler for PostTool events across all agents.
func OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().PostToolUse(adaptClaudePostTool(fn))
	sdkcopilot.Adapter().PostToolUse(adaptCopilotPostTool(fn))
	sdkcursor.Adapter().
		PostToolUse(adaptCursorPostTool(fn)).
		AfterShellExecution(adaptCursorAfterShell(fn)).
		AfterMCPExecution(adaptCursorAfterMCP(fn)).
		AfterFileEdit(adaptCursorAfterFileEdit(fn))
	return &Chain{}
}

// OnPostTool registers another PostTool handler on the chain.
func (c *Chain) OnPostTool(fn PostToolHandler) *Chain {
	return OnPostTool(fn)
}

func adaptClaudePostTool(fn PostToolHandler) func(context.Context, sdkclaude.Hook[sdkclaude.PostToolUse], sdkclaude.PostToolUseResults) (sdkclaude.PostToolUseOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PostToolUse], native sdkclaude.PostToolUseResults) (sdkclaude.PostToolUseOutput, error) {
		typed, err := PostToolEventFrom(agclaude.MapEvent(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, postToolHook(hook.Invocation(), typed), claudePostToolResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudePostToolResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PostTool result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type claudePostToolResults struct {
	native sdkclaude.PostToolUseResults
}

func (claudePostToolResults) isPostToolResults() {}

// Context returns a context-injection result.
func (w claudePostToolResults) Context(text string) PostToolResult {
	return claudePostToolResult{native: w.native.Context(text)}
}

type claudePostToolResult struct {
	native sdkclaude.PostToolUseOutput
}

func (claudePostToolResult) isPostToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r claudePostToolResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
func (r claudePostToolResult) WithUpdatedOutput(output string) PostToolResult {
	r.native = r.native.WithUpdatedToolOutput(output)
	return r
}

func adaptCopilotPostTool(fn PostToolHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.PostToolUse], sdkcopilot.PostToolResults) (sdkcopilot.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PostToolUse], native sdkcopilot.PostToolResults) (sdkcopilot.PostToolOutput, error) {
		typed, err := PostToolEventFrom(agcopilot.MapEvent(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, postToolHook(hook.Invocation(), typed), copilotPostToolResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotPostToolResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PostTool result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type copilotPostToolResults struct {
	native sdkcopilot.PostToolResults
}

func (copilotPostToolResults) isPostToolResults() {}

// Context returns a context-injection result.
func (w copilotPostToolResults) Context(text string) PostToolResult {
	return copilotPostToolResult{native: w.native.Context(text)}
}

type copilotPostToolResult struct {
	native sdkcopilot.PostToolOutput
}

func (copilotPostToolResult) isPostToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r copilotPostToolResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
func (r copilotPostToolResult) WithUpdatedOutput(output string) PostToolResult {
	r.native = r.native.WithModifiedResult(output)
	return r
}

func adaptCursorPostTool(fn PostToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.PostToolUse], sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PostToolUse], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		return callCursorPostTool(ctx, hook.Invocation(), agcursor.MapEvent(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorAfterShell(fn PostToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.AfterShellExecution], sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterShellExecution], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		return callCursorPostTool(ctx, hook.Invocation(), agcursor.MapEvent(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorAfterMCP(fn PostToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.AfterMCPExecution], sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterMCPExecution], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		return callCursorPostTool(ctx, hook.Invocation(), agcursor.MapEvent(hook.Event, hook.Raw()), native, fn)
	}
}

func adaptCursorAfterFileEdit(fn PostToolHandler) func(context.Context, sdkcursor.Hook[sdkcursor.AfterFileEdit], sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.AfterFileEdit], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		return callCursorPostTool(ctx, hook.Invocation(), agcursor.MapEvent(hook.Event, hook.Raw()), native, fn)
	}
}

func callCursorPostTool(ctx context.Context, inv run.Invocation, ev *model.Event, native sdkcursor.PostToolResults, fn PostToolHandler) (sdkcursor.PostToolOutput, error) {
	typed, err := PostToolEventFrom(ev)
	if err != nil {
		return nil, err
	}
	out, err := fn(ctx, postToolHook(inv, typed), cursorPostToolResults{native: native})
	if err != nil || out == nil {
		return nil, err
	}
	r, ok := out.(cursorPostToolResult)
	if !ok {
		return nil, fmt.Errorf("agnostic: PostTool result must come from the injected Results builder")
	}
	return r.native, nil
}

type cursorPostToolResults struct {
	native sdkcursor.PostToolResults
}

func (cursorPostToolResults) isPostToolResults() {}

// Context returns a context-injection result.
func (w cursorPostToolResults) Context(text string) PostToolResult {
	return cursorPostToolResult{native: w.native.Context(text)}
}

type cursorPostToolResult struct {
	native sdkcursor.PostToolOutput
}

func (cursorPostToolResult) isPostToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r cursorPostToolResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }

// WithUpdatedOutput replaces tool result text when set.
// On Cursor this maps to updated_mcp_tool_output (MCP post-tool only).
func (r cursorPostToolResult) WithUpdatedOutput(output string) PostToolResult {
	r.native = r.native.WithUpdatedMCPOutput(output)
	return r
}
