package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func newEvent(native sdkclaude.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkclaude.EnvelopeOf(native)
	return &model.Event{
		Agent:          model.Claude,
		Kind:           kind,
		Name:           native.EventName(),
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            adapter.CloneRaw(raw),
	}
}
