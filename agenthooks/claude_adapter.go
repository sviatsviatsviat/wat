package agenthooks

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/claudehook"
)

func mapClaudeEvent(native claudehook.Event, raw []byte) *Event {
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
	case claudehook.SessionStart:
		ev.Life = &Lifecycle{Source: e.Source, Model: e.Model}
	case claudehook.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason}
	case claudehook.UserPromptSubmit:
		ev.Prompt = e.Prompt
	case claudehook.PreToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case claudehook.PermissionRequest:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	case claudehook.PostToolUse:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Raw: cloneRaw(e.ToolResponse), Text: rawToText(e.ToolResponse)}
	case claudehook.PostToolUseFailure:
		ev.Tool = newToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
		ev.Result = &ToolResult{Error: e.Error}
	case claudehook.SubagentStart:
		ev.Subagent = &Subagent{ID: e.AgentID, Type: e.AgentType}
	case claudehook.SubagentStop:
		ev.Subagent = &Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage}
		ev.Turn = &TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case claudehook.Stop:
		ev.Turn = &TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	case claudehook.PreCompact:
		ev.Compact = &CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
	case claudehook.Notification:
		ev.Note = &Note{Type: e.NotificationType, Message: e.Message}
	case claudehook.StopFailure:
		ev.Note = &Note{Type: e.ErrorType, Message: e.Message}
	}
	return ev
}

func envelopeOf(native claudehook.Event) claudehook.Envelope {
	switch e := native.(type) {
	case claudehook.RawEvent:
		return e.Envelope
	default:
		b, err := json.Marshal(native)
		if err != nil {
			return claudehook.Envelope{HookEventName: native.EventName()}
		}
		var env claudehook.Envelope
		_ = json.Unmarshal(b, &env)
		return env
	}
}

func mapClaudeOutput(ev *Event, res Result) any {
	common := claudehook.Common{}
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
		out := claudehook.PreToolUseOutput{Common: common}
		if d := res.Decision.String(); d != "" {
			out.Decision = claudehook.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.UpdatedInput = res.UpdatedInput
		out.AdditionalContext = res.Context
		return out
	case KindPermissionRequest:
		out := claudehook.PermissionRequestOutput{Common: common}
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
		out := claudehook.PostToolUseOutput{Common: common}
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
		out := claudehook.UserPromptSubmitOutput{Common: common}
		if res.BlockPrompt || res.Decision == DecisionDeny {
			out.Block = true
			out.Reason = res.Reason
		}
		out.AdditionalContext = res.Context
		out.SessionTitle = res.SetTitle
		return out
	case KindStop, KindSubagentStop:
		out := claudehook.StopOutput{Common: common}
		if res.FollowUp != "" {
			out.Block = true
			out.Reason = res.FollowUp
		}
		out.AdditionalContext = res.Context
		return out
	case KindSessionStart:
		return claudehook.SessionStartOutput{
			Common:            common,
			AdditionalContext: res.Context,
			SessionTitle:      res.SetTitle,
			Env:               res.Env,
		}
	default:
		return claudehook.CommonOutput{Common: common, AdditionalContext: res.Context}
	}
}
