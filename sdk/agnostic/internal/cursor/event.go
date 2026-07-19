package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func envelope(env sdkcursor.Envelope, name string) model.Envelope {
	return model.Envelope{
		Agent:          sdkcursor.Dialect,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
	}
}
