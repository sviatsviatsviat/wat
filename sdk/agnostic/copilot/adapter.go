package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapEvent(native sdkcopilot.Event, raw []byte) *model.Event {
	name := nativeReceivedName(native)
	kind, ok := KindForEvent(native.EventName())
	if !ok {
		kind = model.KindOther
	}
	env := sdkcopilot.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Copilot,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
	switch e := native.(type) {
	case sdkcopilot.SessionStart:
		ev.Life = &model.Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
	case sdkcopilot.SessionEnd:
		ev.Life = &model.Lifecycle{Reason: e.Reason}
	case sdkcopilot.UserPromptSubmitted:
		ev.Prompt = e.Prompt
	case sdkcopilot.PreToolUse:
		ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input(), "")
	case sdkcopilot.PermissionRequest:
		ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input(), "")
	case sdkcopilot.PostToolUse:
		ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input(), "")
		ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ResultRaw()), Text: e.ResultText()}
	case sdkcopilot.PostToolUseFailure:
		ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input(), "")
		if msg := e.ErrorMessage(); msg != "" {
			ev.Result = &model.ToolResult{Error: msg}
		}
	case sdkcopilot.SubagentStart:
		ev.Subagent = &model.Subagent{
			Type:    e.Name(),
			Task:    e.DisplayName(),
			Summary: e.AgentDescription,
		}
	case sdkcopilot.SubagentStop:
		ev.Subagent = &model.Subagent{
			Type: e.Name(),
			Task: e.DisplayName(),
		}
		ev.Turn = &model.TurnEnd{Status: e.Reason()}
	case sdkcopilot.AgentStop:
		ev.Turn = &model.TurnEnd{Status: e.Reason()}
	case sdkcopilot.PreCompact:
		ev.Compact = &model.CompactInfo{
			Trigger:            e.Trigger,
			CustomInstructions: e.Instructions(),
		}
	case sdkcopilot.Notification:
		ev.Note = &model.Note{Type: e.NotificationType, Title: e.Title, Message: e.Message}
	case sdkcopilot.ErrorOccurred:
		if detail, ok := e.Detail(); ok {
			noteType := detail.Name
			if noteType == "" {
				noteType = e.Context()
			}
			ev.Note = &model.Note{
				Type:        noteType,
				Message:     detail.Message,
				Recoverable: e.Recoverable,
			}
		}
	}
	return ev
}

func nativeReceivedName(native sdkcopilot.Event) string {
	if name := sdkcopilot.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}

func mapOutput(ev *model.Event, res model.Result) any {
	switch ev.Kind {
	case model.KindPreTool:
		out := sdkcopilot.PreToolOutput{}
		if d := res.Decision.String(); d != "" {
			out.Decision = sdkcopilot.PermissionDecision(d)
			out.Reason = res.Reason
		}
		out.ModifiedArgs = res.UpdatedInput
		return out
	case model.KindPostTool:
		out := sdkcopilot.PostToolOutput{AdditionalContext: res.Context}
		if res.UpdatedOutput != nil {
			out.ModifiedResult = *res.UpdatedOutput
		}
		return out
	case model.KindStop, model.KindSubagentStop:
		return sdkcopilot.StopOutput{Reason: res.FollowUp}
	case model.KindPermissionRequest:
		if res.Decision == model.DecisionUnset {
			return nil
		}
		out := sdkcopilot.PermissionRequestOutput{Message: res.Reason}
		switch res.Decision {
		case model.DecisionAllow:
			out.Behavior = "allow"
		case model.DecisionDeny:
			out.Behavior = "deny"
			if res.HaltSession {
				out.Interrupt = true
			}
		case model.DecisionAsk:
			out.Behavior = "deny"
			out.SuppressWarnExit = true
		default:
			return nil
		}
		return out
	case model.KindPostToolFailure:
		return sdkcopilot.PostToolFailureOutput{Context: res.Context}
	case model.KindSessionStart:
		return sdkcopilot.SessionStartOutput{AdditionalContext: res.Context}
	case model.KindSubagentStart:
		return sdkcopilot.SubagentStartOutput{AdditionalContext: res.Context}
	case model.KindNotification:
		return sdkcopilot.NotificationOutput{AdditionalContext: res.Context}
	default:
		return nil
	}
}
