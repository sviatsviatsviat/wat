package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func newEvent(native sdkcursor.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkcursor.EnvelopeOf(native)
	return &model.Event{
		Agent:          model.Cursor,
		Kind:           kind,
		Name:           receivedName(native),
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
}
