package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapCopilotEvent(native copilot.Event, raw []byte) *Event {
	name := nativeReceivedName(native)
	kind, ok := CopilotKindForEvent(native.EventName())
	if !ok {
		kind = KindOther
	}
	env := copilot.EnvelopeOf(native)
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
	case copilot.SessionStart:
		ev.Life = &Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
	case copilot.SessionEnd:
		ev.Life = &Lifecycle{Reason: e.Reason}
	case copilot.UserPromptSubmitted:
		ev.Prompt = e.Prompt
	case copilot.PreToolUse:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
	case copilot.PermissionRequest:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
	case copilot.PostToolUse:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
		ev.Result = &ToolResult{Raw: cloneRaw(e.ResultRaw()), Text: e.ResultText()}
	case copilot.PostToolUseFailure:
		ev.Tool = newToolCall(e.NativeToolName(), e.Input(), "")
		if msg := e.ErrorMessage(); msg != "" {
			ev.Result = &ToolResult{Error: msg}
		}
	case copilot.SubagentStart:
		ev.Subagent = &Subagent{
			Type:    e.Name(),
			Task:    e.DisplayName(),
			Summary: e.AgentDescription,
		}
	case copilot.SubagentStop:
		ev.Subagent = &Subagent{
			Type: e.Name(),
			Task: e.DisplayName(),
		}
		ev.Turn = &TurnEnd{Status: e.Reason()}
	case copilot.AgentStop:
		ev.Turn = &TurnEnd{Status: e.Reason()}
	case copilot.PreCompact:
		ev.Compact = &CompactInfo{
			Trigger:            e.Trigger,
			CustomInstructions: e.Instructions(),
		}
	case copilot.Notification:
		ev.Note = &Note{Type: e.NotificationType, Title: e.Title, Message: e.Message}
	case copilot.ErrorOccurred:
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

func nativeReceivedName(native copilot.Event) string {
	if name := copilot.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapCopilotOutput(ev *Event, res Result) any {
	switch ev.Kind {
	case KindPreTool:
		out := copilot.PreToolOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = copilot.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.ModifiedArgs = res.UpdatedInput
		return out
	case KindPostTool:
		out := copilot.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.ModifiedResult = *res.UpdatedOutput
		}
		return out
	case KindStop, KindSubagentStop:
		return copilot.StopOutput{Reason: res.FollowUp}
	case KindPermissionRequest:
		if res.Decision == DecisionUnset {
			return nil
		}
		out := copilot.PermissionRequestOutput{Message: res.Reason}
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
		return copilot.PostToolFailureOutput{Context: res.Context}
	case KindSessionStart:
		return copilot.SessionStartOutput{AdditionalContext: res.Context}
	case KindSubagentStart:
		return copilot.SubagentStartOutput{AdditionalContext: res.Context}
	case KindNotification:
		return copilot.NotificationOutput{AdditionalContext: res.Context}
	default:
		return nil
	}
}
