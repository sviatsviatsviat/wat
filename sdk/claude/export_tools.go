package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/tools"
)

// Input is the tool input payload on a Claude Code hook event.
type Input = tools.Input

// BashInput is the typed input for the Bash tool.
type BashInput = tools.BashInput

// WriteInput is the typed input for the Write tool.
type WriteInput = tools.WriteInput

// EditInput is the typed input for the Edit tool.
type EditInput = tools.EditInput

// ReadInput is the typed input for the Read tool.
type ReadInput = tools.ReadInput

// GlobInput is the typed input for the Glob tool.
type GlobInput = tools.GlobInput

// GrepInput is the typed input for the Grep tool.
type GrepInput = tools.GrepInput

// WebFetchInput is the typed input for the WebFetch tool.
type WebFetchInput = tools.WebFetchInput

// WebSearchInput is the typed input for the WebSearch tool.
type WebSearchInput = tools.WebSearchInput

// AskUserQuestionInput is the typed input for the AskUserQuestion tool.
type AskUserQuestionInput = tools.AskUserQuestionInput

// Question is one AskUserQuestion prompt entry.
type Question = tools.Question

// Option is one AskUserQuestion choice.
type Option = tools.Option

// ExitPlanModeInput is the typed input for the ExitPlanMode tool.
type ExitPlanModeInput = tools.ExitPlanModeInput

// AgentInput is the typed input for the Agent tool.
type AgentInput = tools.AgentInput

// MCPTool is a parsed mcp__server__tool name and input payload.
type MCPTool = tools.MCPTool

const (
	// ToolBash is the Bash tool name.
	ToolBash = tools.ToolBash
	// ToolWrite is the Write tool name.
	ToolWrite = tools.ToolWrite
	// ToolEdit is the Edit tool name.
	ToolEdit = tools.ToolEdit
	// ToolRead is the Read tool name.
	ToolRead = tools.ToolRead
	// ToolGlob is the Glob tool name.
	ToolGlob = tools.ToolGlob
	// ToolGrep is the Grep tool name.
	ToolGrep = tools.ToolGrep
	// ToolWebFetch is the WebFetch tool name.
	ToolWebFetch = tools.ToolWebFetch
	// ToolWebSearch is the WebSearch tool name.
	ToolWebSearch = tools.ToolWebSearch
	// ToolAskUserQuestion is the AskUserQuestion tool name.
	ToolAskUserQuestion = tools.ToolAskUserQuestion
	// ToolExitPlanMode is the ExitPlanMode tool name.
	ToolExitPlanMode = tools.ToolExitPlanMode
	// ToolAgent is the Agent tool name.
	ToolAgent = tools.ToolAgent
)
