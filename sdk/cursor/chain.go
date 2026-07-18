package cursor

// chain supports fluent handler registration into the shared run registry.
// Callers obtain a chain only via package-level On* helpers (and further methods).
type chain struct{}
