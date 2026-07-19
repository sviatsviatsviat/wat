package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// chain supports fluent handler registration into a run.Registry.
// Callers obtain a chain via UseHooks.
type chain struct {
	reg *run.Registry
}

func newChain(r *run.Registry) *chain {
	return &chain{reg: r}
}

// UseHooks returns a fluent registrar. With no arguments (or a nil / default
// registry) it returns the package default chain; otherwise it attaches this
// dialect to regs[0] and returns a new chain.
func UseHooks(regs ...*run.Registry) *chain {
	switch len(regs) {
	case 0:
		return defaultChain
	case 1:
		if regs[0] == nil || regs[0] == defaultReg {
			return defaultChain
		}
		ensureDialect(regs[0])
		return newChain(regs[0])
	default:
		panic("claude: UseHooks: at most one registry")
	}
}
