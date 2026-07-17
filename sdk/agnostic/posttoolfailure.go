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

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *model.ToolCall
	Result *model.ToolResult
}

// PostToolFailureEventFrom maps a decoded Event to PostToolFailureEvent.
func PostToolFailureEventFrom(ev *model.Event) (PostToolFailureEvent, error) {
	if ev == nil {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindPostToolFailure {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: expected PostToolFailure kind, got %s", ev.Kind)
	}
	return PostToolFailureEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool, Result: ev.Result}, nil
}

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook struct {
	PostToolFailureEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolFailureHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolFailureHook) Raw() json.RawMessage { return h.PostToolFailureEvent.Raw }

func postToolFailureHook(ctx run.Invocation, ev PostToolFailureEvent) PostToolFailureHook {
	return PostToolFailureHook{PostToolFailureEvent: ev, inv: ctx}
}

// PostToolFailureResult is the portable hook response for PostToolFailure events.
// Construct only via PostToolFailureResults (Context).
// A nil value is a no-op.
type PostToolFailureResult interface {
	isPostToolFailureResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// PostToolFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolFailureResults interface {
	// Context returns recovery guidance for PostToolFailure events.
	Context(text string) PostToolFailureResult
	isPostToolFailureResults()
}

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, hook PostToolFailureHook, results PostToolFailureResults) (PostToolFailureResult, error)

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().PostToolUseFailure(adaptClaudePostToolFailure(fn))
	sdkcopilot.Adapter().PostToolUseFailure(adaptCopilotPostToolFailure(fn))
	sdkcursor.Adapter().PostToolUseFailure(adaptCursorPostToolFailure(fn))
	return &Chain{}
}

// OnPostToolFailure registers another PostToolFailure handler on the chain.
func (c *Chain) OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	return OnPostToolFailure(fn)
}

func adaptClaudePostToolFailure(fn PostToolFailureHandler) func(context.Context, sdkclaude.Hook[sdkclaude.PostToolUseFailure], sdkclaude.PostToolUseFailureResults) (sdkclaude.PostToolUseOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PostToolUseFailure], native sdkclaude.PostToolUseFailureResults) (sdkclaude.PostToolUseOutput, error) {
		typed, err := PostToolFailureEventFrom(agclaude.MapPostToolUseFailure(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, postToolFailureHook(hook.Invocation(), typed), claudePostToolFailureResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudePostToolFailureResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PostToolFailure result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type claudePostToolFailureResults struct {
	native sdkclaude.PostToolUseFailureResults
}

func (claudePostToolFailureResults) isPostToolFailureResults() {}

// Context returns a context-injection result.
func (w claudePostToolFailureResults) Context(text string) PostToolFailureResult {
	return claudePostToolFailureResult{native: w.native.Context(text)}
}

type claudePostToolFailureResult struct {
	native sdkclaude.PostToolUseOutput
}

func (claudePostToolFailureResult) isPostToolFailureResult() {}

// IsZero reports whether the result carries no instruction.
func (r claudePostToolFailureResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

func adaptCopilotPostToolFailure(fn PostToolFailureHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.PostToolUseFailure], sdkcopilot.PostToolFailureResults) (sdkcopilot.PostToolFailureOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PostToolUseFailure], native sdkcopilot.PostToolFailureResults) (sdkcopilot.PostToolFailureOutput, error) {
		typed, err := PostToolFailureEventFrom(agcopilot.MapPostToolUseFailure(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, postToolFailureHook(hook.Invocation(), typed), copilotPostToolFailureResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotPostToolFailureResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PostToolFailure result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type copilotPostToolFailureResults struct {
	native sdkcopilot.PostToolFailureResults
}

func (copilotPostToolFailureResults) isPostToolFailureResults() {}

// Context returns a context-injection result.
func (w copilotPostToolFailureResults) Context(text string) PostToolFailureResult {
	return copilotPostToolFailureResult{native: w.native.Context(text)}
}

type copilotPostToolFailureResult struct {
	native sdkcopilot.PostToolFailureOutput
}

func (copilotPostToolFailureResult) isPostToolFailureResult() {}

// IsZero reports whether the result carries no instruction.
func (r copilotPostToolFailureResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

func adaptCursorPostToolFailure(fn PostToolFailureHandler) func(context.Context, sdkcursor.Hook[sdkcursor.PostToolUseFailure], sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PostToolUseFailure], native sdkcursor.PostToolResults) (sdkcursor.PostToolOutput, error) {
		typed, err := PostToolFailureEventFrom(agcursor.MapPostToolUseFailure(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, postToolFailureHook(hook.Invocation(), typed), cursorPostToolFailureResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(cursorPostToolFailureResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: PostToolFailure result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type cursorPostToolFailureResults struct {
	native sdkcursor.PostToolResults
}

func (cursorPostToolFailureResults) isPostToolFailureResults() {}

// Context returns a context-injection result.
func (w cursorPostToolFailureResults) Context(text string) PostToolFailureResult {
	return cursorPostToolFailureResult{native: w.native.Context(text)}
}

type cursorPostToolFailureResult struct {
	native sdkcursor.PostToolOutput
}

func (cursorPostToolFailureResult) isPostToolFailureResult() {}

// IsZero reports whether the result carries no instruction.
func (r cursorPostToolFailureResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }
