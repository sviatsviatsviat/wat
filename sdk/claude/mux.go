package claude

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
		cfg:      defaultRuntimeConfig(),
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
		panic(fmt.Sprintf("claude: duplicate handler for %s", name))
	}
	m.handlers[name] = registeredHandler{
		eventName: name,
		fn: func(ctx context.Context, ev Event) (any, error) {
			typed, ok := ev.(E)
			if !ok {
				return nil, fmt.Errorf("claude: handler for %s received %T", name, ev)
			}
			return fn(ctx, typed)
		},
	}
}

// Serve reads stdin, dispatches the registered handler, writes stdout, and returns the exit code.
func (m *Mux) Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...Option) int {
	if m == nil {
		_, _ = fmt.Fprintln(errw, "claude: nil mux")
		return HandlerErrorExit
	}
	cfg := m.cfg
	applyOptions(&cfg, opts...)

	return hookkit.ServeLoop(ctx, in, out, errw, hookkit.ServeHooks{
		Label:            "claude",
		HandlerErrorExit: HandlerErrorExit,
		Decode: func(raw []byte) (string, any, error) {
			ev, err := Decode(raw)
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
		OnHandlerError: func(string, error) int {
			if cfg.policy == FailBlock {
				return FailBlockExit
			}
			return HandlerErrorExit
		},
		Encode: func(eventName string, result any) ([]byte, int, error) {
			stdout, err := Encode(eventName, result, cfg.encodeOpts()...)
			return stdout, 0, err
		},
	})
}

// Main runs Serve with os.Stdin, os.Stdout, and os.Stderr, then os.Exit.
func (m *Mux) Main(opts ...Option) {
	os.Exit(m.Serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, opts...))
}
