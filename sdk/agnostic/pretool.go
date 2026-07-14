package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	// Tool holds tool invocation details.
	Tool *model.ToolCall
}

// PreToolEventFrom maps a decoded Event to PreToolEvent.
func PreToolEventFrom(ev *model.Event) (PreToolEvent, error) {
	if ev == nil {
		return PreToolEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindPreTool {
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

// PreToolResults is the hook-scoped response builder supplied to PreToolHandler by registration.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() model.PreToolResult
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) model.PreToolResult
	// Ask returns an escalate-to-user verdict with an agent-facing reason.
	Ask(reason string) model.PreToolResult
	isPreToolResults()
}

type preToolResults struct{}

func (preToolResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (preToolResults) Allow() model.PreToolResult { return model.PreToolAllow() }

// Deny returns a deny verdict with an agent-facing reason.
func (preToolResults) Deny(reason string) model.PreToolResult { return model.PreToolDeny(reason) }

// Ask returns an escalate-to-user verdict with an agent-facing reason.
func (preToolResults) Ask(reason string) model.PreToolResult { return model.PreToolAsk(reason) }

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, hook PreToolHook, results PreToolResults) (model.PreToolResult, error)

// OnPreTool registers a handler for PreTool events across all agents.
func OnPreTool(fn PreToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(model.KindPreTool, func(ctx context.Context, ev *model.Event) (model.PreToolResult, error) {
		typed, err := PreToolEventFrom(ev)
		if err != nil {
			return model.PreToolResult{}, err
		}
		return fn(ctx, preToolHook(run.InvocationFrom(ctx), typed), preToolResults{})
	})
	return &Chain{}
}

// OnPreTool registers another PreTool handler on the chain.
func (c *Chain) OnPreTool(fn PreToolHandler) *Chain {
	return OnPreTool(fn)
}
