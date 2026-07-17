package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapEvent adapts a decoded Claude Code hook event into a unified Event.
func MapEvent(native sdkclaude.Event, raw []byte) *model.Event {
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
		mapSessionStart(e, ev)
	case sdkclaude.SessionEnd:
		mapSessionEnd(e, ev)
	case sdkclaude.UserPromptSubmit:
		mapUserPromptSubmit(e, ev)
	case sdkclaude.PreToolUse:
		mapPreToolUse(e, ev)
	case sdkclaude.PermissionRequest:
		mapPermissionRequest(e, ev)
	case sdkclaude.PostToolUse:
		mapPostToolUse(e, ev)
	case sdkclaude.PostToolUseFailure:
		mapPostToolUseFailure(e, ev)
	case sdkclaude.SubagentStart:
		mapSubagentStart(e, ev)
	case sdkclaude.SubagentStop:
		mapSubagentStop(e, ev)
	case sdkclaude.Stop:
		mapStop(e, ev)
	case sdkclaude.PreCompact:
		mapPreCompact(e, ev)
	case sdkclaude.Notification:
		mapNotification(e, ev)
	case sdkclaude.StopFailure:
		mapStopFailure(e, ev)
	}
	return ev
}
