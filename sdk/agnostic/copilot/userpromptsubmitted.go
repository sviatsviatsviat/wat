package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapUserPromptSubmitted maps a Copilot UserPromptSubmitted hook into a unified Event.
func MapUserPromptSubmitted(e sdkcopilot.UserPromptSubmitted, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindUserPrompt)
	ev.Prompt = e.Prompt
	return ev
}
