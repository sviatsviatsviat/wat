package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapEvent adapts a decoded Claude Code hook event into a unified Event.
func MapEvent(native sdkclaude.Event, raw []byte) *model.Event {
	name := native.EventName()
	env := sdkclaude.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Claude,
		Kind:           model.KindOther,
		Name:           name,
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            adapter.CloneRaw(raw),
	}
	switch e := native.(type) {
	case sdkclaude.SessionStart:
		ev.Kind = model.KindSessionStart
		mapSessionStart(e, ev)
	case sdkclaude.SessionEnd:
		ev.Kind = model.KindSessionEnd
		mapSessionEnd(e, ev)
	case sdkclaude.UserPromptSubmit:
		ev.Kind = model.KindUserPrompt
		mapUserPromptSubmit(e, ev)
	case sdkclaude.PreToolUse:
		ev.Kind = model.KindPreTool
		mapPreToolUse(e, ev)
	case sdkclaude.PermissionRequest:
		ev.Kind = model.KindPermissionRequest
		mapPermissionRequest(e, ev)
	case sdkclaude.PostToolUse:
		ev.Kind = model.KindPostTool
		mapPostToolUse(e, ev)
	case sdkclaude.PostToolUseFailure:
		ev.Kind = model.KindPostToolFailure
		mapPostToolUseFailure(e, ev)
	case sdkclaude.SubagentStart:
		ev.Kind = model.KindSubagentStart
		mapSubagentStart(e, ev)
	case sdkclaude.SubagentStop:
		ev.Kind = model.KindSubagentStop
		mapSubagentStop(e, ev)
	case sdkclaude.Stop:
		ev.Kind = model.KindStop
		mapStop(e, ev)
	case sdkclaude.PreCompact:
		ev.Kind = model.KindPreCompact
		mapPreCompact(e, ev)
	case sdkclaude.Notification:
		ev.Kind = model.KindNotification
		mapNotification(e, ev)
	case sdkclaude.StopFailure:
		ev.Kind = model.KindAgentError
		mapStopFailure(e, ev)
	}
	return ev
}
