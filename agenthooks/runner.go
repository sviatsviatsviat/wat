package agenthooks

import (
	"context"
	"fmt"
	"io"
	"os"
)

// HandlerFunc handles a decoded hook event and returns a unified Result.
type HandlerFunc func(ctx context.Context, ev *Event) (Result, error)

// Mux registers hook handlers and runs the hook process lifecycle.
type Mux struct {
	anyHandlers  []HandlerFunc
	kindHandlers map[Kind][]HandlerFunc
}

// NewMux returns an empty handler mux.
func NewMux() *Mux {
	return &Mux{
		kindHandlers: make(map[Kind][]HandlerFunc),
	}
}

// On registers a handler for a normalized event kind.
func (m *Mux) On(kind Kind, fn HandlerFunc) {
	if fn == nil {
		return
	}
	m.kindHandlers[kind] = append(m.kindHandlers[kind], fn)
}

// OnAny registers a handler invoked for every event before kind-specific handlers.
func (m *Mux) OnAny(fn HandlerFunc) {
	if fn == nil {
		return
	}
	m.anyHandlers = append(m.anyHandlers, fn)
}

// Option configures Serve and Main.
type Option func(*serveConfig)

type serveConfig struct {
	dialect   Dialect
	eventHint string
	getenv    func(string) string
}

// WithDialect forces a dialect instead of Detect.
func WithDialect(d Dialect) Option {
	return func(c *serveConfig) {
		c.dialect = d
	}
}

// WithEvent supplies the native event name for Copilot camelCase payloads that
// omit hook_event_name. It is passed to Codec.Decode as eventHint.
func WithEvent(name string) Option {
	return func(c *serveConfig) {
		c.eventHint = name
	}
}

// WithGetenv injects environment lookup for Detect and ClaudeCodec encode.
func WithGetenv(getenv func(string) string) Option {
	return func(c *serveConfig) {
		c.getenv = getenv
	}
}

// Serve reads a hook payload from in, dispatches registered handlers, writes
// encoded stdout to out, diagnostics to errw, and returns the process exit code.
func (m *Mux) Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...Option) int {
	cfg := serveConfig{getenv: os.Getenv}
	for _, opt := range opts {
		opt(&cfg)
	}

	payload, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "agenthooks: read stdin: %v\n", err)
		return 1
	}
	if len(payload) == 0 {
		_, _ = fmt.Fprintln(errw, "agenthooks: empty stdin")
		return 1
	}

	dialect := cfg.dialect
	if dialect == Unknown {
		dialect = Detect(payload, cfg.getenv)
	}
	if dialect == Unknown {
		_, _ = fmt.Fprintln(errw, "agenthooks: unknown dialect")
		return 1
	}

	codec, err := CodecFor(dialect)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "agenthooks: %v\n", err)
		return 1
	}
	if cc, ok := codec.(*ClaudeCodec); ok && cfg.getenv != nil {
		cc.Getenv = cfg.getenv
	}

	ev, err := codec.Decode(payload, cfg.eventHint)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "agenthooks: decode: %v\n", err)
		return 1
	}

	merged, exit, ok := m.dispatch(ctx, ev, errw)
	if !ok {
		return exit
	}

	stdout, code, err := codec.Encode(ev, merged)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "agenthooks: encode: %v\n", err)
		return 1
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "agenthooks: write stdout: %v\n", err)
			return 1
		}
	}
	return code
}

func (m *Mux) dispatch(ctx context.Context, ev *Event, errw io.Writer) (Result, int, bool) {
	var merged Result
	for _, fn := range m.anyHandlers {
		res, err := fn(ctx, ev)
		if err != nil {
			_, _ = fmt.Fprintln(errw, err.Error())
			return Result{}, handlerErrorExit(ev, err), false
		}
		merged = Merge(merged, res)
	}
	for _, fn := range m.kindHandlers[ev.Kind] {
		res, err := fn(ctx, ev)
		if err != nil {
			_, _ = fmt.Fprintln(errw, err.Error())
			return Result{}, handlerErrorExit(ev, err), false
		}
		merged = Merge(merged, res)
	}
	return merged, 0, true
}

func handlerErrorExit(ev *Event, _ error) int {
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

// Main runs Serve with os.Stdin, os.Stdout, and os.Stderr, then os.Exit.
func (m *Mux) Main(opts ...Option) {
	os.Exit(m.Serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, opts...))
}
