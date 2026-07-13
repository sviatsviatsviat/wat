package agenthooks

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/cursorhook"
)

func mapCursorEvent(native cursorhook.Event, raw []byte) *Event {
	name := cursorReceivedName(native)
	kind, ok := CursorKindForEvent(native.EventName())
	if !ok {
		kind = KindOther
	}
	env := cursorhook.EnvelopeOf(native)
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
	case cursorhook.SessionStart:
		ev.Life = &Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
	case cursorhook.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
	case cursorhook.BeforeSubmitPrompt:
		ev.Prompt = e.Prompt
	case cursorhook.PreToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		if shell := e.ShellCommand(); shell != "" {
			ev.Tool.Shell = shell
		}
	case cursorhook.PostToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
	case cursorhook.PostToolUseFailure:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{
			Error:       e.ErrorMessage,
			FailureType: e.FailureType,
			DurationMs:  e.DurationMillis(),
		}
	case cursorhook.BeforeShellExecution:
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: e.Command}
	case cursorhook.AfterShellExecution:
		ev.Tool = &ToolCall{Name: ToolBash, Native: name, Shell: e.Command}
		ev.Result = &ToolResult{Text: e.Output, DurationMs: e.DurationMillis()}
	case cursorhook.BeforeMCPExecution:
		nameNorm, _ := NormalizeToolName(e.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(e.ToolInput),
			MCP:    true,
		}
	case cursorhook.AfterMCPExecution:
		nameNorm, _ := NormalizeToolName(e.ToolName)
		ev.Tool = &ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  cloneRaw(e.ToolInput),
			MCP:    true,
		}
		ev.Result = &ToolResult{Text: e.ResultJSON, DurationMs: e.DurationMillis()}
	case cursorhook.BeforeReadFile:
		input, err := json.Marshal(map[string]any{
			"file_path":   e.FilePath,
			"content":     e.Content,
			"attachments": e.Attachments,
		})
		if err != nil {
			return ev
		}
		ev.Tool = &ToolCall{Name: ToolRead, Native: name, Input: input}
	case cursorhook.AfterFileEdit:
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
	case cursorhook.SubagentStart:
		ev.Subagent = &Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task}
	case cursorhook.SubagentStop:
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
	case cursorhook.Stop:
		ev.Turn = &TurnEnd{Status: e.Status, LoopCount: e.LoopCount}
	case cursorhook.PreCompact:
		ev.Compact = &CompactInfo{Trigger: e.Trigger}
	case cursorhook.AfterAgentResponse:
		ev.Note = &Note{Message: e.Text}
	case cursorhook.AfterAgentThought:
		ev.Note = &Note{Type: "thought", Message: e.Text}
	}
	return ev
}

func cursorReceivedName(native cursorhook.Event) string {
	if name := cursorhook.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapCursorOutput(ev *Event, res Result) any {
	switch {
	case cursorPermissionEvent(ev):
		out := cursorhook.PermissionOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = cursorhook.PermissionDecision(d)
		}
		out.UserMessage = res.UserMessage
		out.AgentMessage = res.Reason
		if res.UpdatedInput != nil && ev.Name == cursorhook.EventPreToolUse {
			out.UpdatedInput = res.UpdatedInput
		}
		return out
	case ev.Kind == KindUserPrompt:
		out := cursorhook.BeforeSubmitPromptOutput{UserMessage: res.UserMessage}
		if res.BlockPrompt || res.Decision == DecisionDeny {
			f := false
			out.Continue = &f
			if out.UserMessage == "" {
				out.UserMessage = res.Reason
			}
		}
		return out
	case ev.Kind == KindPostTool:
		out := cursorhook.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.UpdatedMCPOutput = *res.UpdatedOutput
		}
		return out
	case ev.Kind == KindStop, ev.Kind == KindSubagentStop:
		return cursorhook.StopOutput{FollowUpMessage: res.FollowUp}
	case ev.Kind == KindSessionStart:
		return cursorhook.SessionStartOutput{Env: res.Env, AdditionalContext: res.Context}
	case ev.Kind == KindPreCompact:
		return cursorhook.PreCompactOutput{UserMessage: res.UserMessage}
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
		return ev.Name == cursorhook.EventBeforeTabFileRead
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
