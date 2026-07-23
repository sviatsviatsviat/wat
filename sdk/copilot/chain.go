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

// chain supports fluent handler registration on the package default dialect.
// Callers obtain a chain via UseHooks.
type chain struct{}

// UseHooks returns a fluent registrar bound to this package's default dialect.
func UseHooks() *chain {
	return &chain{}
}

// SessionStart registers a SessionStart handler on the chain.
func (c *chain) SessionStart(fn func(context.Context, run.Hook[SessionStart], SessionStartResults) (SessionStartOutput, error)) *chain {
	sessionstart.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SessionEnd registers a SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, run.Hook[SessionEnd]) error) *chain {
	sessionend.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// UserPromptSubmitted registers a UserPromptSubmitted handler on the chain.
func (c *chain) UserPromptSubmitted(fn func(context.Context, run.Hook[UserPromptSubmitted]) error) *chain {
	userpromptsubmitted.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PreToolResults) (PreToolOutput, error)) *chain {
	pretooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	posttooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolFailureResults) (PostToolFailureOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PermissionRequest registers a PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, run.Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	permissionrequest.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (SubagentStartOutput, error)) *chain {
	subagentstart.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AgentStop registers an AgentStop handler on the chain.
func (c *chain) AgentStop(fn func(context.Context, run.Hook[AgentStop], StopResults) (StopOutput, error)) *chain {
	agentstop.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact]) error) *chain {
	precompact.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, run.Hook[Notification], NotificationResults) (NotificationOutput, error)) *chain {
	notification.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// ErrorOccurred registers an ErrorOccurred handler on the chain.
func (c *chain) ErrorOccurred(fn func(context.Context, run.Hook[ErrorOccurred]) error) *chain {
	erroroccurred.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}
