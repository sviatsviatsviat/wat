package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hook is the handler context for a typed Claude hook event.
type Hook[E Event] struct {
	// Event is the decoded native hook payload.
	Event E
	inv   run.Invocation
}

// NewHook wraps ev with serve-time invocation settings.
func NewHook[E Event](inv run.Invocation, ev E) Hook[E] {
	return Hook[E]{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h Hook[E]) Invocation() run.Invocation { return h.inv }
