package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapEvent adapts a decoded Cursor hook event into a unified Event.
func MapEvent(native sdkcursor.Event, raw []byte) *model.Event {
	name := receivedName(native)
	kind, ok := KindForEvent(native.EventName())
	if !ok {
		kind = model.KindOther
	}
	env := sdkcursor.EnvelopeOf(native)
	ev := &model.Event{
		Agent:          model.Cursor,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}

	switch e := native.(type) {
	case sdkcursor.SessionStart:
		mapSessionStart(e, ev)
	case sdkcursor.SessionEnd:
		mapSessionEnd(e, ev)
	case sdkcursor.BeforeSubmitPrompt:
		mapBeforeSubmitPrompt(e, ev)
	case sdkcursor.PreToolUse:
		mapPreToolUse(e, ev)
	case sdkcursor.PostToolUse:
		mapPostToolUse(e, ev)
	case sdkcursor.PostToolUseFailure:
		mapPostToolUseFailure(e, ev)
	case sdkcursor.BeforeShellExecution:
		mapBeforeShellExecution(e, ev, name)
	case sdkcursor.AfterShellExecution:
		mapAfterShellExecution(e, ev, name)
	case sdkcursor.BeforeMCPExecution:
		mapBeforeMCPExecution(e, ev, name)
	case sdkcursor.AfterMCPExecution:
		mapAfterMCPExecution(e, ev, name)
	case sdkcursor.BeforeReadFile:
		mapBeforeReadFile(e, ev, name)
	case sdkcursor.AfterFileEdit:
		mapAfterFileEdit(e, ev, name)
	case sdkcursor.SubagentStart:
		mapSubagentStart(e, ev)
	case sdkcursor.SubagentStop:
		mapSubagentStop(e, ev)
	case sdkcursor.Stop:
		mapStop(e, ev)
	case sdkcursor.PreCompact:
		mapPreCompact(e, ev)
	case sdkcursor.AfterAgentResponse:
		mapAfterAgentResponse(e, ev)
	case sdkcursor.AfterAgentThought:
		mapAfterAgentThought(e, ev)
	}
	return ev
}

// MapOutput renders a unified Result as a Cursor stdout payload.
func MapOutput(ev *model.Event, res model.Result) any {
	switch ev.Kind {
	case model.KindPreTool:
		return mapPreToolOutput(ev, res)
	case model.KindPostTool:
		return mapPostToolOutput(res)
	case model.KindPostToolFailure:
		return mapPostToolFailureOutput(res)
	case model.KindStop, model.KindSubagentStop:
		return mapStopOutput(res)
	case model.KindSessionStart:
		return mapSessionStartOutput(res)
	default:
		return nil
	}
}
