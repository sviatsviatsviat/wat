package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func envelope(env sdkcursor.Envelope, name string) model.Envelope {
	session := env.ConversationID
	if session == "" {
		session = env.SessionID
	}
	transcript := ""
	if env.TranscriptPath != nil {
		transcript = *env.TranscriptPath
	}
	return model.Envelope{
		Agent:          sdkcursor.Dialect,
		Name:           name,
		Session:        session,
		Cwd:            env.Cwd,
		TranscriptPath: transcript,
	}
}
