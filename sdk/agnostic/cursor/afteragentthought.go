package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterAgentThought(e sdkcursor.AfterAgentThought, ev *model.Event) {
	ev.Note = &model.Note{Type: "thought", Message: e.Text}
}
