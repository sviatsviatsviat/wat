package agenthooks

import (
	"github.com/sviatsviatsviat/wat/copilothook"
)

func mapCopilotEvent(native copilothook.Event, raw []byte) *Event {
	name := nativeReceivedName(native)
	kind, ok := CopilotKindForEvent(native.EventName())
	if !ok {
		kind = KindOther
	}
	env := copilothook.EnvelopeOf(native)
	ev := &Event{
		Agent:          Copilot,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            cloneRaw(raw),
	}
	switch e := native.(type) {
	case copilothook.SessionStart:
		ev.Life = &Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
	case copilothook.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason}
	case copilothook.UserPromptSubmitted:
		ev.Prompt = e.Prompt
	case copilothook.PreToolUse:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
	case copilothook.PermissionRequest:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
	case copilothook.PostToolUse:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
		ev.Result = &ToolResult{Raw: cloneRaw(e.ResultRaw()), Text: e.ResultText()}
	case copilothook.PostToolUseFailure:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
		if msg := e.ErrorMessage(); msg != "" {
			ev.Result = &ToolResult{Error: msg}
		}
	case copilothook.SubagentStart:
		ev.Subagent = &Subagent{
			Type:    e.Name(),
			Task:    e.DisplayName(),
			Summary: e.AgentDescription,
		}
	case copilothook.SubagentStop:
		ev.Subagent = &Subagent{
			Type: e.Name(),
			Task: e.DisplayName(),
		}
		ev.Turn = &TurnEnd{Status: e.Reason()}
	case copilothook.AgentStop:
		ev.Turn = &TurnEnd{Status: e.Reason()}
	case copilothook.PreCompact:
		ev.Compact = &CompactInfo{
			Trigger:            e.Trigger,
			CustomInstructions: e.Instructions(),
		}
	case copilothook.Notification:
		ev.Note = &Note{Type: e.NotificationType, Title: e.Title, Message: e.Message}
	case copilothook.ErrorOccurred:
		if detail, ok := e.Detail(); ok {
			noteType := detail.Name
			if noteType == "" {
				noteType = e.Context()
			}
			ev.Note = &Note{
				Type:        noteType,
				Message:     detail.Message,
				Recoverable: e.Recoverable,
			}
		}
	}
	return ev
}

func nativeReceivedName(native copilothook.Event) string {
	if name := copilothook.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapCopilotOutput(ev *Event, res Result) any {
	switch ev.Kind {
	case KindPreTool:
		out := copilothook.PreToolOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = copilothook.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.ModifiedArgs = res.UpdatedInput
		return out
	case KindPostTool:
		out := copilothook.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.ModifiedResult = *res.UpdatedOutput
		}
		return out
	case KindStop, KindSubagentStop:
		return copilothook.StopOutput{Reason: res.FollowUp}
	case KindPermissionRequest:
		if res.Decision == DecisionUnset {
			return nil
		}
		out := copilothook.PermissionRequestOutput{Message: res.Reason}
		switch res.Decision {
		case DecisionAllow:
			out.Behavior = "allow"
		case DecisionDeny:
			out.Behavior = "deny"
			if res.HaltSession {
				out.Interrupt = true
			}
		case DecisionAsk:
			out.Behavior = "deny"
			out.SuppressWarnExit = true
		default:
			return nil
		}
		return out
	case KindPostToolFailure:
		return copilothook.PostToolFailureOutput{Context: res.Context}
	case KindSessionStart:
		return copilothook.SessionStartOutput{AdditionalContext: res.Context}
	case KindSubagentStart:
		return copilothook.SubagentStartOutput{AdditionalContext: res.Context}
	case KindNotification:
		return copilothook.NotificationOutput{AdditionalContext: res.Context}
	default:
		return nil
	}
}
