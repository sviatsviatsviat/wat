package run

import "github.com/sviatsviatsviat/wat/internal/hookkit"

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
