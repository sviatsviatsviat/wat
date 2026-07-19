package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func envelope(env sdkcopilot.Envelope, name string) model.Envelope {
	return model.Envelope{
		Agent:          sdkcopilot.Dialect,
		Name:           name,
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
	}
}
