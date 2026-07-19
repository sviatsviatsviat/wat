package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func envelope(native sdkcopilot.Event) model.Envelope {
	env := sdkcopilot.EnvelopeOf(native)
	return model.Envelope{
		Agent:          sdkcopilot.Dialect,
		Name:           native.EventName(),
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
	}
}
