package agenthooks

import (
	"encoding/json"
	"fmt"
)

// CursorWarnExit is exit code 2. Cursor treats it as block/deny, equivalent to
// returning permission:"deny" on permission-gating events.
const CursorWarnExit = 2

// CursorHandlerErrorExit is exit code 1. The runner should use this when a
// handler returns an error under Cursor's default fail-open policy.
const CursorHandlerErrorExit = 1

// CursorCodec implements Codec for Cursor hooks.
// Reference: https://cursor.com/docs/hooks
//
// Cursor's dedicated surfaces are folded into the unified tool kinds:
// beforeShellExecution/beforeMCPExecution/beforeReadFile → KindPreTool,
// afterShellExecution/afterMCPExecution/afterFileEdit → KindPostTool.
// Event.Name preserves the native surface; Event.Raw holds the full payload.
type CursorCodec struct{}

// Dialect returns Cursor.
func (c *CursorCodec) Dialect() Dialect { return Cursor }

var cursorKinds = map[string]Kind{
	"sessionStart":         KindSessionStart,
	"sessionEnd":           KindSessionEnd,
	"beforeSubmitPrompt":   KindUserPrompt,
	"preToolUse":           KindPreTool,
	"postToolUse":          KindPostTool,
	"postToolUseFailure":   KindPostToolFailure,
	"beforeShellExecution": KindPreTool,
	"afterShellExecution":  KindPostTool,
	"beforeMCPExecution":   KindPreTool,
	"afterMCPExecution":    KindPostTool,
	"beforeReadFile":       KindPreTool,
	"afterFileEdit":        KindPostTool,
	"subagentStart":        KindSubagentStart,
	"subagentStop":         KindSubagentStop,
	"stop":                 KindStop,
	"preCompact":           KindPreCompact,
	"afterAgentResponse":   KindOther,
	"afterAgentThought":    KindOther,
	"beforeTabFileRead":    KindOther,
	"afterTabFileEdit":     KindOther,
	"workspaceOpen":        KindOther,
}

type cursorAttachment struct {
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}

type cursorEdit struct {
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}

// cursorPayload is a superset decode of every Cursor hook event we normalize.
type cursorPayload struct {
	ConversationID      string             `json:"conversation_id"`
	GenerationID        string             `json:"generation_id"`
	Model               string             `json:"model"`
	HookEventName       string             `json:"hook_event_name"`
	CursorVersion       string             `json:"cursor_version"`
	WorkspaceRoots      []string           `json:"workspace_roots"`
	UserEmail           *string            `json:"user_email"`
	TranscriptPath      *string            `json:"transcript_path"`
	Cwd                 string             `json:"cwd"`
	ToolName            string             `json:"tool_name"`
	ToolInput           json.RawMessage    `json:"tool_input"`
	ToolUseID           string             `json:"tool_use_id"`
	ToolOutput          string             `json:"tool_output"`
	Command             string             `json:"command"`
	Output              string             `json:"output"`
	ResultJSON          string             `json:"result_json"`
	FilePath            string             `json:"file_path"`
	Content             string             `json:"content"`
	Attachments         []cursorAttachment `json:"attachments"`
	Edits               []cursorEdit       `json:"edits"`
	Prompt              string             `json:"prompt"`
	Text                string             `json:"text"`
	ErrorMessage        string             `json:"error_message"`
	FailureType         string             `json:"failure_type"`
	Duration            int64              `json:"duration"`
	DurationMs          int64              `json:"duration_ms"`
	SubagentID          string             `json:"subagent_id"`
	SubagentType        string             `json:"subagent_type"`
	Task                string             `json:"task"`
	Summary             string             `json:"summary"`
	Status              string             `json:"status"`
	LoopCount           int                `json:"loop_count"`
	AgentTranscriptPath *string            `json:"agent_transcript_path"`
	SessionID           string             `json:"session_id"`
	Reason              string             `json:"reason"`
	IsBackgroundAgent   bool               `json:"is_background_agent"`
	Trigger             string             `json:"trigger"`
}

func (p *cursorPayload) durationMs() int64 {
	if p.DurationMs != 0 {
		return p.DurationMs
	}
	return p.Duration
}

func cursorPermissionEvent(ev *Event) bool {
	if ev == nil {
		return false
	}
	switch ev.Kind {
	case KindPreTool, KindSubagentStart:
		return true
	case KindOther:
		return ev.Name == "beforeTabFileRead"
	default:
		return false
	}
}

// Decode parses a Cursor hook stdin payload into a unified Event.
func (c *CursorCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	var p cursorPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("cursor: decode payload: %w", err)
	}
	name := p.HookEventName
	if name == "" {
		name = eventHint
	}
	kind, ok := cursorKinds[name]
	if !ok {
		kind = KindOther
	}
	session := p.ConversationID
	if session == "" {
		session = p.SessionID
	}
	transcript := ""
	if p.TranscriptPath != nil {
		transcript = *p.TranscriptPath
	}
	ev := &Event{
		Agent:          Cursor,
		Kind:           kind,
		Name:           name,
		Session:        session,
		Cwd:            p.Cwd,
		TranscriptPath: transcript,
		Raw:            cloneRaw(raw),
	}

	switch name {
	case "sessionStart":
		ev.Life = &Lifecycle{Model: p.Model, Background: p.IsBackgroundAgent}
	case "sessionEnd":
		ev.Life = &Lifecycle{Reason: p.Reason, Background: p.IsBackgroundAgent}
	case "beforeSubmitPrompt":
		ev.Prompt = p.Prompt
	case "preToolUse":
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
	case "postToolUse":
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
		ev.Result = &ToolResult{Text: p.ToolOutput, DurationMs: p.durationMs()}
	case "postToolUseFailure":
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
		ev.Result = &ToolResult{
			Error:       p.ErrorMessage,
			FailureType: p.FailureType,
			DurationMs:  p.durationMs(),
		}
	case "beforeShellExecution":
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: p.Command}
	case "afterShellExecution":
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: p.Command}
		ev.Result = &ToolResult{Text: p.Output, DurationMs: p.durationMs()}
	case "beforeMCPExecution":
		nameNorm, _ := NormalizeToolName(p.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(p.ToolInput),
			MCP:    true,
		}
	case "afterMCPExecution":
		nameNorm, _ := NormalizeToolName(p.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(p.ToolInput),
			MCP:    true,
		}
		ev.Result = &ToolResult{Text: p.ResultJSON, DurationMs: p.durationMs()}
	case "beforeReadFile":
		input, err := json.Marshal(map[string]any{
			"file_path":   p.FilePath,
			"content":     p.Content,
			"attachments": p.Attachments,
		})
		if err != nil {
			return nil, fmt.Errorf("cursor: decode beforeReadFile input: %w", err)
		}
		ev.Tool = &ToolCall{Name: ToolRead, Native: name, Input: input}
	case "afterFileEdit":
		input, err := json.Marshal(map[string]any{
			"file_path": p.FilePath,
			"edits":     p.Edits,
		})
		if err != nil {
			return nil, fmt.Errorf("cursor: decode afterFileEdit input: %w", err)
		}
		editsRaw, err := json.Marshal(p.Edits)
		if err != nil {
			return nil, fmt.Errorf("cursor: decode afterFileEdit edits: %w", err)
		}
		ev.Tool = &ToolCall{Name: ToolEdit, Native: name, Input: input}
		ev.Result = &ToolResult{Raw: editsRaw}
	case "subagentStart":
		ev.Subagent = &Subagent{ID: p.SubagentID, Type: p.SubagentType, Task: p.Task}
	case "subagentStop":
		tp := ""
		if p.AgentTranscriptPath != nil {
			tp = *p.AgentTranscriptPath
		}
		ev.Subagent = &Subagent{
			ID:             p.SubagentID,
			Type:           p.SubagentType,
			Task:           p.Task,
			Summary:        p.Summary,
			Status:         p.Status,
			TranscriptPath: tp,
			LoopCount:      p.LoopCount,
		}
	case "stop":
		ev.Turn = &TurnEnd{Status: p.Status, LoopCount: p.LoopCount}
	case "preCompact":
		ev.Compact = &CompactInfo{Trigger: p.Trigger}
	case "afterAgentResponse", "afterAgentThought":
		ev.Note = &Note{Message: p.Text}
		if name == "afterAgentThought" {
			ev.Note.Type = "thought"
		}
	}
	return ev, nil
}

// Encode renders a unified Result as Cursor stdout JSON and exit code.
// ev must be non-nil.
func (c *CursorCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("cursor: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}

	out := map[string]any{}
	exitCode := 0

	switch {
	case cursorPermissionEvent(ev):
		if d := res.Decision.String(); d != "" {
			out["permission"] = d
		}
		if res.UserMessage != "" {
			out["user_message"] = res.UserMessage
		}
		if res.Reason != "" {
			out["agent_message"] = res.Reason
		}
		if res.UpdatedInput != nil && ev.Name == "preToolUse" {
			out["updated_input"] = res.UpdatedInput
		}
		if res.Decision == DecisionDeny {
			exitCode = CursorWarnExit
		}
	case ev.Kind == KindUserPrompt:
		if res.BlockPrompt || res.Decision == DecisionDeny {
			out["continue"] = false
			msg := res.UserMessage
			if msg == "" {
				msg = res.Reason
			}
			if msg != "" {
				out["user_message"] = msg
			}
		}
	case ev.Kind == KindPostTool:
		if res.UpdatedOutput != nil {
			out["updated_mcp_tool_output"] = *res.UpdatedOutput
		}
		if res.Context != "" {
			out["additional_context"] = res.Context
		}
	case ev.Kind == KindStop, ev.Kind == KindSubagentStop:
		if res.FollowUp != "" {
			out["followup_message"] = res.FollowUp
		}
	case ev.Kind == KindSessionStart:
		if len(res.Env) > 0 {
			out["env"] = res.Env
		}
		if res.Context != "" {
			out["additional_context"] = res.Context
		}
	case ev.Kind == KindPreCompact:
		if res.UserMessage != "" {
			out["user_message"] = res.UserMessage
		}
	}

	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, exitCode, err
}
