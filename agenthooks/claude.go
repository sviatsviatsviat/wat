package agenthooks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// ClaudeCodec implements Codec for Claude Code hooks.
// Reference: https://code.claude.com/docs/en/hooks
type ClaudeCodec struct {
	// Getenv and AppendFile are injectable for tests. They back the
	// CLAUDE_ENV_FILE side effect used to express Result.Env.
	Getenv     func(string) string
	AppendFile func(path string, data []byte) error
}

// Dialect returns Claude.
func (c *ClaudeCodec) Dialect() Dialect { return Claude }

// claudePayload is a superset decode of every Claude Code event we normalize.
type claudePayload struct {
	SessionID            string          `json:"session_id"`
	TranscriptPath       string          `json:"transcript_path"`
	Cwd                  string          `json:"cwd"`
	HookEventName        string          `json:"hook_event_name"`
	AgentID              string          `json:"agent_id"`
	AgentType            string          `json:"agent_type"`
	Source               string          `json:"source"`
	Reason               string          `json:"reason"`
	Model                string          `json:"model"`
	Prompt               string          `json:"prompt"`
	ToolName             string          `json:"tool_name"`
	ToolInput            json.RawMessage `json:"tool_input"`
	ToolUseID            string          `json:"tool_use_id"`
	ToolResponse         json.RawMessage `json:"tool_response"`
	Error                string          `json:"error"`
	StopHookActive       bool            `json:"stop_hook_active"`
	LastAssistantMessage string          `json:"last_assistant_message"`
	Trigger              string          `json:"trigger"`
	CustomInstructions   string          `json:"custom_instructions"`
	Message              string          `json:"message"`
	NotificationType     string          `json:"notification_type"`
	ErrorType            string          `json:"error_type"`
}

var claudeKinds = map[string]Kind{
	"SessionStart":       KindSessionStart,
	"SessionEnd":         KindSessionEnd,
	"UserPromptSubmit":   KindUserPrompt,
	"PreToolUse":         KindPreTool,
	"PostToolUse":        KindPostTool,
	"PostToolUseFailure": KindPostToolFailure,
	"PermissionRequest":  KindPermissionRequest,
	"SubagentStart":      KindSubagentStart,
	"SubagentStop":       KindSubagentStop,
	"Stop":               KindStop,
	"PreCompact":         KindPreCompact,
	"Notification":       KindNotification,
	"StopFailure":        KindAgentError,
}

// Decode parses a Claude Code hook stdin payload into a unified Event.
func (c *ClaudeCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	var p claudePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("claude: decode payload: %w", err)
	}
	name := p.HookEventName
	if name == "" {
		name = eventHint
	}
	kind, ok := claudeKinds[name]
	if !ok {
		kind = KindOther
	}
	ev := &Event{
		Agent:          Claude,
		Kind:           kind,
		Name:           name,
		Session:        p.SessionID,
		Cwd:            p.Cwd,
		TranscriptPath: p.TranscriptPath,
		Raw:            append(json.RawMessage(nil), raw...),
	}
	switch kind {
	case KindSessionStart:
		ev.Life = &Lifecycle{Source: p.Source, Model: p.Model}
	case KindSessionEnd:
		ev.Life = &Lifecycle{Reason: p.Reason}
	case KindUserPrompt:
		ev.Prompt = p.Prompt
	case KindPreTool, KindPermissionRequest:
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
	case KindPostTool:
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
		ev.Result = &ToolResult{Raw: cloneRaw(p.ToolResponse), Text: rawToText(p.ToolResponse)}
	case KindPostToolFailure:
		ev.Tool = newToolCall(p.ToolName, p.ToolInput, p.ToolUseID)
		ev.Result = &ToolResult{Error: p.Error}
	case KindSubagentStart:
		ev.Subagent = &Subagent{ID: p.AgentID, Type: p.AgentType}
	case KindSubagentStop:
		ev.Subagent = &Subagent{ID: p.AgentID, Type: p.AgentType, Summary: p.LastAssistantMessage}
		ev.Turn = &TurnEnd{StopHookActive: p.StopHookActive, LastAssistantMessage: p.LastAssistantMessage}
	case KindStop:
		ev.Turn = &TurnEnd{StopHookActive: p.StopHookActive, LastAssistantMessage: p.LastAssistantMessage}
	case KindPreCompact:
		ev.Compact = &CompactInfo{Trigger: p.Trigger, CustomInstructions: p.CustomInstructions}
	case KindNotification:
		ev.Note = &Note{Type: p.NotificationType, Message: p.Message}
	case KindAgentError:
		ev.Note = &Note{Type: p.ErrorType, Message: p.Message}
	}
	return ev, nil
}

// Encode renders a unified Result as Claude Code stdout JSON. Claude ignores
// exit 2 with JSON, so blocking is expressed via fields and exit code is always 0.
// ev must be non-nil.
func (c *ClaudeCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("claude: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := map[string]any{}
	hso := map[string]any{}

	if res.HaltSession {
		out["continue"] = false
		if res.Reason != "" {
			out["stopReason"] = res.Reason
		}
	}
	if res.UserMessage != "" {
		out["systemMessage"] = res.UserMessage
	}

	switch ev.Kind {
	case KindPreTool:
		if d := res.Decision.String(); d != "" {
			hso["permissionDecision"] = d
			if res.Reason != "" {
				hso["permissionDecisionReason"] = res.Reason
			}
		} else if res.UpdatedInput != nil {
			hso["permissionDecision"] = "allow"
		}
		if res.UpdatedInput != nil {
			hso["updatedInput"] = res.UpdatedInput
		}
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
	case KindPermissionRequest:
		if d := res.Decision.String(); d != "" {
			dec := map[string]any{"behavior": d}
			if res.UpdatedInput != nil {
				dec["updatedInput"] = res.UpdatedInput
			}
			if res.Reason != "" {
				dec["message"] = res.Reason
			}
			hso["decision"] = dec
		}
	case KindPostTool, KindPostToolFailure:
		if res.Decision == DecisionDeny {
			out["decision"] = "block"
			if res.Reason != "" {
				out["reason"] = res.Reason
			}
		}
		if res.UpdatedOutput != nil {
			hso["updatedToolOutput"] = *res.UpdatedOutput
		}
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
	case KindUserPrompt:
		if res.BlockPrompt || res.Decision == DecisionDeny {
			out["decision"] = "block"
			if res.Reason != "" {
				out["reason"] = res.Reason
			}
		}
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
		if res.SetTitle != "" {
			hso["sessionTitle"] = res.SetTitle
		}
	case KindStop, KindSubagentStop:
		if res.FollowUp != "" {
			out["decision"] = "block"
			out["reason"] = res.FollowUp
		}
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
	case KindSessionStart:
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
		if res.SetTitle != "" {
			hso["sessionTitle"] = res.SetTitle
		}
		if len(res.Env) > 0 {
			if err := c.writeEnvFile(res.Env); err != nil {
				return nil, 0, err
			}
		}
	default:
		if res.Context != "" {
			hso["additionalContext"] = res.Context
		}
	}

	if len(hso) > 0 {
		hso["hookEventName"] = ev.Name
		out["hookSpecificOutput"] = hso
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func (c *ClaudeCodec) writeEnvFile(env map[string]string) error {
	getenv := c.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	path := getenv("CLAUDE_ENV_FILE")
	if path == "" {
		return nil
	}
	appendFile := c.AppendFile
	if appendFile == nil {
		appendFile = func(p string, data []byte) error {
			f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // CLAUDE_ENV_FILE path from agent env
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()
			_, err = f.Write(data)
			return err
		}
	}
	var buf []byte
	for k, v := range env {
		if !validEnvKey(k) {
			return fmt.Errorf("claude: invalid env key %q", k)
		}
		buf = append(buf, []byte(fmt.Sprintf("export %s=%s\n", k, shellSingleQuote(v)))...)
	}
	return appendFile(path, buf)
}

// shellSingleQuote wraps s in single quotes using POSIX-safe escaping.
func shellSingleQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// validEnvKey reports whether k is safe to embed unquoted in a shell export name.
func validEnvKey(k string) bool {
	if k == "" {
		return false
	}
	for i, r := range k {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
