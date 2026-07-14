package tools

import (
	"encoding/json"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// BashInput is the input schema for the Bash tool.
type BashInput struct {
	Command         string `json:"command"`
	Description     string `json:"description,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	RunInBackground bool   `json:"run_in_background,omitempty"`
}

// WriteInput is the input schema for the Write tool.
type WriteInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// EditInput is the input schema for the Edit tool.
type EditInput struct {
	FilePath   string `json:"file_path"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all,omitempty"`
}

// ReadInput is the input schema for the Read tool.
type ReadInput struct {
	FilePath string `json:"file_path"`
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// GlobInput is the input schema for the Glob tool.
type GlobInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// GrepInput is the input schema for the Grep tool.
type GrepInput struct {
	Pattern         string `json:"pattern"`
	Path            string `json:"path,omitempty"`
	Glob            string `json:"glob,omitempty"`
	OutputMode      string `json:"output_mode,omitempty"`
	CaseInsensitive bool   `json:"-i,omitempty"`
	Multiline       bool   `json:"multiline,omitempty"`
}

// WebFetchInput is the input schema for the WebFetch tool.
type WebFetchInput struct {
	URL    string `json:"url"`
	Prompt string `json:"prompt,omitempty"`
}

// WebSearchInput is the input schema for the WebSearch tool.
type WebSearchInput struct {
	Query          string   `json:"query"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	BlockedDomains []string `json:"blocked_domains,omitempty"`
}

// AgentInput is the input schema for the Agent tool.
type AgentInput struct {
	Prompt       string `json:"prompt"`
	Description  string `json:"description,omitempty"`
	SubagentType string `json:"subagent_type,omitempty"`
	Model        string `json:"model,omitempty"`
}

// AskUserQuestionInput is the input schema for the AskUserQuestion tool.
type AskUserQuestionInput struct {
	Questions []Question        `json:"questions"`
	Answers   map[string]string `json:"answers,omitempty"`
}

// Question is one question in AskUserQuestionInput.
type Question struct {
	Question    string   `json:"question"`
	Header      string   `json:"header,omitempty"`
	Options     []Option `json:"options,omitempty"`
	MultiSelect bool     `json:"multiSelect,omitempty"`
}

// Option is one selectable option in a Question.
type Option struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// ExitPlanModeInput is the input schema for the ExitPlanMode tool.
type ExitPlanModeInput struct {
	Plan         string `json:"plan"`
	PlanFilePath string `json:"plan_file_path,omitempty"`
}

// ToolInputAs decodes raw tool input JSON into T.
func ToolInputAs[T any](raw json.RawMessage) (T, error) {
	return hookkit.ToolInputAs[T](raw)
}

// IsMCPTool reports whether name is an MCP tool and returns server and tool parts.
func IsMCPTool(name string) (server, tool string, ok bool) {
	const prefix = "mcp__"
	if !strings.HasPrefix(name, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(name, prefix)
	parts := strings.SplitN(rest, "__", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}
