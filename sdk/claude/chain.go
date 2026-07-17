package claude

import (
	"context"
)

// adapterOwner is the run registry owner used when agnostic fans out onto this SDK.
const adapterOwner = "agnostic"

// Chain supports fluent handler registration into the shared run registry.
type Chain struct {
	owner string // empty ⇒ "claude"
}

func (c *Chain) registerOwner() string {
	if c != nil && c.owner != "" {
		return c.owner
	}
	return "claude"
}

// Adapter returns a Chain that registers handlers under run owner "agnostic".
// Used by sdk/agnostic On* fan-out so portable and native handlers share one
// dialect registry but remain independently resettable in tests.
func Adapter() *Chain {
	return &Chain{owner: adapterOwner}
}

// OnAny registers an observe-only handler invoked for every event.
func (c *Chain) OnAny(fn func(context.Context, AnyHook) error) *Chain {
	if c == nil {
		c = &Chain{}
	}
	registerAny(c.registerOwner(), fn)
	return c
}
