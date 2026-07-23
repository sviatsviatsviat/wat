package hookkit

// Hook is the handler context for a typed hook event.
type Hook[E Event] struct {
	// Event is the decoded native hook payload.
	Event E
	inv   Invocation
}

// NewHook wraps ev with serve-time invocation settings.
func NewHook[E Event](inv Invocation, ev E) Hook[E] {
	return Hook[E]{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h Hook[E]) Invocation() Invocation { return h.inv }
