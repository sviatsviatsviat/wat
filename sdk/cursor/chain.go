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
		panic("cursor: UseHooks: at most one registry")
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

// BeforeSubmitPrompt registers a BeforeSubmitPrompt handler on the chain.
func (c *chain) BeforeSubmitPrompt(fn func(context.Context, run.Hook[BeforeSubmitPrompt], BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *chain {
	beforesubmitprompt.RegisterHandler(c.reg, fn)
	return c
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PermissionResults) (PermissionOutput, error)) *chain {
	pretooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	posttooluse.RegisterHandler(c.reg, fn)
	return c
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolResults) (PostToolOutput, error)) *chain {
	posttoolusefailure.RegisterHandler(c.reg, fn)
	return c
}

// BeforeShellExecution registers a BeforeShellExecution handler on the chain.
func (c *chain) BeforeShellExecution(fn func(context.Context, run.Hook[BeforeShellExecution], PermissionResults) (PermissionOutput, error)) *chain {
	beforeshellexecution.RegisterHandler(c.reg, fn)
	return c
}

// AfterShellExecution registers an AfterShellExecution handler on the chain.
func (c *chain) AfterShellExecution(fn func(context.Context, run.Hook[AfterShellExecution], PostToolResults) (PostToolOutput, error)) *chain {
	aftershellexecution.RegisterHandler(c.reg, fn)
	return c
}

// BeforeMCPExecution registers a BeforeMCPExecution handler on the chain.
func (c *chain) BeforeMCPExecution(fn func(context.Context, run.Hook[BeforeMCPExecution], PermissionResults) (PermissionOutput, error)) *chain {
	beforemcpexecution.RegisterHandler(c.reg, fn)
	return c
}

// AfterMCPExecution registers an AfterMCPExecution handler on the chain.
func (c *chain) AfterMCPExecution(fn func(context.Context, run.Hook[AfterMCPExecution], PostToolResults) (PostToolOutput, error)) *chain {
	aftermcpexecution.RegisterHandler(c.reg, fn)
	return c
}

// BeforeReadFile registers a BeforeReadFile handler on the chain.
func (c *chain) BeforeReadFile(fn func(context.Context, run.Hook[BeforeReadFile], BeforeReadFileResults) (PermissionOutput, error)) *chain {
	beforereadfile.RegisterHandler(c.reg, fn)
	return c
}

// AfterFileEdit registers an AfterFileEdit handler on the chain.
func (c *chain) AfterFileEdit(fn func(context.Context, run.Hook[AfterFileEdit], PostToolResults) (PostToolOutput, error)) *chain {
	afterfileedit.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (PermissionOutput, error)) *chain {
	subagentstart.RegisterHandler(c.reg, fn)
	return c
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	subagentstop.RegisterHandler(c.reg, fn)
	return c
}

// Stop registers a Stop handler on the chain.
func (c *chain) Stop(fn func(context.Context, run.Hook[Stop], StopResults) (StopOutput, error)) *chain {
	stopevent.RegisterHandler(c.reg, fn)
	return c
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact], PreCompactResults) (PreCompactOutput, error)) *chain {
	precompact.RegisterHandler(c.reg, fn)
	return c
}

// AfterAgentResponse registers an AfterAgentResponse handler on the chain.
func (c *chain) AfterAgentResponse(fn func(context.Context, run.Hook[AfterAgentResponse]) error) *chain {
	afteragentresponse.RegisterHandler(c.reg, fn)
	return c
}

// AfterAgentThought registers an AfterAgentThought handler on the chain.
func (c *chain) AfterAgentThought(fn func(context.Context, run.Hook[AfterAgentThought]) error) *chain {
	afteragentthought.RegisterHandler(c.reg, fn)
	return c
}

// BeforeTabFileRead registers a BeforeTabFileRead handler on the chain.
func (c *chain) BeforeTabFileRead(fn func(context.Context, run.Hook[BeforeTabFileRead], BeforeTabFileReadResults) (PermissionOutput, error)) *chain {
	beforetabfileread.RegisterHandler(c.reg, fn)
	return c
}

// AfterTabFileEdit registers an AfterTabFileEdit handler on the chain.
func (c *chain) AfterTabFileEdit(fn func(context.Context, run.Hook[AfterTabFileEdit]) error) *chain {
	aftertabfileedit.RegisterHandler(c.reg, fn)
	return c
}

// WorkspaceOpen registers a WorkspaceOpen handler on the chain.
func (c *chain) WorkspaceOpen(fn func(context.Context, run.Hook[WorkspaceOpen]) error) *chain {
	workspaceopen.RegisterHandler(c.reg, fn)
	return c
}
