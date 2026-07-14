package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
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

// PostToolResults is the hook-scoped response builder supplied to PostToolHandler by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) model.PostToolResult
	isPostToolResults()
}

type postToolResults struct{}

func (postToolResults) isPostToolResults() {}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) model.PostToolResult { return model.PostToolContext(text) }

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, hook PostToolHook, results PostToolResults) (model.PostToolResult, error)

// OnPostTool registers a handler for PostTool events across all agents.
func OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindPostTool, func(ctx context.Context, ev *model.Event) (model.PostToolResult, error) {
		typed, err := PostToolEventFrom(ev)
		if err != nil {
			return model.PostToolResult{}, err
		}
		return fn(ctx, postToolHook(run.InvocationFrom(ctx), typed), postToolResults{})
	})
	return &Chain{}
}

// OnPostTool registers another PostTool handler on the chain.
func (c *Chain) OnPostTool(fn PostToolHandler) *Chain {
	return OnPostTool(fn)
}
