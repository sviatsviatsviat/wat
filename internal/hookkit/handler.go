package hookkit

import (
	"context"
	"fmt"
)

type handler[E Event, O Output] struct {
	fn func(context.Context, E) (O, error)
}

// EventName returns the native event name for E.
func (h handler[E, O]) EventName() string {
	var zero E
	return zero.EventName()
}

// Invoke runs the typed handler.
func (h handler[E, O]) Invoke(ctx context.Context, event Event) (Output, error) {
	typed, ok := event.(E)
	if !ok {
		panic(fmt.Sprintf("handler for %s received %T", h.EventName(), event))
	}
	result, err := h.fn(ctx, typed)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type observeHandler[E Event] struct {
	fn func(context.Context, E) error
}

// EventName returns the native event name for E.
func (h observeHandler[E]) EventName() string {
	var zero E
	return zero.EventName()
}

// Invoke runs the typed observe handler.
func (h observeHandler[E]) Invoke(ctx context.Context, event Event) (Output, error) {
	typed, ok := event.(E)
	if !ok {
		panic(fmt.Sprintf("handler for %s received %T", h.EventName(), event))
	}
	if err := h.fn(ctx, typed); err != nil {
		return nil, err
	}
	return nil, nil
}

// Handler returns a HookHandler that produces a typed output for event E.
// A nil fn yields a nil HookHandler (Dialect.Register no-ops).
func Handler[E Event, O Output](fn func(context.Context, E) (O, error)) HookHandler {
	if fn == nil {
		return nil
	}
	return handler[E, O]{fn: fn}
}

// ObserveHandler returns a HookHandler that observes event E with no host JSON.
// A nil fn yields a nil HookHandler (Dialect.Register no-ops).
func ObserveHandler[E Event](fn func(context.Context, E) error) HookHandler {
	if fn == nil {
		return nil
	}
	return observeHandler[E]{fn: fn}
}
