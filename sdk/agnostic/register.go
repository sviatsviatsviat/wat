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
type PreToolHandler func(ctx context.Context, hook PreToolHook, results PreToolResults) (PreToolResult, error)

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, hook PostToolHook, results PostToolResults) (PostToolResult, error)

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, hook PostToolFailureHook, results PostToolFailureResults) (PostToolFailureResult, error)

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler func(ctx context.Context, hook StopHook, results StopResults) (StopResult, error)

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, hook SessionStartHook, results SessionStartResults) (SessionStartResult, error)

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler func(ctx context.Context, hook SessionEndHook) error

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler func(ctx context.Context, hook UserPromptHook) error

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler func(ctx context.Context, hook PreCompactHook) error

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler func(ctx context.Context, hook SubagentStartHook) error

// AnyHandler handles every event before kind-specific handlers with no hook response.
type AnyHandler func(ctx context.Context, hook AnyHook) error

// Chain supports fluent handler registration into the shared run registry.
type Chain struct{}

// OnPreTool registers a handler for PreTool events across all agents.
func OnPreTool(fn PreToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPreTool, func(ctx context.Context, ev *Event) (PreToolResult, error) {
		typed, err := PreToolEventFrom(ev)
		if err != nil {
			return PreToolResult{}, err
		}
		return fn(ctx, preToolHook(run.InvocationFrom(ctx), typed), preToolResults{})
	})
	return &Chain{}
}

// OnPostTool registers a handler for PostTool events across all agents.
func OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPostTool, func(ctx context.Context, ev *Event) (PostToolResult, error) {
		typed, err := PostToolEventFrom(ev)
		if err != nil {
			return PostToolResult{}, err
		}
		return fn(ctx, postToolHook(run.InvocationFrom(ctx), typed), postToolResults{})
	})
	return &Chain{}
}

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindPostToolFailure, func(ctx context.Context, ev *Event) (PostToolFailureResult, error) {
		typed, err := PostToolFailureEventFrom(ev)
		if err != nil {
			return PostToolFailureResult{}, err
		}
		return fn(ctx, postToolFailureHook(run.InvocationFrom(ctx), typed), postToolFailureResults{})
	})
	return &Chain{}
}

// OnStop registers a handler for Stop events across all agents.
func OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		typed, err := StopEventFrom(ev)
		if err != nil {
			return StopResult{}, err
		}
		return fn(ctx, stopHook(run.InvocationFrom(ctx), typed), stopResults{})
	})
	return &Chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindSubagentStop, func(ctx context.Context, ev *Event) (StopResult, error) {
		typed, err := StopEventFrom(ev)
		if err != nil {
			return StopResult{}, err
		}
		return fn(ctx, stopHook(run.InvocationFrom(ctx), typed), stopResults{})
	})
	return &Chain{}
}

// OnSessionStart registers a handler for SessionStart events across all agents.
func OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerResultHandler(KindSessionStart, func(ctx context.Context, ev *Event) (SessionStartResult, error) {
		typed, err := SessionStartEventFrom(ev)
		if err != nil {
			return SessionStartResult{}, err
		}
		return fn(ctx, sessionStartHook(run.InvocationFrom(ctx), typed), sessionStartResults{})
	})
	return &Chain{}
}

// OnSessionEnd registers an observe-only handler for SessionEnd events.
func OnSessionEnd(fn SessionEndHandler) *Chain {
	registerObserveHandler(KindSessionEnd, func(ctx context.Context, ev *Event) error {
		typed, err := SessionEndEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, sessionEndHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnUserPrompt registers an observe-only handler for UserPrompt events.
func OnUserPrompt(fn UserPromptHandler) *Chain {
	registerObserveHandler(KindUserPrompt, func(ctx context.Context, ev *Event) error {
		typed, err := UserPromptEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, userPromptHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnPreCompact registers an observe-only handler for PreCompact events.
func OnPreCompact(fn PreCompactHandler) *Chain {
	registerObserveHandler(KindPreCompact, func(ctx context.Context, ev *Event) error {
		typed, err := PreCompactEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, preCompactHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnSubagentStart registers an observe-only handler for SubagentStart events.
func OnSubagentStart(fn SubagentStartHandler) *Chain {
	registerObserveHandler(KindSubagentStart, func(ctx context.Context, ev *Event) error {
		typed, err := SubagentStartEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, subagentStartHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnAny registers an observe-only handler invoked for every event before kind-specific handlers.
func OnAny(fn AnyHandler) *Chain {
	registerAny(fn)
	return &Chain{}
}

// OnPostTool registers another PostTool handler on the chain.
func (c *Chain) OnPostTool(fn PostToolHandler) *Chain {
	return OnPostTool(fn)
}

// OnPostToolFailure registers another PostToolFailure handler on the chain.
func (c *Chain) OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	return OnPostToolFailure(fn)
}

// OnStop registers another Stop handler on the chain.
func (c *Chain) OnStop(fn StopHandler) *Chain {
	return OnStop(fn)
}

// OnSubagentStop registers another SubagentStop handler on the chain.
func (c *Chain) OnSubagentStop(fn StopHandler) *Chain {
	return OnSubagentStop(fn)
}

// OnSessionStart registers another SessionStart handler on the chain.
func (c *Chain) OnSessionStart(fn SessionStartHandler) *Chain {
	return OnSessionStart(fn)
}

// OnSessionEnd registers another observe-only SessionEnd handler on the chain.
func (c *Chain) OnSessionEnd(fn SessionEndHandler) *Chain {
	return OnSessionEnd(fn)
}

// OnUserPrompt registers another observe-only UserPrompt handler on the chain.
func (c *Chain) OnUserPrompt(fn UserPromptHandler) *Chain {
	return OnUserPrompt(fn)
}

// OnPreCompact registers another observe-only PreCompact handler on the chain.
func (c *Chain) OnPreCompact(fn PreCompactHandler) *Chain {
	return OnPreCompact(fn)
}

// OnSubagentStart registers another observe-only SubagentStart handler on the chain.
func (c *Chain) OnSubagentStart(fn SubagentStartHandler) *Chain {
	return OnSubagentStart(fn)
}

// OnAny registers another observe-only catch-all handler on the chain.
func (c *Chain) OnAny(fn AnyHandler) *Chain {
	return OnAny(fn)
}

// OnPreTool registers another PreTool handler on the chain.
func (c *Chain) OnPreTool(fn PreToolHandler) *Chain {
	return OnPreTool(fn)
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

func registerObserveHandler(kind Kind, fn func(context.Context, *Event) error) {
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

func registerAny(fn AnyHandler) {
	if fn == nil {
		return
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		run.RegisterAnyHandler("agnostic", agent.String(), makeAnyProducer(agent, fn))
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

func makeAnyProducer(agent Dialect, fn AnyHandler) run.Producer {
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
		if err := fn(ctx, anyHook(run.InvocationFrom(ctx), AnyEventFrom(ev))); err != nil {
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
