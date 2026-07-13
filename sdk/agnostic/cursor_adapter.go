package agnostic

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapCursorEvent(native cursor.Event, raw []byte) *Event {
	name := cursorReceivedName(native)
	kind, ok := CursorKindForEvent(native.EventName())
	if !ok {
		kind = KindOther
	}
	env := cursor.EnvelopeOf(native)
	ev := &Event{
		Agent:          Cursor,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            cloneRaw(raw),
	}

	switch e := native.(type) {
	case cursor.SessionStart:
		ev.Life = &Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
	case cursor.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
	case cursor.BeforeSubmitPrompt:
		ev.Prompt = e.Prompt
	case cursor.PreToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		if shell := e.ShellCommand(); shell != "" {
			ev.Tool.Shell = shell
		}
	case cursor.PostToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
	case cursor.PostToolUseFailure:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{
			Error:       e.ErrorMessage,
			FailureType: e.FailureType,
			DurationMs:  e.DurationMillis(),
		}
	case cursor.BeforeShellExecution:
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: e.Command}
	case cursor.AfterShellExecution:
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: e.Command}
		ev.Result = &ToolResult{Text: e.Output, DurationMs: e.DurationMillis()}
	case cursor.BeforeMCPExecution:
		nameNorm, _ := NormalizeToolName(e.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(e.ToolInput),
			MCP:    true,
		}
	case cursor.AfterMCPExecution:
		nameNorm, _ := NormalizeToolName(e.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(e.ToolInput),
			MCP:    true,
		}
		ev.Result = &ToolResult{Text: e.ResultJSON, DurationMs: e.DurationMillis()}
	case cursor.BeforeReadFile:
		input, err := json.Marshal(map[string]any{
			"file_path":   e.FilePath,
			"content":     e.Content,
			"attachments": e.Attachments,
		})
		if err != nil {
			return ev
		}
		ev.Tool = &ToolCall{Name: ToolRead, Native: name, Input: input}
	case cursor.AfterFileEdit:
		input, err := json.Marshal(map[string]any{
			"file_path": e.FilePath,
			"edits":     e.Edits,
		})
		if err != nil {
			return ev
		}
		editsRaw, err := json.Marshal(e.Edits)
		if err != nil {
			return ev
		}
		ev.Tool = &ToolCall{Name: ToolEdit, Native: name, Input: input}
		ev.Result = &ToolResult{Raw: editsRaw}
	case cursor.SubagentStart:
		ev.Subagent = &Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task}
	case cursor.SubagentStop:
		tp := ""
		if e.AgentTranscriptPath != nil {
			tp = *e.AgentTranscriptPath
		}
		ev.Subagent = &Subagent{
			ID:             e.SubagentID,
			Type:           e.SubagentType,
			Task:           e.Task,
			Summary:        e.Summary,
			Status:         e.Status,
			TranscriptPath: tp,
			LoopCount:      e.LoopCount,
		}
	case cursor.Stop:
		ev.Turn = &TurnEnd{Status: e.Status, LoopCount: e.LoopCount}
	case cursor.PreCompact:
		ev.Compact = &CompactInfo{Trigger: e.Trigger}
	case cursor.AfterAgentResponse:
		ev.Note = &Note{Message: e.Text}
	case cursor.AfterAgentThought:
		ev.Note = &Note{Type: "thought", Message: e.Text}
	}
	return ev
}

func cursorReceivedName(native cursor.Event) string {
	if name := cursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapCursorOutput(ev *Event, res Result) any {
	switch {
	case cursorPermissionEvent(ev):
		out := cursor.PermissionOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = cursor.PermissionDecision(d)
		}
		out.UserMessage = res.UserMessage
		out.AgentMessage = res.Reason
		if res.UpdatedInput != nil && ev.Name == cursor.EventPreToolUse {
			out.UpdatedInput = res.UpdatedInput
		}
		return out
	case ev.Kind == KindUserPrompt:
		out := cursor.BeforeSubmitPromptOutput{UserMessage: res.UserMessage}
		if res.BlockPrompt || res.Decision == DecisionDeny {
			f := false
			out.Continue = &f
			if out.UserMessage == "" {
				out.UserMessage = res.Reason
			}
		}
		return out
	case ev.Kind == KindPostTool:
		out := cursor.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.UpdatedMCPOutput = *res.UpdatedOutput
		}
		return out
	case ev.Kind == KindStop, ev.Kind == KindSubagentStop:
		return cursor.StopOutput{FollowUpMessage: res.FollowUp}
	case ev.Kind == KindSessionStart:
		return cursor.SessionStartOutput{Env: res.Env, AdditionalContext: res.Context}
	case ev.Kind == KindPreCompact:
		return cursor.PreCompactOutput{UserMessage: res.UserMessage}
	default:
		return nil
	}
}

func cursorPermissionEvent(ev *Event) bool {
	if ev == nil {
		return false
	}
	switch ev.Kind {
	case KindPreTool, KindSubagentStart:
		return true
	case KindOther:
		return ev.Name == cursor.EventBeforeTabFileRead
	default:
		return false
	}
}

func cursorEncodeEventName(ev *Event) string {
	if ev.Name != "" {
		return ev.Name
	}
	return CursorEventForKind[ev.Kind]
}
