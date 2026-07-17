package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapBeforeSubmitPrompt maps a Cursor BeforeSubmitPrompt hook into a unified Event.
func MapBeforeSubmitPrompt(e sdkcursor.BeforeSubmitPrompt, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindUserPrompt)
	ev.Prompt = e.Prompt
	return ev
}
