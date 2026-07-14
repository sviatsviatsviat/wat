package copilot

import (
	"context"
)

// Chain supports fluent handler registration into the shared run registry.
type Chain struct{}

// OnAny registers an observe-only handler invoked for every event.
func (c *Chain) OnAny(fn func(context.Context, AnyHook) error) *Chain {
	registerAny(fn)
	return &Chain{}
}
