package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapStopFailure(e sdkclaude.StopFailure, ev *model.Event) {
	ev.Note = &model.Note{Type: e.ErrorType, Message: e.Message}
}
