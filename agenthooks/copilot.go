package agenthooks

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CopilotPreToolErrorExit is the exit code when a preToolUse handler returns
// an error. Copilot command hooks fail-closed on non-zero exits other than 2.
const CopilotPreToolErrorExit = 1

// CopilotWarnExit is exit code 2. Copilot treats it as a warning by default;
// for permissionRequest it means deny, and for postToolUseFailure it carries
// additionalContext in stdout.
const CopilotWarnExit = 2

// CopilotCodec implements Codec for GitHub Copilot hooks in camelCase CLI and
// VS Code compatible (PascalCase event name, snake_case fields) formats.
//
// Handler errors on preToolUse should exit CopilotPreToolErrorExit (fail-closed).
// Encode returns CopilotWarnExit only for documented output paths on
// permissionRequest deny and postToolUseFailure context.
//
// Reference: https://docs.github.com/en/copilot/reference/hooks-reference
type CopilotCodec struct{}

// Dialect returns Copilot.
func (c *CopilotCodec) Dialect() Dialect { return Copilot }

type copilotFormat int

const (
	copilotFormatUnknown copilotFormat = iota
	copilotFormatCamel
	copilotFormatVSCode
)

// copilotTimestamp accepts ms-epoch numbers (camelCase) and ISO-8601 strings (VS Code).
type copilotTimestamp struct {
	time.Time
}

// UnmarshalJSON accepts ms-epoch numbers or ISO-8601 strings.
func (t *copilotTimestamp) UnmarshalJSON(data []byte) error {
	data = bytesTrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		parsed, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil
		}
		t.Time = parsed
		return nil
	}
	var ms json.Number
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}
	n, err := ms.Int64()
	if err != nil {
		return nil
	}
	t.Time = time.UnixMilli(n)
	return nil
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

type copilotToolResult struct {
	ResultType       string `json:"resultType"`
	TextResultForLLM string `json:"textResultForLlm"`
	ResultTypeSnake  string `json:"result_type"`
	TextResultSnake  string `json:"text_result_for_llm"`
}

func (r copilotToolResult) text() string {
	if r.TextResultForLLM != "" {
		return r.TextResultForLLM
	}
	return r.TextResultSnake
}

type copilotErrorDetail struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Stack   string `json:"stack"`
}

// copilotPayload is a superset decode of Copilot hook events in both wire formats.
type copilotPayload struct {
	HookEventName           string            `json:"hook_event_name"`
	SessionID               string            `json:"session_id"`
	SessionIDCaml           string            `json:"sessionId"`
	Timestamp               copilotTimestamp  `json:"timestamp"`
	Cwd                     string            `json:"cwd"`
	TranscriptPath          string            `json:"transcript_path"`
	TranscriptCamel         string            `json:"transcriptPath"`
	Source                  string            `json:"source"`
	InitialPrompt           string            `json:"initial_prompt"`
	InitialPromptCamel      string            `json:"initialPrompt"`
	Reason                  string            `json:"reason"`
	Prompt                  string            `json:"prompt"`
	ToolName                string            `json:"tool_name"`
	ToolNameCamel           string            `json:"toolName"`
	ToolInput               json.RawMessage   `json:"tool_input"`
	ToolArgs                json.RawMessage   `json:"toolArgs"`
	ToolResult              copilotToolResult `json:"toolResult"`
	ToolResultSnake         copilotToolResult `json:"tool_result"`
	Error                   json.RawMessage   `json:"error"`
	ErrorContext            string            `json:"error_context"`
	ErrorContextCamel       string            `json:"errorContext"`
	Recoverable             *bool             `json:"recoverable"`
	StopReason              string            `json:"stop_reason"`
	StopReasonCamel         string            `json:"stopReason"`
	AgentName               string            `json:"agent_name"`
	AgentNameCamel          string            `json:"agentName"`
	AgentDisplayName        string            `json:"agent_display_name"`
	AgentDisplayNameCamel   string            `json:"agentDisplayName"`
	AgentDescription        string            `json:"agentDescription"`
	Trigger                 string            `json:"trigger"`
	CustomInstructions      string            `json:"custom_instructions"`
	CustomInstructionsCamel string            `json:"customInstructions"`
	Message                 string            `json:"message"`
	Title                   string            `json:"title"`
	NotificationType        string            `json:"notification_type"`
}

var copilotKinds = map[string]Kind{
	"sessionStart":        KindSessionStart,
	"SessionStart":        KindSessionStart,
	"sessionEnd":          KindSessionEnd,
	"SessionEnd":          KindSessionEnd,
	"userPromptSubmitted": KindUserPrompt,
	"UserPromptSubmit":    KindUserPrompt,
	"preToolUse":          KindPreTool,
	"PreToolUse":          KindPreTool,
	"postToolUse":         KindPostTool,
	"PostToolUse":         KindPostTool,
	"postToolUseFailure":  KindPostToolFailure,
	"PostToolUseFailure":  KindPostToolFailure,
	"permissionRequest":   KindPermissionRequest,
	"PermissionRequest":   KindPermissionRequest,
	"subagentStart":       KindSubagentStart,
	"SubagentStart":       KindSubagentStart,
	"subagentStop":        KindSubagentStop,
	"SubagentStop":        KindSubagentStop,
	"agentStop":           KindStop,
	"Stop":                KindStop,
	"preCompact":          KindPreCompact,
	"PreCompact":          KindPreCompact,
	"notification":        KindNotification,
	"Notification":        KindNotification,
	"errorOccurred":       KindAgentError,
	"ErrorOccurred":       KindAgentError,
}

func sniffCopilotFormat(raw []byte) copilotFormat {
	var peek struct {
		HookEventName string `json:"hook_event_name"`
		SessionID     string `json:"sessionId"`
	}
	if json.Unmarshal(raw, &peek) != nil {
		return copilotFormatUnknown
	}
	if peek.HookEventName != "" {
		return copilotFormatVSCode
	}
	if peek.SessionID != "" {
		return copilotFormatCamel
	}
	return copilotFormatUnknown
}

func (p *copilotPayload) session() string {
	if p.SessionID != "" {
		return p.SessionID
	}
	return p.SessionIDCaml
}

func (p *copilotPayload) transcript() string {
	if p.TranscriptPath != "" {
		return p.TranscriptPath
	}
	return p.TranscriptCamel
}

func (p *copilotPayload) toolName() string {
	if p.ToolName != "" {
		return p.ToolName
	}
	return p.ToolNameCamel
}

func (p *copilotPayload) toolInput() json.RawMessage {
	if len(p.ToolInput) > 0 {
		return p.ToolInput
	}
	return p.ToolArgs
}

func (p *copilotPayload) initialPrompt() string {
	if p.InitialPrompt != "" {
		return p.InitialPrompt
	}
	return p.InitialPromptCamel
}

func (p *copilotPayload) stopReason() string {
	if p.StopReason != "" {
		return p.StopReason
	}
	return p.StopReasonCamel
}

func (p *copilotPayload) agentName() string {
	if p.AgentName != "" {
		return p.AgentName
	}
	return p.AgentNameCamel
}

func (p *copilotPayload) agentDisplayName() string {
	if p.AgentDisplayName != "" {
		return p.AgentDisplayName
	}
	return p.AgentDisplayNameCamel
}

func (p *copilotPayload) customInstructions() string {
	if p.CustomInstructions != "" {
		return p.CustomInstructions
	}
	return p.CustomInstructionsCamel
}

func (p *copilotPayload) errorContext() string {
	if p.ErrorContext != "" {
		return p.ErrorContext
	}
	return p.ErrorContextCamel
}

func (p *copilotPayload) toolResultText() string {
	if t := p.ToolResult.text(); t != "" {
		return t
	}
	return p.ToolResultSnake.text()
}

func (p *copilotPayload) toolResultRaw() json.RawMessage {
	if p.ToolResult.TextResultForLLM != "" || p.ToolResult.ResultType != "" {
		b, err := json.Marshal(p.ToolResult)
		if err == nil {
			return b
		}
	}
	if p.ToolResultSnake.TextResultSnake != "" || p.ToolResultSnake.ResultTypeSnake != "" {
		b, err := json.Marshal(p.ToolResultSnake)
		if err == nil {
			return b
		}
	}
	return nil
}

func (p *copilotPayload) errorDetail() (copilotErrorDetail, bool) {
	if len(p.Error) == 0 {
		return copilotErrorDetail{}, false
	}
	var s string
	if json.Unmarshal(p.Error, &s) == nil {
		return copilotErrorDetail{Message: s}, true
	}
	var detail copilotErrorDetail
	if json.Unmarshal(p.Error, &detail) != nil {
		return copilotErrorDetail{}, false
	}
	return detail, true
}

// Decode parses a GitHub Copilot hook stdin payload into a unified Event.
func (c *CopilotCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	format := sniffCopilotFormat(raw)
	if format == copilotFormatUnknown {
		return nil, fmt.Errorf("copilot: decode payload: unrecognized format")
	}

	var p copilotPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("copilot: decode payload: %w", err)
	}

	name := p.HookEventName
	if name == "" {
		name = eventHint
	}
	if name == "" {
		return nil, fmt.Errorf("copilot: decode: event name required (camelCase payloads need eventHint)")
	}

	kind, ok := copilotKinds[name]
	if !ok {
		kind = KindOther
	}

	ev := &Event{
		Agent:          Copilot,
		Kind:           kind,
		Name:           name,
		Session:        p.session(),
		Cwd:            p.Cwd,
		TranscriptPath: p.transcript(),
		Raw:            append(json.RawMessage(nil), raw...),
	}

	switch kind {
	case KindSessionStart:
		ev.Life = &Lifecycle{Source: p.Source, InitialPrompt: p.initialPrompt()}
	case KindSessionEnd:
		ev.Life = &Lifecycle{Reason: p.Reason}
	case KindUserPrompt:
		ev.Prompt = p.Prompt
	case KindPreTool, KindPermissionRequest:
		ev.Tool = newToolCall(p.toolName(), p.toolInput(), "")
	case KindPostTool:
		ev.Tool = newToolCall(p.toolName(), p.toolInput(), "")
		resultRaw := p.toolResultRaw()
		ev.Result = &ToolResult{Raw: cloneRaw(resultRaw), Text: p.toolResultText()}
	case KindPostToolFailure:
		ev.Tool = newToolCall(p.toolName(), p.toolInput(), "")
		if detail, ok := p.errorDetail(); ok {
			ev.Result = &ToolResult{Error: detail.Message}
		} else if len(p.Error) > 0 {
			ev.Result = &ToolResult{Error: string(p.Error)}
		}
	case KindSubagentStart:
		ev.Subagent = &Subagent{
			Type:    p.agentName(),
			Task:    p.agentDisplayName(),
			Summary: p.AgentDescription,
		}
	case KindSubagentStop:
		ev.Subagent = &Subagent{
			Type: p.agentName(),
			Task: p.agentDisplayName(),
		}
		ev.Turn = &TurnEnd{Status: p.stopReason()}
	case KindStop:
		ev.Turn = &TurnEnd{Status: p.stopReason()}
	case KindPreCompact:
		ev.Compact = &CompactInfo{
			Trigger:            p.Trigger,
			CustomInstructions: p.customInstructions(),
		}
	case KindNotification:
		ev.Note = &Note{Type: p.NotificationType, Title: p.Title, Message: p.Message}
	case KindAgentError:
		if detail, ok := p.errorDetail(); ok {
			noteType := detail.Name
			if noteType == "" {
				noteType = p.errorContext()
			}
			ev.Note = &Note{
				Type:        noteType,
				Message:     detail.Message,
				Recoverable: p.Recoverable,
			}
		}
	}

	return ev, nil
}

// Encode renders a unified Result as Copilot stdout JSON and exit code.
// ev must be non-nil.
func (c *CopilotCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("copilot: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}

	switch ev.Kind {
	case KindPreTool:
		return c.encodePreTool(res)
	case KindPostTool:
		return c.encodePostTool(res)
	case KindStop, KindSubagentStop:
		return c.encodeStop(res)
	case KindPermissionRequest:
		return c.encodePermissionRequest(res)
	case KindPostToolFailure:
		return c.encodePostToolFailure(res)
	case KindSessionStart, KindSubagentStart, KindNotification:
		return c.encodeAdditionalContext(res)
	default:
		return nil, 0, nil
	}
}

func (c *CopilotCodec) encodePreTool(res Result) ([]byte, int, error) {
	out := map[string]any{}
	if d := res.Decision.String(); d != "" {
		out["permissionDecision"] = d
		if res.Reason != "" {
			out["permissionDecisionReason"] = res.Reason
		}
	}
	if res.UpdatedInput != nil {
		out["modifiedArgs"] = res.UpdatedInput
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func (c *CopilotCodec) encodePostTool(res Result) ([]byte, int, error) {
	out := map[string]any{}
	if res.UpdatedOutput != nil {
		out["modifiedResult"] = map[string]any{
			"resultType":       "success",
			"textResultForLlm": *res.UpdatedOutput,
		}
	}
	if res.Context != "" {
		out["additionalContext"] = res.Context
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func (c *CopilotCodec) encodeStop(res Result) ([]byte, int, error) {
	if res.FollowUp == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"decision": "block",
		"reason":   res.FollowUp,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func (c *CopilotCodec) encodePermissionRequest(res Result) ([]byte, int, error) {
	if res.Decision == DecisionUnset {
		return nil, 0, nil
	}
	out := map[string]any{}
	switch res.Decision {
	case DecisionAllow:
		out["behavior"] = "allow"
	case DecisionDeny:
		out["behavior"] = "deny"
		if res.Reason != "" {
			out["message"] = res.Reason
		}
		if res.HaltSession {
			out["interrupt"] = true
		}
	case DecisionAsk:
		out["behavior"] = "deny"
		if res.Reason != "" {
			out["message"] = res.Reason
		}
	default:
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if res.Decision == DecisionDeny {
		exitCode = CopilotWarnExit
	}
	return b, exitCode, nil
}

func (c *CopilotCodec) encodePostToolFailure(res Result) ([]byte, int, error) {
	if res.Context == "" {
		return nil, 0, nil
	}
	return []byte(res.Context), CopilotWarnExit, nil
}

func (c *CopilotCodec) encodeAdditionalContext(res Result) ([]byte, int, error) {
	if res.Context == "" {
		return nil, 0, nil
	}
	out := map[string]any{"additionalContext": res.Context}
	b, err := json.Marshal(out)
	return b, 0, err
}
