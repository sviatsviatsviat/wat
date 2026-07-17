package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapUserPromptSubmit maps a Claude UserPromptSubmit hook into a unified Event.
func MapUserPromptSubmit(e sdkclaude.UserPromptSubmit, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindUserPrompt)
	ev.Prompt = e.Prompt
	return ev
}
