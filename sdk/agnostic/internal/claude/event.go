package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func envelope(native sdkclaude.Event) model.Envelope {
	env := sdkclaude.EnvelopeOf(native)
	return model.Envelope{
		Agent:          sdkclaude.Dialect,
		Name:           native.EventName(),
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
	}
}
