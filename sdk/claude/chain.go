package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcompleted"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcreated"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact/precompact"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitation"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptexpansion"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptsubmit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissiondenied"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/pretooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/messagedisplay"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/notification"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreecreate"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
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
		panic("claude: UseHooks: at most one registry")
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

// UserPromptSubmit registers a UserPromptSubmit handler on the chain.
func (c *chain) UserPromptSubmit(fn func(context.Context, run.Hook[UserPromptSubmit], UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *chain {
	userpromptsubmit.RegisterHandler(c.reg, fn)
	return c
}

// UserPromptExpansion registers a UserPromptExpansion handler on the chain.
func (c *chain) UserPromptExpansion(fn func(context.Context, run.Hook[UserPromptExpansion], UserPromptExpansionResults) (CommonOutput, error)) *chain {
	userpromptexpansion.RegisterHandler(c.reg, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PreToolUseResults) (PreToolUseOutput, error)) *chain {
	pretooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolUseResults) (PostToolUseOutput, error)) *chain {
	posttooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolUseFailureResults) (PostToolUseOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(c.reg, fn)
	return c
}

// PermissionRequest registers a PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, run.Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	permissionrequest.RegisterHandler(c.reg, fn)
	return c
}

// PermissionDenied registers a PermissionDenied handler on the chain.
func (c *chain) PermissionDenied(fn func(context.Context, run.Hook[PermissionDenied], PermissionDeniedResults) (PermissionDeniedOutput, error)) *chain {
	permissiondenied.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (CommonOutput, error)) *chain {
	subagentstart.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(c.reg, fn)
	return c
}

// TaskCreated registers a TaskCreated handler on the chain.
func (c *chain) TaskCreated(fn func(context.Context, run.Hook[TaskCreated], TaskCreatedResults) (CommonOutput, error)) *chain {
	taskcreated.RegisterHandler(c.reg, fn)
	return c
}

// TaskCompleted registers a TaskCompleted handler on the chain.
func (c *chain) TaskCompleted(fn func(context.Context, run.Hook[TaskCompleted], TaskCompletedResults) (CommonOutput, error)) *chain {
	taskcompleted.RegisterHandler(c.reg, fn)
	return c
}

// Stop registers a Stop handler on the chain.
func (c *chain) Stop(fn func(context.Context, run.Hook[Stop], StopResults) (StopOutput, error)) *chain {
	stopevent.RegisterHandler(c.reg, fn)
	return c
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, run.Hook[Notification], NotificationResults) (CommonOutput, error)) *chain {
	notification.RegisterHandler(c.reg, fn)
	return c
}

// MessageDisplay registers a MessageDisplay handler on the chain.
func (c *chain) MessageDisplay(fn func(context.Context, run.Hook[MessageDisplay], MessageDisplayResults) (MessageDisplayOutput, error)) *chain {
	messagedisplay.RegisterHandler(c.reg, fn)
	return c
}

// WorktreeCreate registers a WorktreeCreate handler on the chain.
func (c *chain) WorktreeCreate(fn func(context.Context, run.Hook[WorktreeCreate], WorktreeCreateResults) (WorktreeCreateOutput, error)) *chain {
	worktreecreate.RegisterHandler(c.reg, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact], PreCompactResults) (CommonOutput, error)) *chain {
	precompact.RegisterHandler(c.reg, fn)
	return c
}

// Elicitation registers an Elicitation handler on the chain.
func (c *chain) Elicitation(fn func(context.Context, run.Hook[Elicitation], ElicitationResults) (ElicitationOutput, error)) *chain {
	elicitation.RegisterHandler(c.reg, fn)
	return c
}
