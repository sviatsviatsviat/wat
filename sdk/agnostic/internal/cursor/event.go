package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func envelope(native sdkcursor.Event, raw []byte) model.Envelope {
	env := sdkcursor.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return model.Envelope{
		Agent:          sdkcursor.Dialect,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            hookkit.CloneRaw(raw),
	}
}

func receivedName(native sdkcursor.Event) string {
	if name := sdkcursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}
