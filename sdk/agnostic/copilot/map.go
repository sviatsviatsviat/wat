package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapEvent adapts a decoded GitHub Copilot hook event into a unified Event.
func MapEvent(native sdkcopilot.Event, raw []byte) *model.Event {
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
		mapSessionStart(e, ev)
	case sdkcopilot.SessionEnd:
		mapSessionEnd(e, ev)
	case sdkcopilot.UserPromptSubmitted:
		mapUserPromptSubmitted(e, ev)
	case sdkcopilot.PreToolUse:
		mapPreToolUse(e, ev)
	case sdkcopilot.PermissionRequest:
		mapPermissionRequest(e, ev)
	case sdkcopilot.PostToolUse:
		mapPostToolUse(e, ev)
	case sdkcopilot.PostToolUseFailure:
		mapPostToolUseFailure(e, ev)
	case sdkcopilot.SubagentStart:
		mapSubagentStart(e, ev)
	case sdkcopilot.SubagentStop:
		mapSubagentStop(e, ev)
	case sdkcopilot.AgentStop:
		mapAgentStop(e, ev)
	case sdkcopilot.PreCompact:
		mapPreCompact(e, ev)
	case sdkcopilot.Notification:
		mapNotification(e, ev)
	case sdkcopilot.ErrorOccurred:
		mapErrorOccurred(e, ev)
	}
	return ev
}
