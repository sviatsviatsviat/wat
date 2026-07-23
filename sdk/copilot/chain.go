package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/compact/precompact"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/errors/erroroccurred"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/prompt/userpromptsubmitted"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/pretooluse"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/ui/notification"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// chain supports fluent handler registration into a run.Registry.
// Callers obtain a chain via UseHooks.
type chain struct {
	reg *run.Registry
}

var defaultChain = &chain{reg: runtime.DefaultReg}

// UseHooks returns a fluent registrar. With no arguments (or a nil / default
// registry) it returns the package default chain; otherwise it attaches this
// dialect to regs[0] and returns a new chain.
func UseHooks(regs ...*run.Registry) *chain {
	switch len(regs) {
	case 0:
		return defaultChain
	case 1:
		if regs[0] == nil || regs[0] == runtime.DefaultReg {
			return defaultChain
		}
		runtime.EnsureDialect(regs[0])
		return &chain{reg: regs[0]}
	default:
		panic("copilot: UseHooks: at most one registry")
	}
}

// SessionStart registers a SessionStart handler on the chain.
func (c *chain) SessionStart(fn func(context.Context, run.Hook[SessionStart], SessionStartResults) (SessionStartOutput, error)) *chain {
	sessionstart.RegisterHandler(c.reg, fn)
	return c
}

// SessionEnd registers a SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, run.Hook[SessionEnd]) error) *chain {
	sessionend.RegisterHandler(c.reg, fn)
	return c
}

// UserPromptSubmitted registers a UserPromptSubmitted handler on the chain.
func (c *chain) UserPromptSubmitted(fn func(context.Context, run.Hook[UserPromptSubmitted]) error) *chain {
	userpromptsubmitted.RegisterHandler(c.reg, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PreToolResults) (PreToolOutput, error)) *chain {
	pretooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	posttooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolFailureResults) (PostToolFailureOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(c.reg, fn)
	return c
}

// PermissionRequest registers a PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, run.Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	permissionrequest.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (SubagentStartOutput, error)) *chain {
	subagentstart.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(c.reg, fn)
	return c
}

// AgentStop registers an AgentStop handler on the chain.
func (c *chain) AgentStop(fn func(context.Context, run.Hook[AgentStop], StopResults) (StopOutput, error)) *chain {
	agentstop.RegisterHandler(c.reg, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact]) error) *chain {
	precompact.RegisterHandler(c.reg, fn)
	return c
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, run.Hook[Notification], NotificationResults) (NotificationOutput, error)) *chain {
	notification.RegisterHandler(c.reg, fn)
	return c
}

// ErrorOccurred registers an ErrorOccurred handler on the chain.
func (c *chain) ErrorOccurred(fn func(context.Context, run.Hook[ErrorOccurred]) error) *chain {
	erroroccurred.RegisterHandler(c.reg, fn)
	return c
}
