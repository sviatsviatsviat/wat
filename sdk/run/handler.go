package run

import (
	"context"
	"fmt"
)

// hookRegistration is a typed hook handler ready to attach to a Registry.
type hookRegistration interface {
	eventName() string
	handle(ctx context.Context, event Event) (output []byte, exit int, err error)
}

type handler[E Event, O Output] struct {
	fn func(context.Context, Hook[E]) (O, error)
}

func (h handler[E, O]) eventName() string {
	var zero E
	return zero.EventName()
}

func (h handler[E, O]) handle(ctx context.Context, event Event) ([]byte, int, error) {
	typed, ok := event.(E)
	if !ok {
		panic(fmt.Sprintf("handler for %s received %T", h.eventName(), event))
	}
	result, err := h.fn(ctx, NewHook(InvocationFrom(ctx), typed))
	if err != nil {
		return nil, 1, err
	}
	return encodeOutput(result)
}

type observeHandler[E Event] struct {
	fn func(context.Context, Hook[E]) error
}

func (h observeHandler[E]) eventName() string {
	var zero E
	return zero.EventName()
}

func (h observeHandler[E]) handle(ctx context.Context, event Event) ([]byte, int, error) {
	typed, ok := event.(E)
	if !ok {
		panic(fmt.Sprintf("handler for %s received %T", h.eventName(), event))
	}
	if err := h.fn(ctx, NewHook(InvocationFrom(ctx), typed)); err != nil {
		return nil, 1, err
	}
	return nil, 0, nil
}

func encodeOutput(out Output) ([]byte, int, error) {
	if out.IsZero() {
		return nil, 0, nil
	}
	return out.Encode()
}

// Handler returns a hookRegistration that encodes a typed output for event E.
// A nil fn yields a nil hookRegistration (RegisterHandler no-ops).
func Handler[E Event, O Output](fn func(context.Context, Hook[E]) (O, error)) hookRegistration {
	if fn == nil {
		return nil
	}
	return handler[E, O]{fn: fn}
}

// ObserveHandler returns a hookRegistration that observes event E with no host JSON.
// A nil fn yields a nil hookRegistration (RegisterObserveHandler no-ops).
func ObserveHandler[E Event](fn func(context.Context, Hook[E]) error) hookRegistration {
	if fn == nil {
		return nil
	}
	return observeHandler[E]{fn: fn}
}
