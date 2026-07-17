package claude

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func envelope(native sdkclaude.Event, raw []byte) model.Envelope {
	env := sdkclaude.EnvelopeOf(native)
	return model.Envelope{
		Agent:          model.Claude,
		Name:           native.EventName(),
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            hookkit.CloneRaw(raw),
	}
}
