package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func envelope(native sdkcursor.Event) model.Envelope {
	env := sdkcursor.EnvelopeOf(native)
	return model.Envelope{
		Agent:          sdkcursor.Dialect,
		Name:           native.EventName(),
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
	}
}
