package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapClaudeEvent(native claude.Event, raw []byte) *Event {
	name := native.EventName()
	kind, ok := ClaudeKindForEvent(name)
	if !ok {
		kind = KindOther
	}
	env := envelopeOf(native)
	ev := &Event{
		Agent:          Claude,
		Kind:           kind,
		Name:           name,
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            cloneRaw(raw),
	}
	switch e := native.(type) {
	case claude.SessionStart:
		ev.Life = &Lifecycle{Source: e.Source, Model: e.Model}
	case claude.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason}
	case claude.UserPromptSubmit:
		ev.Prompt = e.Prompt
	case claude.PreToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case claude.PermissionRequest:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case claude.PostToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Raw: cloneRaw(e.ToolResponse), Text: rawToText(e.ToolResponse)}
	case claude.PostToolUseFailure:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Error: e.Error}
	case claude.SubagentStart:
		ev.Subagent = &Subagent{ID: e.AgentID, Type: e.AgentType}
	case claude.SubagentStop:
		ev.Subagent = &Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage}
		ev.Turn = &TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case claude.Stop:
		ev.Turn = &TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case claude.PreCompact:
		ev.Compact = &CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
	case claude.Notification:
		ev.Note = &Note{Type: e.NotificationType, Message: e.Message}
	case claude.StopFailure:
		ev.Note = &Note{Type: e.ErrorType, Message: e.Message}
	}
	return ev
}

func envelopeOf(native claude.Event) claude.Envelope {
	return claude.EnvelopeOf(native)
}

func mapClaudeOutput(ev *Event, res Result) any {
	common := claude.Common{}
	if res.HaltSession {
		f := false
		common.Continue = &f
		common.StopReason = res.Reason
	}
	if res.UserMessage != "" {
		common.SystemMessage = res.UserMessage
	}

	switch ev.Kind {
	case KindPreTool:
		out := claude.PreToolUseOutput{Common: common}
		if d := res.Decision.String(); d != "" {
			out.Decision = claude.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.UpdatedInput = res.UpdatedInput
		out.AdditionalContext = res.Context
		return out
	case KindPermissionRequest:
		out := claude.PermissionRequestOutput{Common: common}
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
	case KindPostTool, KindPostToolFailure:
		out := claude.PostToolUseOutput{Common: common}
		if res.Decision == DecisionDeny {
			out.Block = true
			out.Reason = res.Reason
		}
		if res.UpdatedOutput != nil {
			out.UpdatedToolOutput = *res.UpdatedOutput
		}
		out.AdditionalContext = res.Context
		return out
	case KindUserPrompt:
		out := claude.UserPromptSubmitOutput{Common: common}
		if res.BlockPrompt || res.Decision == DecisionDeny {
			out.Block = true
			out.Reason = res.Reason
		}
		out.AdditionalContext = res.Context
		out.SessionTitle = res.SetTitle
		return out
	case KindStop, KindSubagentStop:
		out := claude.StopOutput{Common: common}
		if res.FollowUp != "" {
			out.Block = true
			out.Reason = res.FollowUp
		}
		out.AdditionalContext = res.Context
		return out
	case KindSessionStart:
		return claude.SessionStartOutput{
			Common:            common,
			AdditionalContext: res.Context,
			SessionTitle:      res.SetTitle,
			Env:               res.Env,
		}
	default:
		return claude.CommonOutput{Common: common, AdditionalContext: res.Context}
	}
}
