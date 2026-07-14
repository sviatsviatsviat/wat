package agnostic

import (
	"context"
	"fmt"
	"io"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// HandlerFunc handles a decoded hook event and returns a unified Result.
type HandlerFunc func(ctx context.Context, ev *Event) (Result, error)

// Chain supports fluent handler registration into the shared run registry.
type Chain struct{}

// On registers a handler for a normalized event kind across all agents.
func On(kind Kind, fn HandlerFunc) *Chain {
	registerKind(kind, fn)
	return &Chain{}
}

// OnAny registers a handler invoked for every event before kind-specific handlers.
func OnAny(fn HandlerFunc) *Chain {
	registerAny(fn)
	return &Chain{}
}

// On registers another kind handler on the chain.
func (c *Chain) On(kind Kind, fn HandlerFunc) *Chain {
	registerKind(kind, fn)
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

// WithGetenv injects environment lookup for Detect and ClaudeCodec encode.
func WithGetenv(getenv func(string) string) run.Option {
	return run.WithGetenv(getenv)
}

func registerKind(kind Kind, fn HandlerFunc) {
	if fn == nil {
		return
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		for _, eventName := range eventsForKind(agent, kind) {
			registerAgnosticHandler(agent, eventName, kind, fn)
		}
	}
}

func registerAny(fn HandlerFunc) {
	if fn == nil {
		return
	}
	for _, agent := range []Dialect{Claude, Copilot, Cursor} {
		run.RegisterAnyHandler("agnostic", agent.String(), makeAgnosticProducer(agent, fn))
	}
}

func registerAgnosticHandler(agent Dialect, eventName string, kind Kind, fn HandlerFunc) {
	run.RegisterHandler("agnostic", agent.String(), eventName, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, fmt.Errorf("agnostic: decode: %w", err)
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

func makeAgnosticProducer(agent Dialect, fn HandlerFunc) run.Producer {
	return func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, fmt.Errorf("agnostic: decode: %w", err)
		}
		res, err := fn(ctx, ev)
		if err != nil {
			return nil, handlerErrorExit(ev), err
		}
		stdout, code, err := codec.Encode(ev, res)
		return stdout, code, err
	}
}

func codecForServe(agent Dialect, cfg *run.Config) (Codec, error) {
	codec, err := CodecFor(agent)
	if err != nil {
		return nil, err
	}
	if cc, ok := codec.(*claude.Codec); ok {
		claude.ApplyRunConfig(cc, cfg)
	}
	return codec, nil
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
