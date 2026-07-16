package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hook is the handler context for a typed Copilot hook event.
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

// Raw returns the untouched native JSON payload when available.
func (h Hook[E]) Raw() json.RawMessage { return RawBytes(h.Event) }

// AnyHook is the handler context for catch-all OnAny handlers.
type AnyHook struct {
	Event
	inv run.Invocation
}

// NewAnyHook wraps ev with serve-time invocation settings.
func NewAnyHook(inv run.Invocation, ev Event) AnyHook {
	return AnyHook{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h AnyHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload when available.
func (h AnyHook) Raw() json.RawMessage { return RawBytes(h.Event) }
