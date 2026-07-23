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
)

// chain supports fluent handler registration on the package default dialect.
// Callers obtain a chain via UseHooks.
type chain struct{}

// UseHooks returns a fluent registrar bound to this package's default dialect.
func UseHooks() *chain {
	return &chain{}
}

// SessionStart registers a SessionStart handler on the chain.
func (c *chain) SessionStart(fn func(context.Context, SessionStart, SessionStartResults) (SessionStartOutput, error)) *chain {
	sessionstart.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SessionEnd registers a SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, SessionEnd) error) *chain {
	sessionend.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// UserPromptSubmit registers a UserPromptSubmit handler on the chain.
func (c *chain) UserPromptSubmit(fn func(context.Context, UserPromptSubmit, UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *chain {
	userpromptsubmit.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// UserPromptExpansion registers a UserPromptExpansion handler on the chain.
func (c *chain) UserPromptExpansion(fn func(context.Context, UserPromptExpansion, UserPromptExpansionResults) (CommonOutput, error)) *chain {
	userpromptexpansion.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, PreToolUse, PreToolUseResults) (PreToolUseOutput, error)) *chain {
	pretooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, PostToolUse, PostToolUseResults) (PostToolUseOutput, error)) *chain {
	posttooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailure, PostToolUseFailureResults) (PostToolUseOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PermissionRequest registers a PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, PermissionRequest, PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	permissionrequest.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PermissionDenied registers a PermissionDenied handler on the chain.
func (c *chain) PermissionDenied(fn func(context.Context, PermissionDenied, PermissionDeniedResults) (PermissionDeniedOutput, error)) *chain {
	permissiondenied.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, SubagentStart, SubagentStartResults) (CommonOutput, error)) *chain {
	subagentstart.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, SubagentStop, StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// TaskCreated registers a TaskCreated handler on the chain.
func (c *chain) TaskCreated(fn func(context.Context, TaskCreated, TaskCreatedResults) (CommonOutput, error)) *chain {
	taskcreated.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// TaskCompleted registers a TaskCompleted handler on the chain.
func (c *chain) TaskCompleted(fn func(context.Context, TaskCompleted, TaskCompletedResults) (CommonOutput, error)) *chain {
	taskcompleted.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// Stop registers a Stop handler on the chain.
func (c *chain) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *chain {
	stopevent.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, Notification, NotificationResults) (CommonOutput, error)) *chain {
	notification.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// MessageDisplay registers a MessageDisplay handler on the chain.
func (c *chain) MessageDisplay(fn func(context.Context, MessageDisplay, MessageDisplayResults) (MessageDisplayOutput, error)) *chain {
	messagedisplay.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// WorktreeCreate registers a WorktreeCreate handler on the chain.
func (c *chain) WorktreeCreate(fn func(context.Context, WorktreeCreate, WorktreeCreateResults) (WorktreeCreateOutput, error)) *chain {
	worktreecreate.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, PreCompact, PreCompactResults) (CommonOutput, error)) *chain {
	precompact.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// Elicitation registers an Elicitation handler on the chain.
func (c *chain) Elicitation(fn func(context.Context, Elicitation, ElicitationResults) (ElicitationOutput, error)) *chain {
	elicitation.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}
