package agnostic

// chain supports fluent handler registration into the shared run registry.
// Callers obtain a chain only via package-level On* helpers (and further On* methods).
type chain struct{}
