package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapEvent adapts a decoded Cursor hook event into a unified Event.
func MapEvent(native sdkcursor.Event, raw []byte) *model.Event {
	name := receivedName(native)
	env := sdkcursor.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Cursor,
		Kind:           model.KindOther,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}

	switch e := native.(type) {
	case sdkcursor.SessionStart:
		ev.Kind = model.KindSessionStart
		mapSessionStart(e, ev)
	case sdkcursor.SessionEnd:
		ev.Kind = model.KindSessionEnd
		mapSessionEnd(e, ev)
	case sdkcursor.BeforeSubmitPrompt:
		ev.Kind = model.KindUserPrompt
		mapBeforeSubmitPrompt(e, ev)
	case sdkcursor.PreToolUse:
		ev.Kind = model.KindPreTool
		mapPreToolUse(e, ev)
	case sdkcursor.PostToolUse:
		ev.Kind = model.KindPostTool
		mapPostToolUse(e, ev)
	case sdkcursor.PostToolUseFailure:
		ev.Kind = model.KindPostToolFailure
		mapPostToolUseFailure(e, ev)
	case sdkcursor.BeforeShellExecution:
		ev.Kind = model.KindPreTool
		mapBeforeShellExecution(e, ev, name)
	case sdkcursor.AfterShellExecution:
		ev.Kind = model.KindPostTool
		mapAfterShellExecution(e, ev, name)
	case sdkcursor.BeforeMCPExecution:
		ev.Kind = model.KindPreTool
		mapBeforeMCPExecution(e, ev, name)
	case sdkcursor.AfterMCPExecution:
		ev.Kind = model.KindPostTool
		mapAfterMCPExecution(e, ev, name)
	case sdkcursor.BeforeReadFile:
		ev.Kind = model.KindPreTool
		mapBeforeReadFile(e, ev, name)
	case sdkcursor.AfterFileEdit:
		ev.Kind = model.KindPostTool
		mapAfterFileEdit(e, ev, name)
	case sdkcursor.SubagentStart:
		ev.Kind = model.KindSubagentStart
		mapSubagentStart(e, ev)
	case sdkcursor.SubagentStop:
		ev.Kind = model.KindSubagentStop
		mapSubagentStop(e, ev)
	case sdkcursor.Stop:
		ev.Kind = model.KindStop
		mapStop(e, ev)
	case sdkcursor.PreCompact:
		ev.Kind = model.KindPreCompact
		mapPreCompact(e, ev)
	case sdkcursor.AfterAgentResponse:
		mapAfterAgentResponse(e, ev)
	case sdkcursor.AfterAgentThought:
		mapAfterAgentThought(e, ev)
	}
	return ev
}
