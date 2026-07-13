package claudehook

import (
	"context"
	"fmt"
	"io"
	"os"
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
	m.handlers[name] = registeredHandler{
		eventName: name,
		fn: func(ctx context.Context, ev Event) (any, error) {
			typed, ok := ev.(E)
			if !ok {
				return nil, fmt.Errorf("claudehook: handler for %s received %T", name, ev)
			}
			return fn(ctx, typed)
		},
	}
}

// Serve reads stdin, dispatches the registered handler, writes stdout, and returns the exit code.
func (m *Mux) Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...Option) int {
	if m == nil {
		_, _ = fmt.Fprintln(errw, "claudehook: nil mux")
		return 1
	}
	cfg := m.cfg
	applyOptions(&cfg, opts...)

	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "claudehook: read stdin: %v\n", err)
		return 1
	}
	ev, err := Decode(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "claudehook: decode: %v\n", err)
		return 1
	}
	h, ok := m.handlers[ev.EventName()]
	if !ok {
		return 0
	}
	result, err := h.fn(ctx, ev)
	if err != nil {
		_, _ = fmt.Fprintln(errw, err.Error())
		if cfg.policy == FailBlock {
			return 2
		}
		return 1
	}
	if result == nil {
		return 0
	}
	if z, ok := result.(interface{ isZero() bool }); ok && z.isZero() {
		return 0
	}
	stdout, err := Encode(ev.EventName(), result, cfg.encodeOpts()...)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "claudehook: encode: %v\n", err)
		return 1
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "claudehook: write stdout: %v\n", err)
			return 1
		}
	}
	return 0
}

// Main runs Serve with os.Stdin, os.Stdout, and os.Stderr, then os.Exit.
func (m *Mux) Main(opts ...Option) {
	os.Exit(m.Serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, opts...))
}
