package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func envelope(env sdkclaude.Envelope, name string) model.Envelope {
	return model.Envelope{
		Agent:          sdkclaude.Dialect,
		Name:           name,
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
	}
}
