package run

import (
	"sort"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Registry ensures dialect slots and returns handler bags for registration.
// Same name merges across Ensure calls: detect and codec from the first call win;
// later calls return the existing Dialect for attaching more handlers.
//
// Registry is sealed: only sdk/run can provide implementations (via Serve).
type Registry interface {
	Ensure(name string, detect hookkit.DetectFunc, codec *hookkit.Codec) *hookkit.Dialect
	registry()
}

// Hooks contributes one or more dialect handler bags via Registry.
// Values that share a dialect name merge: handlers append in Contribute order.
type Hooks interface {
	Contribute(Registry)
}

// ManifestVersion is the current registration manifest schema version.
const ManifestVersion = 1

// Registration describes one native agent event with authored handlers.
type Registration struct {
	// Dialect is the native agent dialect.
	Dialect string `json:"dialect"`
	// Event is the native hook event name.
	Event string `json:"event"`
	// HandlerCount is the number of handlers registered for the event.
	HandlerCount int `json:"handler_count"`
}

// Manifest describes the native registrations contributed by authored hooks.
type Manifest struct {
	// Version is the manifest schema version.
	Version int `json:"version"`
	// Registrations contains one sorted entry per native dialect and event.
	Registrations []Registration `json:"registrations"`
}

// EventsFor returns the sorted native events registered for dialect.
// Events with HandlerCount < 1 are omitted.
func (m Manifest) EventsFor(dialect string) []string {
	var events []string
	for _, registration := range m.Registrations {
		if registration.Dialect == dialect && registration.HandlerCount > 0 {
			events = append(events, registration.Event)
		}
	}
	sort.Strings(events)
	return events
}

// Has reports whether dialect and event have at least one registered handler.
func (m Manifest) Has(dialect, event string) bool {
	for _, registration := range m.Registrations {
		if registration.Dialect == dialect && registration.Event == event && registration.HandlerCount > 0 {
			return true
		}
	}
	return false
}

// Inspect contributes hooks and returns their native registration manifest.
func Inspect(hooks ...Hooks) Manifest {
	r := newRouter()
	contribute(r, hooks)
	return r.manifest()
}
