package copilot

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

// Mux registers typed hook handlers and runs the hook process lifecycle.
type Mux struct {
	handlers map[string]registeredHandler
	cfg      runtimeConfig
}

type registeredHandler struct {
	fn        func(context.Context, Event) (any, error)
	eventName string
}

// NewMux returns an empty handler mux.
func NewMux() *Mux {
	return &Mux{
		handlers: make(map[string]registeredHandler),
		cfg:      defaultDecodeConfig(),
	}
}

// On registers a typed handler for event type E.
func On[E Event, O any](m *Mux, fn func(context.Context, E) (O, error)) {
	if m == nil || fn == nil {
		return
	}
	var zero E
	name := zero.EventName()
	if _, exists := m.handlers[name]; exists {
		panic(fmt.Sprintf("copilot: duplicate handler for %s", name))
	}
	m.handlers[name] = registeredHandler{
		eventName: name,
		fn: func(ctx context.Context, ev Event) (any, error) {
			typed, ok := ev.(E)
			if !ok {
				return nil, fmt.Errorf("copilot: handler for %s received %T", name, ev)
			}
			return fn(ctx, typed)
		},
	}
}

// Serve reads stdin, dispatches the registered handler, writes stdout, and returns the exit code.
func (m *Mux) Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...Option) int {
	if m == nil {
		_, _ = fmt.Fprintln(errw, "copilot: nil mux")
		return HandlerErrorExit
	}
	cfg := m.cfg
	applyOptions(&cfg, opts...)

	return hookkit.ServeLoop(ctx, in, out, errw, hookkit.ServeHooks{
		Label:            "copilot",
		HandlerErrorExit: HandlerErrorExit,
		Decode: func(raw []byte) (string, any, error) {
			ev, err := Decode(raw, WithEvent(cfg.eventHint.Hint))
			if err != nil {
				return "", nil, err
			}
			return ev.EventName(), ev, nil
		},
		Lookup: func(eventName string) (func(context.Context, any) (any, error), bool) {
			h, ok := m.handlers[eventName]
			if !ok {
				return nil, false
			}
			return func(ctx context.Context, ev any) (any, error) {
				return h.fn(ctx, ev.(Event))
			}, true
		},
		IsZeroResult: isZeroOutput,
		OnHandlerError: func(eventName string, err error) int {
			if eventName == EventPreToolUse {
				return PreToolErrorExit
			}
			return HandlerErrorExit
		},
		Encode: func(eventName string, result any) ([]byte, int, error) {
			return Encode(eventName, result)
		},
	})
}

// Main runs Serve with os.Stdin, os.Stdout, and os.Stderr, then os.Exit.
func (m *Mux) Main(opts ...Option) {
	os.Exit(m.Serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, opts...))
}
