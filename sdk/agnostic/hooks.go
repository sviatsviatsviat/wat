package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// hooks collects deferred portable handler registrations until Contribute
// installs them (same dialect name appends handlers).
// Callers obtain a registrar via UseHooks.
type hooks struct {
	parts []run.Hooks
}

// UseHooks returns a fluent registrar that fans out onto each agent SDK.
func UseHooks() *hooks {
	return &hooks{}
}

func (c *hooks) appendParts(parts ...run.Hooks) *hooks {
	for _, p := range parts {
		if p != nil {
			c.parts = append(c.parts, p)
		}
	}
	return c
}

// Contribute installs these hooks' dialect contributions via reg.
func (c *hooks) Contribute(reg run.Registry) {
	if c == nil || reg == nil {
		return
	}
	for _, p := range c.parts {
		p.Contribute(reg)
	}
}
