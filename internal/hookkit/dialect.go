package hookkit

import (
	"context"
)

// HookHandler is a typed hook registration ready for dialect dispatch.
type HookHandler interface {
	EventName() string
	Invoke(ctx context.Context, event Event) (Output, error)
}

// Dialect holds the codec and handlers for one agent.
// Registration must finish before serve; not safe for concurrent use.
type Dialect struct {
	codec    *Codec
	handlers map[string][]HookHandler
}

// NewDialect returns a dialect that decodes with c.
func NewDialect(c *Codec) *Dialect {
	if c == nil {
		panic("hookkit: NewDialect: nil codec")
	}
	return &Dialect{
		codec:    c,
		handlers: make(map[string][]HookHandler),
	}
}

// Register appends a hook handler on d.
func (d *Dialect) Register(h HookHandler) {
	if h == nil {
		return
	}
	eventName := h.EventName()
	d.handlers[eventName] = append(d.handlers[eventName], h)
}

// RegisterWith appends a handler that receives results when invoked.
// Pass results as type R (typically the package Results interface), matching fn.
// A nil fn is a no-op.
func RegisterWith[E Event, O Output, R any](d *Dialect, results R, fn func(context.Context, Hook[E], R) (O, error)) {
	if fn == nil {
		return
	}
	d.Register(Handler(func(ctx context.Context, hook Hook[E]) (O, error) {
		return fn(ctx, hook, results)
	}))
}

// RegisterObserve appends an observe-only handler for event E.
// A nil fn is a no-op.
func RegisterObserve[E Event](d *Dialect, fn func(context.Context, Hook[E]) error) {
	d.Register(ObserveHandler(fn))
}

// Codec returns the dialect codec.
func (d *Dialect) Codec() *Codec {
	return d.codec
}

// HandlersFor returns a copy of handlers registered for eventName.
func (d *Dialect) HandlersFor(eventName string) []HookHandler {
	return append([]HookHandler(nil), d.handlers[eventName]...)
}
