package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterAgentResponse(e sdkcursor.AfterAgentResponse, ev *model.Event) {
	ev.Note = &model.Note{Message: e.Text}
}
