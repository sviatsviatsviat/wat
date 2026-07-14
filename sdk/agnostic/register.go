package agnostic

import (
	"context"
	"fmt"
	"io"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, ev *Event, results PreToolResults) (PreToolResult, error)

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, ev *Event, results PostToolResults) (PostToolResult, error)

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, ev *Event, results PostToolFailureResults) (PostToolFailureResult, error)

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, ev *Event, results StopResults) (StopResult, error)

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, ev *Event, results SessionStartResults) (SessionStartResult, error)

// ObserveHandler handles observe-only portable events with no hook response.
type ObserveHandler func(ctx context.Context, ev *Event) error

// Chain supports fluent handler registration into the shared run registry.
type Chain struct{}

// OnPreTool registers a handler for PreTool events across all agents.
func OnPreTool(fn PreToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPreTool, func(ctx context.Context, ev *Event) (PreToolResult, error) {
		return fn(ctx, ev, preToolResults{})
	})
	return &Chain{}
}

// OnPostTool registers a handler for PostTool events across all agents.
func OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPostTool, func(ctx context.Context, ev *Event) (PostToolResult, error) {
		return fn(ctx, ev, postToolResults{})
	})
	return &Chain{}
}

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPostToolFailure, func(ctx context.Context, ev *Event) (PostToolFailureResult, error) {
		return fn(ctx, ev, postToolFailureResults{})
	})
	return &Chain{}
}

// OnStop registers a handler for Stop events across all agents.
func OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		return fn(ctx, ev, stopResults{})
	})
	return &Chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindSubagentStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		return fn(ctx, ev, stopResults{})
	})
	return &Chain{}
}

// OnSessionStart registers a handler for SessionStart events across all agents.
func OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindSessionStart, func(ctx context.Context, ev *Event) (SessionStartResult, error) {
		return fn(ctx, ev, sessionStartResults{})
	})
	return &Chain{}
}

// OnSessionEnd registers an observe-only handler for SessionEnd events.
func OnSessionEnd(fn ObserveHandler) *Chain {
	registerObserveHandler(KindSessionEnd, fn)
	return &Chain{}
}

// OnUserPrompt registers an observe-only handler for UserPrompt events.
func OnUserPrompt(fn ObserveHandler) *Chain {
	registerObserveHandler(KindUserPrompt, fn)
	return &Chain{}
}

// OnPreCompact registers an observe-only handler for PreCompact events.
func OnPreCompact(fn ObserveHandler) *Chain {
	registerObserveHandler(KindPreCompact, fn)
	return &Chain{}
}

// OnSubagentStart registers an observe-only handler for SubagentStart events.
func OnSubagentStart(fn ObserveHandler) *Chain {
	registerObserveHandler(KindSubagentStart, fn)
	return &Chain{}
}

// OnAny registers an observe-only handler invoked for every event before kind-specific handlers.
func OnAny(fn ObserveHandler) *Chain {
	registerAny(fn)
	return &Chain{}
}

// OnPostTool registers another PostTool handler on the chain.
func (c *Chain) OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindPostTool, func(ctx context.Context, ev *Event) (PostToolResult, error) {
		return fn(ctx, ev, postToolResults{})
	})
	return c
}

// OnPostToolFailure registers another PostToolFailure handler on the chain.
func (c *Chain) OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindPostToolFailure, func(ctx context.Context, ev *Event) (PostToolFailureResult, error) {
		return fn(ctx, ev, postToolFailureResults{})
	})
	return c
}

// OnStop registers another Stop handler on the chain.
func (c *Chain) OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		return fn(ctx, ev, stopResults{})
	})
	return c
}

// OnSubagentStop registers another SubagentStop handler on the chain.
func (c *Chain) OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindSubagentStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		return fn(ctx, ev, stopResults{})
	})
	return c
}

// OnSessionStart registers another SessionStart handler on the chain.
func (c *Chain) OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindSessionStart, func(ctx context.Context, ev *Event) (SessionStartResult, error) {
		return fn(ctx, ev, sessionStartResults{})
	})
	return c
}

// OnSessionEnd registers another observe-only SessionEnd handler on the chain.
func (c *Chain) OnSessionEnd(fn ObserveHandler) *Chain {
	registerObserveHandler(KindSessionEnd, fn)
	return c
}

// OnUserPrompt registers another observe-only UserPrompt handler on the chain.
func (c *Chain) OnUserPrompt(fn ObserveHandler) *Chain {
	registerObserveHandler(KindUserPrompt, fn)
	return c
}

// OnPreCompact registers another observe-only PreCompact handler on the chain.
func (c *Chain) OnPreCompact(fn ObserveHandler) *Chain {
	registerObserveHandler(KindPreCompact, fn)
	return c
}

// OnSubagentStart registers another observe-only SubagentStart handler on the chain.
func (c *Chain) OnSubagentStart(fn ObserveHandler) *Chain {
	registerObserveHandler(KindSubagentStart, fn)
	return c
}

// OnAny registers another observe-only catch-all handler on the chain.
func (c *Chain) OnAny(fn ObserveHandler) *Chain {
	registerAny(fn)
	return c
}

// OnPreTool registers another PreTool handler on the chain.
func (c *Chain) OnPreTool(fn PreToolHandler) *Chain {
	if fn == nil {
		return c
	}
	registerResultHandler(KindPreTool, func(ctx context.Context, ev *Event) (PreToolResult, error) {
		return fn(ctx, ev, preToolResults{})
	})
	return c
}

// Serve reads a hook payload from in, dispatches registered handlers, writes
// encoded stdout to out, diagnostics to errw, and returns the process exit code.
func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...run.Option) int {
	return run.Serve(ctx, in, out, errw, opts...)
}

// WithDialect forces a dialect instead of auto-detection.
func WithDialect(d Dialect) run.Option {
	return run.WithDialect(d.String())
}

// WithEvent supplies the native event name for Copilot camelCase payloads that
// omit hook_event_name.
func WithEvent(name string) run.Option {
	return run.WithEvent(name)
}

// WithGetenv injects environment lookup for Detect.
func WithGetenv(getenv func(string) string) run.Option {
	return run.WithGetenv(getenv)
}

type wireResult interface {
	Result() model.Result
}

func registerResultHandler[R wireResult](kind Kind, fn func(context.Context, *Event) (R, error)) {
	if fn == nil {
		return
	}
	wrap := func(ctx context.Context, ev *Event) (model.Result, error) {
		res, err := fn(ctx, ev)
		if err != nil {
			return model.Result{}, err
		}
		return res.Result(), nil
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		for _, eventName := range eventsForKind(agent, kind) {
			registerAgnosticHandler(agent, eventName, kind, wrap)
		}
	}
}

func registerObserveHandler(kind Kind, fn ObserveHandler) {
	if fn == nil {
		return
	}
	wrap := func(ctx context.Context, ev *Event) (model.Result, error) {
		if err := fn(ctx, ev); err != nil {
			return model.Result{}, err
		}
		return model.Result{}, nil
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		for _, eventName := range eventsForKind(agent, kind) {
			registerAgnosticHandler(agent, eventName, kind, wrap)
		}
	}
}

func registerAny(fn ObserveHandler) {
	if fn == nil {
		return
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		run.RegisterAnyHandler("agnostic", agent.String(), makeObserveProducer(agent, fn))
	}
}

func registerAgnosticHandler(agent Dialect, eventName string, kind Kind, fn func(context.Context, *Event) (model.Result, error)) {
	run.RegisterHandler("agnostic", agent.String(), eventName, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, err
		}
		if ev.Kind != kind {
			return nil, 0, nil
		}
		res, err := fn(ctx, ev)
		if err != nil {
			return nil, handlerErrorExit(ev), err
		}
		stdout, code, err := codec.Encode(ev, res)
		return stdout, code, err
	})
}

func makeObserveProducer(agent Dialect, fn ObserveHandler) run.Producer {
	return func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, err
		}
		if err := fn(ctx, ev); err != nil {
			return nil, handlerErrorExit(ev), err
		}
		return nil, 0, nil
	}
}

func codecForServe(agent Dialect, _ *run.Config) (Codec, error) {
	switch agent {
	case Claude:
		return &claude.Codec{}, nil
	case Copilot:
		return &copilot.Codec{}, nil
	case Cursor:
		return &cursor.Codec{}, nil
	default:
		return nil, fmt.Errorf("agnostic: no codec for dialect %q", agent)
	}
}

func eventsForKind(agent Dialect, kind Kind) []string {
	switch agent {
	case Claude:
		if name, ok := ClaudeEventForKind[kind]; ok {
			return []string{name}
		}
	case Copilot:
		if name, ok := CopilotEventForKind[kind]; ok {
			return []string{name}
		}
	case Cursor:
		var out []string
		for event, k := range CursorKindForEventMap {
			if k == kind {
				out = append(out, event)
			}
		}
		return out
	}
	return nil
}

func handlerErrorExit(ev *Event) int {
	switch ev.Agent {
	case Copilot:
		if ev.Kind == KindPreTool {
			return CopilotPreToolErrorExit
		}
		return 1
	case Cursor:
		return CursorHandlerErrorExit
	default:
		return 1
	}
}

// ResetHandlers clears only agnostic-registered handlers from the shared run registry.
func ResetHandlers() {
	run.ResetOwner("agnostic")
}
