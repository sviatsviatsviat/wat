package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	// Life holds the lifecycle information associated with this session start event.
	Life *lifecycle
}

// SessionStartEventFrom maps a decoded Event to SessionStartEvent.
func SessionStartEventFrom(ev *Event) (SessionStartEvent, error) {
	if ev == nil {
		return SessionStartEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindSessionStart {
		return SessionStartEvent{}, fmt.Errorf("agnostic: expected SessionStart kind, got %s", ev.Kind)
	}
	return SessionStartEvent{Envelope: envelopeFrom(ev), Life: ev.Life}, nil
}

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook struct {
	SessionStartEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionStartHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h SessionStartHook) Raw() json.RawMessage { return h.SessionStartEvent.Raw }

func sessionStartHook(ctx run.Invocation, ev SessionStartEvent) SessionStartHook {
	return SessionStartHook{SessionStartEvent: ev, inv: ctx}
}

// SessionStartResult is the portable hook response for SessionStart events.
// Construct only via SessionStartResults (Context).
// A nil value is a no-op.
type SessionStartResult interface {
	isSessionStartResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// SessionStartResults is the hook-scoped response builder supplied to OnSessionStart handlers by registration.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartResult
	isSessionStartResults()
}

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, hook SessionStartHook, results SessionStartResults) (SessionStartResult, error)

// OnSessionStart registers a handler for SessionStart events across all agents.
func OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	sdkclaude.Adapter().SessionStart(adaptClaudeSessionStart(fn))
	sdkcopilot.Adapter().SessionStart(adaptCopilotSessionStart(fn))
	sdkcursor.Adapter().SessionStart(adaptCursorSessionStart(fn))
	return &Chain{}
}

// OnSessionStart registers another SessionStart handler on the chain.
func (c *Chain) OnSessionStart(fn SessionStartHandler) *Chain {
	return OnSessionStart(fn)
}

func adaptClaudeSessionStart(fn SessionStartHandler) func(context.Context, sdkclaude.Hook[sdkclaude.SessionStart], sdkclaude.SessionStartResults) (sdkclaude.SessionStartOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SessionStart], native sdkclaude.SessionStartResults) (sdkclaude.SessionStartOutput, error) {
		typed, err := SessionStartEventFrom(mapClaudeSessionStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, sessionStartHook(hook.Invocation(), typed), claudeSessionStartResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(claudeSessionStartResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: SessionStart result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type claudeSessionStartResults struct {
	native sdkclaude.SessionStartResults
}

func (claudeSessionStartResults) isSessionStartResults() {}

// Context returns a context-injection result.
func (w claudeSessionStartResults) Context(text string) SessionStartResult {
	return claudeSessionStartResult{native: w.native.Context(text)}
}

type claudeSessionStartResult struct {
	native sdkclaude.SessionStartOutput
}

func (claudeSessionStartResult) isSessionStartResult() {}

// IsZero reports whether the result carries no instruction.
func (r claudeSessionStartResult) IsZero() bool { return sdkclaude.IsZeroOutput(r.native) }

func adaptCopilotSessionStart(fn SessionStartHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.SessionStart], sdkcopilot.SessionStartResults) (sdkcopilot.SessionStartOutput, error) {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SessionStart], native sdkcopilot.SessionStartResults) (sdkcopilot.SessionStartOutput, error) {
		typed, err := SessionStartEventFrom(mapCopilotSessionStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, sessionStartHook(hook.Invocation(), typed), copilotSessionStartResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(copilotSessionStartResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: SessionStart result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type copilotSessionStartResults struct {
	native sdkcopilot.SessionStartResults
}

func (copilotSessionStartResults) isSessionStartResults() {}

// Context returns a context-injection result.
func (w copilotSessionStartResults) Context(text string) SessionStartResult {
	return copilotSessionStartResult{native: w.native.Context(text)}
}

type copilotSessionStartResult struct {
	native sdkcopilot.SessionStartOutput
}

func (copilotSessionStartResult) isSessionStartResult() {}

// IsZero reports whether the result carries no instruction.
func (r copilotSessionStartResult) IsZero() bool { return sdkcopilot.IsZeroOutput(r.native) }

func adaptCursorSessionStart(fn SessionStartHandler) func(context.Context, sdkcursor.Hook[sdkcursor.SessionStart], sdkcursor.SessionStartResults) (sdkcursor.SessionStartOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SessionStart], native sdkcursor.SessionStartResults) (sdkcursor.SessionStartOutput, error) {
		typed, err := SessionStartEventFrom(mapCursorSessionStart(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, sessionStartHook(hook.Invocation(), typed), cursorSessionStartResults{native: native})
		if err != nil || out == nil {
			return nil, err
		}
		r, ok := out.(cursorSessionStartResult)
		if !ok {
			return nil, fmt.Errorf("agnostic: SessionStart result must come from the injected Results builder")
		}
		return r.native, nil
	}
}

type cursorSessionStartResults struct {
	native sdkcursor.SessionStartResults
}

func (cursorSessionStartResults) isSessionStartResults() {}

// Context returns a context-injection result.
func (w cursorSessionStartResults) Context(text string) SessionStartResult {
	return cursorSessionStartResult{native: w.native.Context(text)}
}

type cursorSessionStartResult struct {
	native sdkcursor.SessionStartOutput
}

func (cursorSessionStartResult) isSessionStartResult() {}

// IsZero reports whether the result carries no instruction.
func (r cursorSessionStartResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }

func mapClaudeSessionStart(e sdkclaude.SessionStart, raw []byte) *Event {
	ev := claudeEvent(e, raw, KindSessionStart)
	ev.Life = &lifecycle{Source: e.Source, Model: e.Model}
	return ev
}

func mapCopilotSessionStart(e sdkcopilot.SessionStart, raw []byte) *Event {
	ev := copilotEvent(e, raw, KindSessionStart)
	ev.Life = &lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
	return ev
}

func mapCursorSessionStart(e sdkcursor.SessionStart, raw []byte) *Event {
	ev := cursorEvent(e, raw, KindSessionStart)
	ev.Life = &lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
	return ev
}
