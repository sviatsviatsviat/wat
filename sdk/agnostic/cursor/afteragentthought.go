package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapAfterAgentThought maps a Cursor AfterAgentThought hook into a unified Event.
func MapAfterAgentThought(e sdkcursor.AfterAgentThought, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindOther)
	ev.Note = &model.Note{Type: "thought", Message: e.Text}
	return ev
}
