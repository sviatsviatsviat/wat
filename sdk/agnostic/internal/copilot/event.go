package copilot

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func envelope(native sdkcopilot.Event, raw []byte) model.Envelope {
	env := sdkcopilot.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return model.Envelope{
		Agent:          sdkcopilot.Dialect,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            hookkit.CloneRaw(raw),
	}
}
