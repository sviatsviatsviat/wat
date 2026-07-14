package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapEvent(native sdkcursor.Event, raw []byte) *model.Event {
	name := receivedName(native)
	kind, ok := KindForEvent(native.EventName())
	if !ok {
		kind = model.KindOther
	}
	env := sdkcursor.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Cursor,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}

	switch e := native.(type) {
	case sdkcursor.SessionStart:
		ev.Life = &model.Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
	case sdkcursor.SessionEnd:
		ev.Life = &model.Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
	case sdkcursor.BeforeSubmitPrompt:
		ev.Prompt = e.Prompt
	case sdkcursor.PreToolUse:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		if shell := e.ShellCommand(); shell != "" {
			ev.Tool.Shell = shell
		}
	case sdkcursor.PostToolUse:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &model.ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
	case sdkcursor.PostToolUseFailure:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &model.ToolResult{
			Error:       e.ErrorMessage,
			FailureType: e.FailureType,
			DurationMs:  e.DurationMillis(),
		}
	case sdkcursor.BeforeShellExecution:
		ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: name, Shell: e.Command}
	case sdkcursor.AfterShellExecution:
		ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: name, Shell: e.Command}
		ev.Result = &model.ToolResult{Text: e.Output, DurationMs: e.DurationMillis()}
	case sdkcursor.BeforeMCPExecution:
		nameNorm, _ := model.NormalizeToolName(e.ToolName)
		ev.Tool = &model.ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  adapter.CloneRaw(e.ToolInput),
			MCP:    true,
		}
	case sdkcursor.AfterMCPExecution:
		nameNorm, _ := model.NormalizeToolName(e.ToolName)
		ev.Tool = &model.ToolCall{
			Name:   nameNorm,
			Native: name,
			Input:  adapter.CloneRaw(e.ToolInput),
			MCP:    true,
		}
		ev.Result = &model.ToolResult{Text: e.ResultJSON, DurationMs: e.DurationMillis()}
	case sdkcursor.BeforeReadFile:
		input, err := json.Marshal(map[string]any{
			"file_path":   e.FilePath,
			"content":     e.Content,
			"attachments": e.Attachments,
		})
		if err != nil {
			return ev
		}
		ev.Tool = &model.ToolCall{Name: model.ToolRead, Native: name, Input: input}
	case sdkcursor.AfterFileEdit:
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
		ev.Tool = &model.ToolCall{Name: model.ToolEdit, Native: name, Input: input}
		ev.Result = &model.ToolResult{Raw: editsRaw}
	case sdkcursor.SubagentStart:
		ev.Subagent = &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task}
	case sdkcursor.SubagentStop:
		tp := ""
		if e.AgentTranscriptPath != nil {
			tp = *e.AgentTranscriptPath
		}
		ev.Subagent = &model.Subagent{
			ID:             e.SubagentID,
			Type:           e.SubagentType,
			Task:           e.Task,
			Summary:        e.Summary,
			Status:         e.Status,
			TranscriptPath: tp,
			LoopCount:      e.LoopCount,
		}
	case sdkcursor.Stop:
		ev.Turn = &model.TurnEnd{Status: e.Status, LoopCount: e.LoopCount}
	case sdkcursor.PreCompact:
		ev.Compact = &model.CompactInfo{Trigger: e.Trigger}
	case sdkcursor.AfterAgentResponse:
		ev.Note = &model.Note{Message: e.Text}
	case sdkcursor.AfterAgentThought:
		ev.Note = &model.Note{Type: "thought", Message: e.Text}
	}
	return ev
}

func receivedName(native sdkcursor.Event) string {
	if name := sdkcursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapOutput(ev *model.Event, res model.Result) any {
	switch {
	case permissionEvent(ev):
		out := sdkcursor.PermissionOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = sdkcursor.PermissionDecision(d)
		}
		out.UserMessage = res.UserMessage
		out.AgentMessage = res.Reason
		if res.UpdatedInput != nil && ev.Name == sdkcursor.EventPreToolUse {
			out.UpdatedInput = res.UpdatedInput
		}
		return out
	case ev.Kind == model.KindUserPrompt:
		out := sdkcursor.BeforeSubmitPromptOutput{UserMessage: res.UserMessage}
		if res.BlockPrompt || res.Decision == model.DecisionDeny {
			f := false
			out.Continue = &f
			if out.UserMessage == "" {
				out.UserMessage = res.Reason
			}
		}
		return out
	case ev.Kind == model.KindPostTool:
		out := sdkcursor.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.UpdatedMCPOutput = *res.UpdatedOutput
		}
		return out
	case ev.Kind == model.KindStop, ev.Kind == model.KindSubagentStop:
		return sdkcursor.StopOutput{FollowUpMessage: res.FollowUp}
	case ev.Kind == model.KindSessionStart:
		return sdkcursor.SessionStartOutput{Env: res.Env, AdditionalContext: res.Context}
	case ev.Kind == model.KindPreCompact:
		return sdkcursor.PreCompactOutput{UserMessage: res.UserMessage}
	default:
		return nil
	}
}

func permissionEvent(ev *model.Event) bool {
	if ev == nil {
		return false
	}
	switch ev.Kind {
	case model.KindPreTool, model.KindSubagentStart:
		return true
	case model.KindOther:
		return ev.Name == sdkcursor.EventBeforeTabFileRead
	default:
		return false
	}
}
