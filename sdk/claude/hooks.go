package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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

// hooks is a typed Claude registrar; callers obtain one via UseHooks.
type hooks struct {
	hookkit.HandlerQueue
}

var codec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

// UseHooks returns a fluent registrar for this package's Claude dialect.
func UseHooks() *hooks {
	return &hooks{}
}

func bind[F any](c *hooks, fn F, register func(*hookkit.Dialect, F)) *hooks {
	if c == nil {
		return nil
	}
	hookkit.Bind(&c.HandlerQueue, fn, register)
	return c
}

// Contribute installs the Claude dialect and registers these hooks' handlers.
func (c *hooks) Contribute(reg run.Registry) {
	if c == nil {
		return
	}
	c.Install(reg, runtime.Dialect, detectPayload, codec)
}

func detectPayload(raw []byte) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return false
	}
	if has("hook_event_name") && has("timestamp") {
		return false
	}
	return has("session_id")
}

// SessionStart registers a SessionStart handler on the registrar.
func (c *hooks) SessionStart(fn func(context.Context, SessionStart, SessionStartResults) (SessionStartOutput, error)) *hooks {
	return bind(c, fn, sessionstart.RegisterHandler)
}

// SessionEnd registers a SessionEnd handler on the registrar.
func (c *hooks) SessionEnd(fn func(context.Context, SessionEnd) error) *hooks {
	return bind(c, fn, sessionend.RegisterHandler)
}

// UserPromptSubmit registers a UserPromptSubmit handler on the registrar.
func (c *hooks) UserPromptSubmit(fn func(context.Context, UserPromptSubmit, UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *hooks {
	return bind(c, fn, userpromptsubmit.RegisterHandler)
}

// UserPromptExpansion registers a UserPromptExpansion handler on the registrar.
func (c *hooks) UserPromptExpansion(fn func(context.Context, UserPromptExpansion, UserPromptExpansionResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, userpromptexpansion.RegisterHandler)
}

// PreToolUse registers a PreToolUse handler on the registrar.
func (c *hooks) PreToolUse(fn func(context.Context, PreToolUse, PreToolUseResults) (PreToolUseOutput, error)) *hooks {
	return bind(c, fn, pretooluse.RegisterHandler)
}

// PostToolUse registers a PostToolUse handler on the registrar.
func (c *hooks) PostToolUse(fn func(context.Context, PostToolUse, PostToolUseResults) (PostToolUseOutput, error)) *hooks {
	return bind(c, fn, posttooluse.RegisterHandler)
}

// PostToolUseFailure registers a PostToolUseFailure handler on the registrar.
func (c *hooks) PostToolUseFailure(fn func(context.Context, PostToolUseFailure, PostToolUseFailureResults) (PostToolUseOutput, error)) *hooks {
	return bind(c, fn, posttoolusefailure.RegisterHandler)
}

// PermissionRequest registers a PermissionRequest handler on the registrar.
func (c *hooks) PermissionRequest(fn func(context.Context, PermissionRequest, PermissionRequestResults) (PermissionRequestOutput, error)) *hooks {
	return bind(c, fn, permissionrequest.RegisterHandler)
}

// PermissionDenied registers a PermissionDenied handler on the registrar.
func (c *hooks) PermissionDenied(fn func(context.Context, PermissionDenied, PermissionDeniedResults) (PermissionDeniedOutput, error)) *hooks {
	return bind(c, fn, permissiondenied.RegisterHandler)
}

// SubagentStart registers a SubagentStart handler on the registrar.
func (c *hooks) SubagentStart(fn func(context.Context, SubagentStart, SubagentStartResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, subagentstart.RegisterHandler)
}

// SubagentStop registers a SubagentStop handler on the registrar.
func (c *hooks) SubagentStop(fn func(context.Context, SubagentStop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, subagentstop.RegisterHandler)
}

// TaskCreated registers a TaskCreated handler on the registrar.
func (c *hooks) TaskCreated(fn func(context.Context, TaskCreated, TaskCreatedResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, taskcreated.RegisterHandler)
}

// TaskCompleted registers a TaskCompleted handler on the registrar.
func (c *hooks) TaskCompleted(fn func(context.Context, TaskCompleted, TaskCompletedResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, taskcompleted.RegisterHandler)
}

// Stop registers a Stop handler on the registrar.
func (c *hooks) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, stopevent.RegisterHandler)
}

// Notification registers a Notification handler on the registrar.
func (c *hooks) Notification(fn func(context.Context, Notification, NotificationResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, notification.RegisterHandler)
}

// MessageDisplay registers a MessageDisplay handler on the registrar.
func (c *hooks) MessageDisplay(fn func(context.Context, MessageDisplay, MessageDisplayResults) (MessageDisplayOutput, error)) *hooks {
	return bind(c, fn, messagedisplay.RegisterHandler)
}

// WorktreeCreate registers a WorktreeCreate handler on the registrar.
// Path results encode as plain stdout for command hooks; see docs/agents/claude.md.
func (c *hooks) WorktreeCreate(fn func(context.Context, WorktreeCreate, WorktreeCreateResults) (WorktreeCreateOutput, error)) *hooks {
	return bind(c, fn, worktreecreate.RegisterHandler)
}

// PreCompact registers a PreCompact handler on the registrar.
func (c *hooks) PreCompact(fn func(context.Context, PreCompact, PreCompactResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, precompact.RegisterHandler)
}

// Elicitation registers an Elicitation handler on the registrar.
func (c *hooks) Elicitation(fn func(context.Context, Elicitation, ElicitationResults) (ElicitationOutput, error)) *hooks {
	return bind(c, fn, elicitation.RegisterHandler)
}
