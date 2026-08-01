package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

// Input is the tool input payload on a Copilot hook event.
type Input = tools.Input

// BashInput is the input schema for the bash tool.
type BashInput = tools.BashInput

// CreateInput is the typed input for the create tool.
type CreateInput = tools.CreateInput

// EditInput is the typed input for the edit tool.
type EditInput = tools.EditInput

// ViewInput is the typed input for the view tool.
type ViewInput = tools.ViewInput

// WebFetchInput is the typed input for the webfetch tool.
type WebFetchInput = tools.WebFetchInput

// GlobInput is the typed input for the glob tool.
type GlobInput = tools.GlobInput

// GrepInput is the typed input for the grep tool.
type GrepInput = tools.GrepInput

// WebSearchInput is the typed input for the web_search tool.
type WebSearchInput = tools.WebSearchInput

// TaskInput is the typed input for the task/agent tool.
type TaskInput = tools.TaskInput

// AskUserInput is the typed input for the ask_user tool.
type AskUserInput = tools.AskUserInput

// AskUserQuestion is one question in AskUserInput.
type AskUserQuestion = tools.AskUserQuestion

// AskUserOption is one option in AskUserQuestion.
type AskUserOption = tools.AskUserOption

// UpdateTodoInput is the typed input for update_todo / TodoWrite.
type UpdateTodoInput = tools.UpdateTodoInput

// TodoItem is one entry in UpdateTodoInput.
type TodoItem = tools.TodoItem

const (
	// ToolBash is the bash tool name.
	ToolBash = tools.ToolBash
	// ToolPowerShell is the PowerShell tool name.
	ToolPowerShell = tools.ToolPowerShell
	// ToolShell is the shell tool name.
	ToolShell = tools.ToolShell
	// ToolCreate is the create tool name.
	ToolCreate = tools.ToolCreate
	// ToolWrite is the Claude-format Write tool name accepted by AsCreate.
	ToolWrite = tools.ToolWrite
	// ToolEdit is the edit tool name.
	ToolEdit = tools.ToolEdit
	// ToolView is the view tool name.
	ToolView = tools.ToolView
	// ToolRead is the Claude-format Read tool name accepted by AsView.
	ToolRead = tools.ToolRead
	// ToolWebFetch is the webfetch tool name.
	ToolWebFetch = tools.ToolWebFetch
	// ToolGlob is the glob tool name.
	ToolGlob = tools.ToolGlob
	// ToolGrep is the grep tool name.
	ToolGrep = tools.ToolGrep
	// ToolWebSearch is the web_search tool name.
	ToolWebSearch = tools.ToolWebSearch
	// ToolTask is the task tool name.
	ToolTask = tools.ToolTask
	// ToolAskUser is the ask_user tool name.
	ToolAskUser = tools.ToolAskUser
	// ToolUpdateTodo is the update_todo tool name.
	ToolUpdateTodo = tools.ToolUpdateTodo
)
