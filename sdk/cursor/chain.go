package cursor

import (
	"context"

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

// BeforeSubmitPrompt registers a BeforeSubmitPrompt handler on the chain.
func (c *chain) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPrompt, BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *chain {
	beforesubmitprompt.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, PreToolUse, PermissionResults) (PermissionOutput, error)) *chain {
	pretooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, PostToolUse, PostToolResults) (PostToolOutput, error)) *chain {
	posttooluse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailure, PostToolResults) (PostToolOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// BeforeShellExecution registers a BeforeShellExecution handler on the chain.
func (c *chain) BeforeShellExecution(fn func(context.Context, BeforeShellExecution, PermissionResults) (PermissionOutput, error)) *chain {
	beforeshellexecution.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterShellExecution registers an AfterShellExecution handler on the chain.
func (c *chain) AfterShellExecution(fn func(context.Context, AfterShellExecution, PostToolResults) (PostToolOutput, error)) *chain {
	aftershellexecution.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// BeforeMCPExecution registers a BeforeMCPExecution handler on the chain.
func (c *chain) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecution, PermissionResults) (PermissionOutput, error)) *chain {
	beforemcpexecution.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterMCPExecution registers an AfterMCPExecution handler on the chain.
func (c *chain) AfterMCPExecution(fn func(context.Context, AfterMCPExecution, PostToolResults) (PostToolOutput, error)) *chain {
	aftermcpexecution.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// BeforeReadFile registers a BeforeReadFile handler on the chain.
func (c *chain) BeforeReadFile(fn func(context.Context, BeforeReadFile, BeforeReadFileResults) (PermissionOutput, error)) *chain {
	beforereadfile.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterFileEdit registers an AfterFileEdit handler on the chain.
func (c *chain) AfterFileEdit(fn func(context.Context, AfterFileEdit, PostToolResults) (PostToolOutput, error)) *chain {
	afterfileedit.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, SubagentStart, SubagentStartResults) (PermissionOutput, error)) *chain {
	subagentstart.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, SubagentStop, StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// Stop registers a Stop handler on the chain.
func (c *chain) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *chain {
	stopevent.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, PreCompact, PreCompactResults) (PreCompactOutput, error)) *chain {
	precompact.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterAgentResponse registers an AfterAgentResponse handler on the chain.
func (c *chain) AfterAgentResponse(fn func(context.Context, AfterAgentResponse) error) *chain {
	afteragentresponse.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterAgentThought registers an AfterAgentThought handler on the chain.
func (c *chain) AfterAgentThought(fn func(context.Context, AfterAgentThought) error) *chain {
	afteragentthought.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// BeforeTabFileRead registers a BeforeTabFileRead handler on the chain.
func (c *chain) BeforeTabFileRead(fn func(context.Context, BeforeTabFileRead, BeforeTabFileReadResults) (PermissionOutput, error)) *chain {
	beforetabfileread.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// AfterTabFileEdit registers an AfterTabFileEdit handler on the chain.
func (c *chain) AfterTabFileEdit(fn func(context.Context, AfterTabFileEdit) error) *chain {
	aftertabfileedit.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}

// WorkspaceOpen registers a WorkspaceOpen handler on the chain.
func (c *chain) WorkspaceOpen(fn func(context.Context, WorkspaceOpen) error) *chain {
	workspaceopen.RegisterHandler(runtime.DefaultDialect, fn)
	return c
}
