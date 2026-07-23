package agnostic

// chain supports fluent portable handler registration.
// Callers obtain a chain via UseHooks.
type chain struct{}

var defaultChain = &chain{}

// UseHooks returns a fluent registrar that fans out onto each agent SDK.
func UseHooks() *chain {
	return defaultChain
}
