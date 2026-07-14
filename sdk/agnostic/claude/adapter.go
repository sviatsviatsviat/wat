package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapEvent(native sdkclaude.Event, raw []byte) *model.Event {
	name := native.EventName()
	kind, ok := KindForEvent(name)
	if !ok {
		kind = model.KindOther
	}
	env := sdkclaude.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Claude,
		Kind:           kind,
		Name:           name,
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            adapter.CloneRaw(raw),
	}
	switch e := native.(type) {
	case sdkclaude.SessionStart:
		ev.Life = &model.Lifecycle{Source: e.Source, Model: e.Model}
	case sdkclaude.SessionEnd:
		ev.Life = &model.Lifecycle{Reason: e.Reason}
	case sdkclaude.UserPromptSubmit:
		ev.Prompt = e.Prompt
	case sdkclaude.PreToolUse:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case sdkclaude.PermissionRequest:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case sdkclaude.PostToolUse:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ToolResponse), Text: adapter.RawToText(e.ToolResponse)}
	case sdkclaude.PostToolUseFailure:
		ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &model.ToolResult{Error: e.Error}
	case sdkclaude.SubagentStart:
		ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType}
	case sdkclaude.SubagentStop:
		ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage}
		ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case sdkclaude.Stop:
		ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case sdkclaude.PreCompact:
		ev.Compact = &model.CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
	case sdkclaude.Notification:
		ev.Note = &model.Note{Type: e.NotificationType, Message: e.Message}
	case sdkclaude.StopFailure:
		ev.Note = &model.Note{Type: e.ErrorType, Message: e.Message}
	}
	return ev
}

func mapOutput(ev *model.Event, res model.Result) any {
	common := sdkclaude.Common{}
	if res.HaltSession {
		f := false
		common.Continue = &f
		common.StopReason = res.Reason
	}
	if res.UserMessage != "" {
		common.SystemMessage = res.UserMessage
	}

	switch ev.Kind {
	case model.KindPreTool:
		out := sdkclaude.PreToolUseOutput{Common: common}
		if d := res.Decision.String(); d != "" {
			out.Decision = sdkclaude.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.UpdatedInput = res.UpdatedInput
		out.AdditionalContext = res.Context
		return out
	case model.KindPermissionRequest:
		out := sdkclaude.PermissionRequestOutput{Common: common}
		if d := res.Decision.String(); d != "" {
			out.Behavior = d
		}
		out.UpdatedInput = res.UpdatedInput
		out.Message = res.Reason
		out.AdditionalContext = res.Context
		if res.HaltSession {
			out.Interrupt = true
		}
		return out
	case model.KindPostTool, model.KindPostToolFailure:
		out := sdkclaude.PostToolUseOutput{Common: common}
		if res.Decision == model.DecisionDeny {
			out.Block = true
			out.Reason = res.Reason
		}
		if res.UpdatedOutput != nil {
			out.UpdatedToolOutput = *res.UpdatedOutput
		}
		out.AdditionalContext = res.Context
		return out
	case model.KindUserPrompt:
		out := sdkclaude.UserPromptSubmitOutput{Common: common}
		if res.BlockPrompt || res.Decision == model.DecisionDeny {
			out.Block = true
			out.Reason = res.Reason
		}
		out.AdditionalContext = res.Context
		out.SessionTitle = res.SetTitle
		return out
	case model.KindStop, model.KindSubagentStop:
		out := sdkclaude.StopOutput{Common: common}
		if res.FollowUp != "" {
			out.Block = true
			out.Reason = res.FollowUp
		}
		out.AdditionalContext = res.Context
		return out
	case model.KindSessionStart:
		return sdkclaude.SessionStartOutput{
			Common:            common,
			AdditionalContext: res.Context,
			SessionTitle:      res.SetTitle,
			Env:               res.Env,
		}
	default:
		return sdkclaude.CommonOutput{Common: common, AdditionalContext: res.Context}
	}
}
