package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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

// hooks is a typed Copilot registrar; callers obtain one via UseHooks.
type hooks struct {
	hookkit.HandlerQueue
}

var codec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

// UseHooks returns a fluent registrar for this package's Copilot dialect.
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

// Contribute installs the Copilot dialect and registers these hooks' handlers.
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
	return has("hook_event_name") && has("timestamp")
}

// SessionStart registers a SessionStart handler on the registrar.
func (c *hooks) SessionStart(fn func(context.Context, SessionStart, SessionStartResults) (SessionStartOutput, error)) *hooks {
	return bind(c, fn, sessionstart.RegisterHandler)
}

// SessionEnd registers a SessionEnd handler on the registrar.
func (c *hooks) SessionEnd(fn func(context.Context, SessionEnd) error) *hooks {
	return bind(c, fn, sessionend.RegisterHandler)
}

// UserPromptSubmitted registers a UserPromptSubmitted handler on the registrar.
func (c *hooks) UserPromptSubmitted(fn func(context.Context, UserPromptSubmitted) error) *hooks {
	return bind(c, fn, userpromptsubmitted.RegisterHandler)
}

// PreToolUse registers a PreToolUse handler on the registrar.
func (c *hooks) PreToolUse(fn func(context.Context, PreToolUse, PreToolResults) (PreToolOutput, error)) *hooks {
	return bind(c, fn, pretooluse.RegisterHandler)
}

// PostToolUse registers a PostToolUse handler on the registrar.
func (c *hooks) PostToolUse(fn func(context.Context, PostToolUse, PostToolResults) (PostToolOutput, error)) *hooks {
	return bind(c, fn, posttooluse.RegisterHandler)
}

// PostToolUseFailure registers a PostToolUseFailure handler on the registrar.
func (c *hooks) PostToolUseFailure(fn func(context.Context, PostToolUseFailure, PostToolFailureResults) (PostToolFailureOutput, error)) *hooks {
	return bind(c, fn, posttoolusefailure.RegisterHandler)
}

// PermissionRequest registers a PermissionRequest handler on the registrar.
func (c *hooks) PermissionRequest(fn func(context.Context, PermissionRequest, PermissionRequestResults) (PermissionRequestOutput, error)) *hooks {
	return bind(c, fn, permissionrequest.RegisterHandler)
}

// SubagentStart registers a SubagentStart handler on the registrar.
func (c *hooks) SubagentStart(fn func(context.Context, SubagentStart, SubagentStartResults) (SubagentStartOutput, error)) *hooks {
	return bind(c, fn, subagentstart.RegisterHandler)
}

// SubagentStop registers a SubagentStop handler on the registrar.
func (c *hooks) SubagentStop(fn func(context.Context, SubagentStop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, subagentstop.RegisterHandler)
}

// AgentStop registers an AgentStop handler on the registrar.
func (c *hooks) AgentStop(fn func(context.Context, AgentStop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, agentstop.RegisterHandler)
}

// PreCompact registers a PreCompact handler on the registrar.
func (c *hooks) PreCompact(fn func(context.Context, PreCompact) error) *hooks {
	return bind(c, fn, precompact.RegisterHandler)
}

// Notification registers a Notification handler on the registrar.
func (c *hooks) Notification(fn func(context.Context, Notification, NotificationResults) (NotificationOutput, error)) *hooks {
	return bind(c, fn, notification.RegisterHandler)
}

// ErrorOccurred registers an ErrorOccurred handler on the registrar.
func (c *hooks) ErrorOccurred(fn func(context.Context, ErrorOccurred) error) *hooks {
	return bind(c, fn, erroroccurred.RegisterHandler)
}
