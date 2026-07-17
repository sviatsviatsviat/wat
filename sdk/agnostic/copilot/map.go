package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapEvent adapts a decoded GitHub Copilot hook event into a unified Event.
func MapEvent(native sdkcopilot.Event, raw []byte) *model.Event {
	name := nativeReceivedName(native)
	env := sdkcopilot.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Copilot,
		Kind:           model.KindOther,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
	switch e := native.(type) {
	case sdkcopilot.SessionStart:
		ev.Kind = model.KindSessionStart
		mapSessionStart(e, ev)
	case sdkcopilot.SessionEnd:
		ev.Kind = model.KindSessionEnd
		mapSessionEnd(e, ev)
	case sdkcopilot.UserPromptSubmitted:
		ev.Kind = model.KindUserPrompt
		mapUserPromptSubmitted(e, ev)
	case sdkcopilot.PreToolUse:
		ev.Kind = model.KindPreTool
		mapPreToolUse(e, ev)
	case sdkcopilot.PermissionRequest:
		ev.Kind = model.KindPermissionRequest
		mapPermissionRequest(e, ev)
	case sdkcopilot.PostToolUse:
		ev.Kind = model.KindPostTool
		mapPostToolUse(e, ev)
	case sdkcopilot.PostToolUseFailure:
		ev.Kind = model.KindPostToolFailure
		mapPostToolUseFailure(e, ev)
	case sdkcopilot.SubagentStart:
		ev.Kind = model.KindSubagentStart
		mapSubagentStart(e, ev)
	case sdkcopilot.SubagentStop:
		ev.Kind = model.KindSubagentStop
		mapSubagentStop(e, ev)
	case sdkcopilot.AgentStop:
		ev.Kind = model.KindStop
		mapAgentStop(e, ev)
	case sdkcopilot.PreCompact:
		ev.Kind = model.KindPreCompact
		mapPreCompact(e, ev)
	case sdkcopilot.Notification:
		ev.Kind = model.KindNotification
		mapNotification(e, ev)
	case sdkcopilot.ErrorOccurred:
		ev.Kind = model.KindAgentError
		mapErrorOccurred(e, ev)
	}
	return ev
}
