package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapStopFailure maps a Claude StopFailure hook into a unified Event.
func MapStopFailure(e sdkclaude.StopFailure, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindAgentError)
	ev.Note = &model.Note{Type: e.ErrorType, Message: e.Message}
	return ev
}
