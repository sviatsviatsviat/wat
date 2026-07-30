package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcompleted"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcreated"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/teammateidle"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact/postcompact"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact/precompact"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitation"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitationresult"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptexpansion"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptsubmit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/setup"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopfailure"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissiondenied"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttoolbatch"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/pretooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/messagedisplay"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/notification"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/configchange"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/cwdchanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/filechanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/instructionsloaded"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreecreate"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreeremove"
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

// Setup registers a Setup handler on the registrar.
func (c *hooks) Setup(fn func(context.Context, Setup, SetupResults) (SetupOutput, error)) *hooks {
	return bind(c, fn, setup.RegisterHandler)
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

// PostToolBatch registers a PostToolBatch handler on the registrar.
func (c *hooks) PostToolBatch(fn func(context.Context, PostToolBatch, PostToolBatchResults) (PostToolBatchOutput, error)) *hooks {
	return bind(c, fn, posttoolbatch.RegisterHandler)
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

// TeammateIdle registers a TeammateIdle handler on the registrar.
func (c *hooks) TeammateIdle(fn func(context.Context, TeammateIdle, TeammateIdleResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, teammateidle.RegisterHandler)
}

// Stop registers a Stop handler on the registrar.
func (c *hooks) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, stopevent.RegisterHandler)
}

// StopFailure registers a StopFailure handler on the registrar.
func (c *hooks) StopFailure(fn func(context.Context, StopFailure) error) *hooks {
	return bind(c, fn, stopfailure.RegisterHandler)
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

// WorktreeRemove registers a WorktreeRemove handler on the registrar.
func (c *hooks) WorktreeRemove(fn func(context.Context, WorktreeRemove) error) *hooks {
	return bind(c, fn, worktreeremove.RegisterHandler)
}

// PreCompact registers a PreCompact handler on the registrar.
func (c *hooks) PreCompact(fn func(context.Context, PreCompact, PreCompactResults) (CommonOutput, error)) *hooks {
	return bind(c, fn, precompact.RegisterHandler)
}

// PostCompact registers a PostCompact handler on the registrar.
func (c *hooks) PostCompact(fn func(context.Context, PostCompact) error) *hooks {
	return bind(c, fn, postcompact.RegisterHandler)
}

// Elicitation registers an Elicitation handler on the registrar.
func (c *hooks) Elicitation(fn func(context.Context, Elicitation, ElicitationResults) (ElicitationOutput, error)) *hooks {
	return bind(c, fn, elicitation.RegisterHandler)
}

// ElicitationResult registers an ElicitationResult handler on the registrar.
func (c *hooks) ElicitationResult(fn func(context.Context, ElicitationResult, ElicitationResultResults) (ElicitationResultOutput, error)) *hooks {
	return bind(c, fn, elicitationresult.RegisterHandler)
}

// InstructionsLoaded registers an InstructionsLoaded handler on the registrar.
func (c *hooks) InstructionsLoaded(fn func(context.Context, InstructionsLoaded) error) *hooks {
	return bind(c, fn, instructionsloaded.RegisterHandler)
}

// ConfigChange registers a ConfigChange handler on the registrar.
func (c *hooks) ConfigChange(fn func(context.Context, ConfigChange, ConfigChangeResults) (ConfigChangeOutput, error)) *hooks {
	return bind(c, fn, configchange.RegisterHandler)
}

// CwdChanged registers a CwdChanged handler on the registrar.
func (c *hooks) CwdChanged(fn func(context.Context, CwdChanged, CwdChangedResults) (CwdChangedOutput, error)) *hooks {
	return bind(c, fn, cwdchanged.RegisterHandler)
}

// FileChanged registers a FileChanged handler on the registrar.
func (c *hooks) FileChanged(fn func(context.Context, FileChanged, FileChangedResults) (FileChangedOutput, error)) *hooks {
	return bind(c, fn, filechanged.RegisterHandler)
}
