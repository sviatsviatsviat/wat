package cursor

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentresponse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentthought"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/compact/precompact"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/prompt/beforesubmitprompt"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/workspaceopen"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/afterfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftermcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftershellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftertabfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforemcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforereadfile"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforeshellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforetabfileread"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/pretooluse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// hooks is a typed Cursor registrar; callers obtain one via UseHooks.
type hooks struct {
	hookkit.HandlerQueue
}

var codec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

// UseHooks returns a fluent registrar for this package's Cursor dialect.
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

// Contribute installs the Cursor dialect and registers these hooks' handlers.
func (c *hooks) Contribute(reg run.Registry) {
	if c == nil {
		return
	}
	c.Install(reg, runtime.Dialect, detectPayload, codec)
}

func detectPayload(raw []byte) bool {
	return detectPayloadWith(raw, os.Getenv)
}

func detectPayloadWith(raw []byte, getenv func(string) string) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return true
	}
	if getenv != nil && getenv("CURSOR_VERSION") != "" {
		return true
	}
	return false
}

// SessionStart registers a SessionStart handler on the registrar.
func (c *hooks) SessionStart(fn func(context.Context, SessionStart, SessionStartResults) (SessionStartOutput, error)) *hooks {
	return bind(c, fn, sessionstart.RegisterHandler)
}

// SessionEnd registers a SessionEnd handler on the registrar.
func (c *hooks) SessionEnd(fn func(context.Context, SessionEnd) error) *hooks {
	return bind(c, fn, sessionend.RegisterHandler)
}

// BeforeSubmitPrompt registers a BeforeSubmitPrompt handler on the registrar.
func (c *hooks) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPrompt, BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *hooks {
	return bind(c, fn, beforesubmitprompt.RegisterHandler)
}

// PreToolUse registers a PreToolUse handler on the registrar.
func (c *hooks) PreToolUse(fn func(context.Context, PreToolUse, PermissionResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, pretooluse.RegisterHandler)
}

// PostToolUse registers a PostToolUse handler on the registrar.
func (c *hooks) PostToolUse(fn func(context.Context, PostToolUse, PostToolResults) (PostToolOutput, error)) *hooks {
	return bind(c, fn, posttooluse.RegisterHandler)
}

// PostToolUseFailure registers a PostToolUseFailure handler on the registrar.
func (c *hooks) PostToolUseFailure(fn func(context.Context, PostToolUseFailure, PostToolResults) (PostToolOutput, error)) *hooks {
	return bind(c, fn, posttoolusefailure.RegisterHandler)
}

// BeforeShellExecution registers a BeforeShellExecution handler on the registrar.
func (c *hooks) BeforeShellExecution(fn func(context.Context, BeforeShellExecution, PermissionResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, beforeshellexecution.RegisterHandler)
}

// AfterShellExecution registers an AfterShellExecution handler on the registrar.
func (c *hooks) AfterShellExecution(fn func(context.Context, AfterShellExecution, PostToolResults) (PostToolOutput, error)) *hooks {
	return bind(c, fn, aftershellexecution.RegisterHandler)
}

// BeforeMCPExecution registers a BeforeMCPExecution handler on the registrar.
func (c *hooks) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecution, PermissionResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, beforemcpexecution.RegisterHandler)
}

// AfterMCPExecution registers an observe-only AfterMCPExecution handler.
// Cursor documents no output fields for this event. Rewrite MCP tool output
// with PostToolUse (updated_mcp_tool_output). Cloud agents do not run MCP hooks.
func (c *hooks) AfterMCPExecution(fn func(context.Context, AfterMCPExecution) error) *hooks {
	return bind(c, fn, aftermcpexecution.RegisterHandler)
}

// BeforeReadFile registers a BeforeReadFile handler on the registrar.
func (c *hooks) BeforeReadFile(fn func(context.Context, BeforeReadFile, BeforeReadFileResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, beforereadfile.RegisterHandler)
}

// AfterFileEdit registers an AfterFileEdit handler on the registrar.
func (c *hooks) AfterFileEdit(fn func(context.Context, AfterFileEdit, PostToolResults) (PostToolOutput, error)) *hooks {
	return bind(c, fn, afterfileedit.RegisterHandler)
}

// SubagentStart registers a SubagentStart handler on the registrar.
func (c *hooks) SubagentStart(fn func(context.Context, SubagentStart, SubagentStartResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, subagentstart.RegisterHandler)
}

// SubagentStop registers a SubagentStop handler on the registrar.
func (c *hooks) SubagentStop(fn func(context.Context, SubagentStop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, subagentstop.RegisterHandler)
}

// Stop registers a Stop handler on the registrar.
func (c *hooks) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *hooks {
	return bind(c, fn, stopevent.RegisterHandler)
}

// PreCompact registers a PreCompact handler on the registrar.
func (c *hooks) PreCompact(fn func(context.Context, PreCompact, PreCompactResults) (PreCompactOutput, error)) *hooks {
	return bind(c, fn, precompact.RegisterHandler)
}

// AfterAgentResponse registers an AfterAgentResponse handler on the registrar.
func (c *hooks) AfterAgentResponse(fn func(context.Context, AfterAgentResponse) error) *hooks {
	return bind(c, fn, afteragentresponse.RegisterHandler)
}

// AfterAgentThought registers an AfterAgentThought handler on the registrar.
func (c *hooks) AfterAgentThought(fn func(context.Context, AfterAgentThought) error) *hooks {
	return bind(c, fn, afteragentthought.RegisterHandler)
}

// BeforeTabFileRead registers a BeforeTabFileRead handler on the registrar.
func (c *hooks) BeforeTabFileRead(fn func(context.Context, BeforeTabFileRead, BeforeTabFileReadResults) (PermissionOutput, error)) *hooks {
	return bind(c, fn, beforetabfileread.RegisterHandler)
}

// AfterTabFileEdit registers an AfterTabFileEdit handler on the registrar.
func (c *hooks) AfterTabFileEdit(fn func(context.Context, AfterTabFileEdit) error) *hooks {
	return bind(c, fn, aftertabfileedit.RegisterHandler)
}

// WorkspaceOpen registers a WorkspaceOpen handler on the registrar.
func (c *hooks) WorkspaceOpen(fn func(context.Context, WorkspaceOpen) error) *hooks {
	return bind(c, fn, workspaceopen.RegisterHandler)
}
