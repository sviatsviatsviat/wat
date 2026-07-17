package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapAfterAgentResponse maps a Cursor AfterAgentResponse hook into a unified Event.
func MapAfterAgentResponse(e sdkcursor.AfterAgentResponse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindOther)
	ev.Note = &model.Note{Message: e.Text}
	return ev
}
