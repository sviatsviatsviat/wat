package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func newEvent(native sdkcopilot.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkcopilot.EnvelopeOf(native)
	return &model.Event{
		Agent:          model.Copilot,
		Kind:           kind,
		Name:           nativeReceivedName(native),
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
}
